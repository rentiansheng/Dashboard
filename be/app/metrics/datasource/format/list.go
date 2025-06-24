package format

import (
	"fmt"

	"github.com/rentiansheng/dashboard/app/metrics/datasource/tools"

	"github.com/flosch/pongo2/v4"
	"github.com/rentiansheng/dashboard/app/metrics/datasource/format/plugins"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"github.com/rentiansheng/passion/lib/array"
)

func AdjustChartDetailPageInfoRowByFormatterTemplate(ctx context.Context, row map[string]interface{},
	detailMeta define.DataSourceDetailMeta, name string) (map[string]interface{}, errors.Error) {
	// 保证有内容返回
	// 值做枚举转换
	// enum 字段转换值

	// 最终值
	result := make(map[string]interface{}, 0)

	if row == nil {
		row = make(map[string]interface{}, 0)
	}

	// 虚拟字段， copy 方法
	copyFieldName := func(fields map[string]string) errors.Error {
		for toFieldName, fromFieldName := range fields {
			if _, ok := row[toFieldName]; ok {
				ctx.Log().Errorf("copy field error. to field has exist. fields: %s, from: %s, to: %s", fields, fromFieldName, toFieldName)
				err := ctx.Error().Errorf(code.MetricDetailCopyRowFieldErrCode, name, fromFieldName, toFieldName)
				return err
			}
			row[toFieldName] = row[fromFieldName]
		}
		return nil
	}

	// 嵌套字段，extra nested field
	extraExistsFieldRela := make(map[string]bool, 0)
	for _, field := range detailMeta.Fields {
		if len(field.Nested) == 0 {
			continue
		}
		// Notice: only one level
		nestedPath := field.Nested
		if _, ok := extraExistsFieldRela[nestedPath]; ok {
			// 同一个字段仅需要展开一次
			continue
		}

		extraExistsFieldRela[nestedPath] = true
		subVal := row[nestedPath]
		if subVal == nil {
			continue
		}
		// es nested key  must be string. es 是不区分数组和对象。这里需要做处理。 使用最后一个元素的值

		var subValKv map[string]interface{}
		if kv, ok := subVal.(map[string]interface{}); ok {
			subValKv = kv
		} else if kvs, ok := subVal.([]map[string]interface{}); ok {
			subValKv = kvs[len(kvs)-1]
		}
		if subValKv != nil {
			for k, v := range subValKv {
				subValKey := fmt.Sprintf("%s.%s", nestedPath, k)
				row[subValKey] = v
			}
		}
	}

	// 处理虚拟字段
	if err := copyFieldName(detailMeta.CopyRawFieldValue); err != nil {
		return nil, err
	}

	for fieldName, virtualField := range detailMeta.VirtualFieldFn {
		waitRow := map[string]interface{}{}
		for field, newField := range virtualField.ValueFields {
			waitRow[newField] = row[field]
		}
		fn := plugins.GetVirtualField(virtualField.FuncName)
		if fn == nil {
			return nil, ctx.Error().Errorf(code.InstNotFoundKVCode, "get virtual field fn", virtualField.FuncName)
		}
		newVal, err := fn(ctx, waitRow)
		if err != nil {
			return nil, err
		}
		row[fieldName] = newVal
	}

	// 枚举转换未展示只
	for key, val := range row {
		if enumVal, ok := detailMeta.Enum[key]; ok {
			if valArrI, ok := val.([]interface{}); ok {
				for _, valI := range valArrI {
					strVal := fmt.Sprintf("%v", valI)
					row[key] = enumVal[strVal]
				}
			} else {
				mapperValArrI := make([]interface{}, 0)
				if err := ctx.Mapper("", val, &mapperValArrI); err == nil {
					for _, valI := range mapperValArrI {
						strVal := fmt.Sprintf("%v", valI)
						row[key] = enumVal[strVal]
					}
				} else {
					strVal := fmt.Sprintf("%v", val)
					if displayVal, ok := enumVal[strVal]; ok {
						row[key] = displayVal
					}
				}

			}

		}

	}
	// 处理虚拟字段
	if err := copyFieldName(detailMeta.CopyFieldValue); err != nil {
		return nil, err
	}

	if val, ok := row[define.MetricDetailIDFieldName]; ok {
		result[define.MetricDetailIDFieldName] = val
	}

	for _, field := range detailMeta.Fields {
		rawValue := row[field.Name]

		displayValue := rawValue
		if field.Formatter != "" {
			displayValue = fmt.Sprintf(field.Formatter, rawValue)
		} else if field.Template != "" {
			tpl, err := pongo2.FromString(field.Template)
			if err != nil {
				ctx.Log().ErrorJSON("template from string error. field: %s, data: %s, err: %s", field, row, err)
				return nil, ctx.Error().Errorf(code.TemplateFormatErr, err.Error())
			}
			displayValue, err = tpl.Execute(row)
			if err != nil {
				ctx.Log().ErrorJSON("template build  error. field: %s, data: %s, err: %s", field, row, err)
				return nil, ctx.Error().Errorf(code.TemplateBuildErr, err.Error())
			}
		}

		switch {
		case field.IsLink():
			result[field.Name] = field.LinkValue(rawValue, displayValue)
		case field.IsSeconds2Date():
			seconds, err := array.Int.ToUint64(displayValue)
			if err != nil {
				ctx.Log().ErrorJSON("template build  error, field: %s, value: %s, err: %s", field, displayValue, err)
				return nil, ctx.Error().Errorf(code.TemplateBuildErr, err.Error())
			}
			result[field.Name] = tools.Second2DateInterval(seconds)
		default:
			result[field.Name] = displayValue
		}
	}

	return result, nil
}
