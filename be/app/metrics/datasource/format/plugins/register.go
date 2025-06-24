package plugins

import (
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
)

// metricDetailAPIFn 获取项内容的数据
type metricDetailAPIFn func(ctx context.Context, fields []string, value []string, offset, limit int64) ([]map[string]interface{}, errors.Error)

// enumHookFn 字段枚举指
type chartDetailEnumAPIFn func(ctx context.Context) (map[string]interface{}, errors.Error)

type virtualFieldFn func(ctx context.Context, row map[string]interface{}) (interface{}, errors.Error)

var (
	chartDetailEnumHandle   = make(map[string]chartDetailEnumAPIFn, 0)
	chartVirtualFieldHandle = make(map[string]virtualFieldFn, 0)
)

// RegisterDetailEnumAPI 新加一种获取枚举的值
func RegisterDetailEnumAPI(enumName string, fn chartDetailEnumAPIFn) {
	if _, ok := chartDetailEnumHandle[enumName]; ok {
		panic(enumName + " duplicate chart detail enum api handler")
	}
	chartDetailEnumHandle[enumName] = fn
}

// GetDetailEnumAPI  根据名字获取对应获取枚举的方法
func GetDetailEnumAPI(detailAPIName string) chartDetailEnumAPIFn {
	if handler := chartDetailEnumHandle[detailAPIName]; handler != nil {
		return handler
	}
	return nil
}

func RegisterVirtualField(virtualField string, fn virtualFieldFn) {
	if _, ok := chartVirtualFieldHandle[virtualField]; ok {
		panic(virtualField + " duplicate chart detail enum api handler")
	}
	chartVirtualFieldHandle[virtualField] = fn
}

func GetVirtualField(virtualField string) virtualFieldFn {
	if handler := chartVirtualFieldHandle[virtualField]; handler != nil {
		return handler
	}
	return nil
}
