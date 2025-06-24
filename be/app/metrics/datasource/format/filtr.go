package format

import (
	"github.com/rentiansheng/dashboard/app/metrics/datasource/tools"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"strings"
)

// AdjustQueryInput
func AdjustQueryInput(ctx context.Context, dataName string, fields []define.MetricDataSourceMeta,
	input *define.QueryReq, isDetail bool) errors.Error {
	return adjustInputHelper(ctx, fields, input, isDetail)
}

func adjustInputHelper(ctx context.Context, fields []define.MetricDataSourceMeta, input *define.QueryReq, isDetail bool) (err errors.Error) {
	nestedFieldRela := make(map[string][]string, 0)
	for _, f := range fields {
		if strings.TrimSpace(f.Nested) != "" {
			nestedFieldRela[f.FieldName] = []string{f.Nested}
		}
	}
	for indexIdx, index := range input.Indexes {
		dateRange := input.Date
		if !isDetail && input.Typ != define.ChartTypePie {
			// 根据周期调整时间范围
			dateRange, err = tools.BuildTimeRangeCycleBorder(ctx, index.Output.Cycle, input.Date)
			if err != nil {
				ctx.Log().ErrorJSON("start cycle error. data: %s, err: %s", input, err)
				return err
			}

		}
		idxTimeFilter := define.QueryFilterRule{
			Field:    index.Output.TimeField,
			Operator: define.QueryFilterRuleBetween,
			Values:   []interface{}{tools.SecondToMillisecond(dateRange.Start), tools.SecondToMillisecond(dateRange.End)},
		}

		for filterIdx, filter := range index.Filters {
			// 等待填充字段转换方法
			_, _ = filterIdx, filter
			// 根据部门限定查询条件
			newRules := make([]define.QueryFilterRule, 0, len(input.Indexes[indexIdx].Filters[filterIdx].Rules)+1)
			for _, rule := range input.Indexes[indexIdx].Filters[filterIdx].Rules {
				if nesteds, ok := nestedFieldRela[rule.Field]; ok {
					rule.Nesteds = nesteds
					newRules = append(newRules, rule)
				} else {
					newRules = append(newRules, rule)
				}
			}
			input.Indexes[indexIdx].Filters[filterIdx].Rules = append(newRules, idxTimeFilter)

		}
		if len(input.Indexes[indexIdx].Filters) == 0 {
			input.Indexes[indexIdx].Filters = append(input.Indexes[indexIdx].Filters, define.QueryFilter{Rules: []define.QueryFilterRule{idxTimeFilter}})
		}

		for aggFilterIdx, filters := range index.Output.Aggregator.Filters {
			for filterIdx, filter := range filters.Terms {
				if nesteds, ok := nestedFieldRela[filter.Field]; ok {
					filters.Terms[filterIdx].Nesteds = nesteds
				}
			}
			for filterIdx, filter := range filters.Range {
				if nesteds, ok := nestedFieldRela[filter.Field]; ok {
					filters.Range[filterIdx].Nesteds = nesteds
				}
			}
			index.Output.Aggregator.Filters[aggFilterIdx] = filters
			if nesteds, ok := nestedFieldRela[filters.ExtraAggField]; ok {
				if len(nesteds) > 0 {
					index.Output.Aggregator.Filters[aggFilterIdx].ExtraAggFieldNested = nesteds[0]
				}
			}
		}

	}

	return nil
}
