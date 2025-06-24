package es

import (
	"github.com/rentiansheng/dashboard/app/metrics/datasource/filter/adjust"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"github.com/rentiansheng/ges"
)

const (
	RuleOperatorEQ      = define.QueryFilterRuleEQ      //"eq"
	RuleOperatorIN      = define.QueryFilterRuleIn      //"in"
	RuleOperatorBetween = define.QueryFilterRuleBetween //"between"
	RuleOperatorLT      = define.QueryFilterRuleLt      // "lt"
	RuleOperatorLTE     = define.QueryFilterRuleLte     //"lte"
	RuleOperatorGT      = define.QueryFilterRuleGt      //"gt"
	RuleOperatorGTE     = define.QueryFilterRuleGte     //"gte"
	// RuleOperatorLike like查询
	RuleOperatorLike = define.QueryFilterRuleLike //"like"
	// RuleOperatorPrefix 前缀模糊查询
	RuleOperatorPrefix = define.QueryFilterRulePrefix
	// RuleOperatorSuffix 后缀模糊查询
	RuleOperatorSuffix = define.QueryFilterRuleSuffix
)

type indexFn struct {
	fieldAdjustFn             map[string]adjust.FieldAdjustFn
	extraAdjustFn             adjust.ExtraFilterFn
	virtualFieldExtraAdjustFn map[string]adjust.VirtualFieldExtraFilterFn
}

func (i indexFn) FieldExtraAdjust(field string) (adjust.FieldAdjustFn, bool) {
	fn, ok := i.fieldAdjustFn[field]
	return fn, ok
}

func (i indexFn) ExtraAdjust(ctx context.Context, info define.MetricDataSource, groupKeyId uint64) (ges.Filter, errors.Error) {
	return i.extraAdjustFn(ctx, info, groupKeyId)
}

func (i indexFn) VirtualFieldExtraAdjust(field string) (adjust.VirtualFieldExtraFilterFn, bool) {
	fn, ok := i.virtualFieldExtraAdjustFn[field]
	return fn, ok
}

type IndexHandle interface {
	FieldExtraAdjust(field string) (adjust.FieldAdjustFn, bool)
	//ExtraAdjust  数据限定函数
	ExtraAdjust(ctx context.Context, info define.MetricDataSource, groupKeyId uint64) (ges.Filter, errors.Error)
	VirtualFieldExtraAdjust(field string) (adjust.VirtualFieldExtraFilterFn, bool)
}

func GetIndexHandle(ctx context.Context, info define.MetricDataSource) (IndexHandle, errors.Error) {

	if info.GroupKey == "" {
		return nil, ctx.Error().Errorf(code.MetricDataSourceUnsupportedErrCode, info.DisplayName)
	}

	fn := indexFn{
		fieldAdjustFn:             make(map[string]adjust.FieldAdjustFn),
		extraAdjustFn:             adjust.ExtraGroupKeyFilterFn,
		virtualFieldExtraAdjustFn: make(map[string]adjust.VirtualFieldExtraFilterFn),
	}

	return fn, nil
}

