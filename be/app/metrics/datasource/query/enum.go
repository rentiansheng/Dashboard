package query

import (
	"fmt"

	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"strings"
)

func EnumES(ctx context.Context, dataSourceName string, req *define.DataSourceFieldEnumReq) (*define.DataSourceFieldEnum, errors.Error) {
	switch dataSourceName {

	}
	return nil, nil
}

func DataSourceEnum(ctx context.Context, metas []define.MetricDataSourceMeta, input *define.QueryReq,
	results []map[string]interface{}) (define.DataSourceDetailMeta, errors.Error) {
	detailMeta := define.DataSourceDetailMeta{
		Fields:            nil,
		Enum:              make(map[string]map[string]interface{}, 0),
		CopyFieldValue:    nil,
		CopyRawFieldValue: nil,
	}
	enumInput := define.DataSourceMetaEnumReq{}
	if err := ctx.Mapper("build data source meta enum api params", input, &enumInput); err != nil {
		return detailMeta, err
	}

	for _, meta := range metas {
		item := define.MetricDetailMeta{}
		if err := ctx.Mapper("convert field meta", meta, &item); err != nil {
			ctx.Log().ErrorJSON("convert field meta err. from: %s, to: %s, err: %s", meta, item, err)
			return detailMeta, err
		}
		// 修改 detail类型，而不是查询类型
		item.Typ = meta.OutputDataType
		item.Template = strings.TrimSpace(item.Template)
		item.Formatter = strings.TrimSpace(item.Formatter)

		detailMeta.Fields = append(detailMeta.Fields, item)
		if meta.Action.API {
			// TODO: use mapper
			values := make([]interface{}, 0)
			for _, row := range results {
				if val, ok := row[meta.FieldName]; ok && val != nil {
					values = append(values, val)
				}
			}
			if len(values) > 0 {
				enumInput.FieldName = meta.FieldName
				// 这里需要转换，填写值没有意义
				enumInput.FieldValue = values
				// TODO:  http request
				enumResp := define.DataSourceMetaFieldEnumData{}
				//if err != nil {
				//	return detailMeta, err
				//}
				enumRela := make(map[string]interface{}, 0)
				if enumResp.Typ == define.DataSourceMetaFieldEnumDataTypeKv {

					if err := ctx.Mapper("enum relation", enumResp.Values, &enumRela); err != nil {
						return detailMeta, err
					}

				} else {
					// array 不需要处理，值和value 是一样

				}
				detailMeta.Enum[meta.FieldName] = enumRela
			}
		} else if meta.Enum.Values.Typ == define.DataSourceFieldEnumTypeKv {
			tmpEnum := make(map[string]interface{}, 0)
			if err := ctx.Mapper("data source enum kv", meta.Enum.Values.Values, &tmpEnum); err != nil {
				ctx.Log().ErrorJSON("field: %s, enum: %s, err: %s", meta.FieldName, meta.Enum, err.Error())
				return detailMeta, err
			}
			// swap tmpEnum object key and value
			detailMeta.Enum[meta.FieldName] = make(map[string]interface{}, len(tmpEnum))
			for key, value := range tmpEnum {
				detailMeta.Enum[meta.FieldName][fmt.Sprintf("%v", value)] = key
			}

		}
	}

	return detailMeta, nil
}

func DataSourceEnumField(ctx context.Context, meta define.MetricDataSourceMeta, values []interface{}) (map[string]interface{}, errors.Error) {

	result := make(map[string]interface{}, 0)
	if meta.Action.API {
		enumInput := define.DataSourceMetaEnumReq{
			FieldName:    meta.FieldName,
			FieldValue:   values,
			DataSourceId: meta.DataSourceID,
		}

		enumInput.FieldName = meta.FieldName
		// 这里需要转换，填写值没有意义
		enumInput.FieldValue = values
		// TODO:  http request
		enumResp := define.DataSourceMetaFieldEnumData{}
		//if err != nil {
		//	return detailMeta, err
		//}

		if enumResp.Typ == define.DataSourceMetaFieldEnumDataTypeKv {
			tmpEnumRela := make(map[string]interface{}, 0)

			if err := ctx.Mapper("enum relation", enumResp.Values, &tmpEnumRela); err != nil {
				return nil, err
			}
			for key, value := range tmpEnumRela {
				result[fmt.Sprintf("%v", value)] = key
			}

		} else {
			// array 不需要处理，值和value 是一样

		}

	} else if meta.Enum.Values.Typ == define.DataSourceFieldEnumTypeKv {
		tmpEnum := make(map[string]interface{}, 0)
		if err := ctx.Mapper("data source enum kv", meta.Enum.Values.Values, &tmpEnum); err != nil {
			ctx.Log().ErrorJSON("field: %s, enum: %s, err: %s", meta.FieldName, meta.Enum, err.Error())
			return nil, err
		}
		// swap tmpEnum object key and value
		for key, value := range tmpEnum {
			result[fmt.Sprintf("%v", value)] = key
		}

	}

	return result, nil
}