// QueryParams 查询参数处理
func QueryParams(ctx context.Context, handler IndexHandle, filters []define.QueryFilterRule) (esFilters, esNotFilters []ges.Filter, err errors.Error) {

	esFilters = make([]ges.Filter, 0, len(filters))
	esNotFilters = make([]ges.Filter, 0, len(filters))
	for _, rule := range filters {

		// 虚拟字段直接生成查询条件， 无需再次处理
		if fn, ok := handler.VirtualFieldExtraAdjust(rule.Field); ok {
			virtualFilter, err := fn(ctx, rule)
			if err != nil {
				return nil, nil, err
			}
			esFilters = append(esFilters, virtualFilter)
			continue
		}

		if fn, ok := handler.FieldExtraAdjust(rule.Field); ok {
			newRule, err := fn(ctx, rule)
			if err != nil {
				return nil, nil, err
			}
			// 使用新的规则
			rule = newRule
		}
		//  需要多次判断的原因是， 有字段是virtual field， 经过adjust后， 会变成pause 字段
		if rule.Pause {
			continue
		}
		if len(rule.Values) == 0 {
			return nil, nil, ctx.Error().Errorf(code.MetricDataSourceQueryFilterNeedValueErrCode, rule.Values, rule.Operator)
		}
		fn, ok := esCondHandle[rule.Operator]
		if !ok {
			return nil, nil, ctx.Error().Errorf(code.MetricDataSourceQueryFilterRoleNotErrCode, rule.Operator)
		}

		filter, err := fn(ctx, rule.Field, rule.Values)
		if err != nil {
			return nil, nil, err
		}
		if len(rule.Nesteds) > 0 {
			filter = ges.NewFilter().Nested(ges.Nested().Path(rule.Nesteds[0]).Must(filter))
		}
		if rule.Not {
			esNotFilters = append(esNotFilters, filter)
		} else {
			esFilters = append(esFilters, filter)

		}
	}
	return esFilters, esNotFilters, nil
}

type esCondHandleFn func(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error)

var esCondHandle = map[define.QueryFilterRuleOperator]esCondHandleFn{
	RuleOperatorEQ:      esEq,
	RuleOperatorIN:      esIn,
	RuleOperatorBetween: esBetween,
	RuleOperatorLT:      esLt,
	RuleOperatorLTE:     esLte,
	RuleOperatorGT:      esGt,
	RuleOperatorGTE:     esGte,
	RuleOperatorLike:    esWildcard,
	RuleOperatorPrefix:  esWildcardPrefix,
	RuleOperatorSuffix:  esWildcardSuffix,
}

func esEq(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	return ges.Term(field, values[0]), nil
}

func esIn(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	return ges.Terms(field, values), nil
}

func esBetween(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var int64Vals []int64
	if err := ctx.Mapper("es orm between", values, &int64Vals); err != nil {
		return nil, err
	}
	if len(int64Vals) < 2 {
		err := ctx.Error().Errorf(code.MetricDataSourceQueryFilterRoleNotErrCode, field)
		ctx.Log().Error(err.Error())
		return nil, err
	}

	return ges.Between(field, int64Vals[0], int64Vals[1]), nil

}

func esGt(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var int64Vals []int64
	if err := ctx.Mapper("es orm gt", values, &int64Vals); err != nil {
		return nil, err
	}

	return ges.Gt(field, int64Vals[0]), nil

}

func esGte(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var int64Vals []int64
	if err := ctx.Mapper("es orm gte", values, &int64Vals); err != nil {
		return nil, err
	}

	return ges.Gte(field, int64Vals[0]), nil

}

func esLt(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var int64Vals []int64
	if err := ctx.Mapper("es orm lt", values, &int64Vals); err != nil {
		return nil, err
	}

	return ges.Lt(field, int64Vals[0]), nil

}

func esLte(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var int64Vals []int64
	if err := ctx.Mapper("es orm lte", values, &int64Vals); err != nil {
		return nil, err
	}

	return ges.Lte(field, int64Vals[0]), nil

}

func esWildcard(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var strs []string
	if err := ctx.Mapper("es orm wildcard", values, &strs); err != nil {
		return nil, err
	}
	return ges.Wildcard(field, "*"+strs[0]+"*"), nil
}

func esWildcardPrefix(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var strs []string
	if err := ctx.Mapper("es orm wildcard", values, &strs); err != nil {
		return nil, err
	}
	return ges.Wildcard(field, strs[0]+"*"), nil
}

func esWildcardSuffix(ctx context.Context, field string, values []interface{}) (ges.Filter, errors.Error) {
	var strs []string
	if err := ctx.Mapper("es orm wildcard", values, &strs); err != nil {
		return nil, err
	}
	return ges.Wildcard(field, "*"+strs[0]), nil
}
