package query

import (
	"fmt"
	"github.com/rentiansheng/dashboard/app/metrics/datasource/filter/es"
	"github.com/rentiansheng/dashboard/app/metrics/datasource/format"
	"github.com/rentiansheng/dashboard/app/metrics/datasource/tools"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"github.com/rentiansheng/ges"
)

func logESEngine(ctx context.Context, engine ges.Client) {
	if engine == nil {
		return
	}
	// TODO: add log
}

func SearchES(ctx context.Context, info define.MetricDataSource, query *define.QueryReq) ([]map[string]interface{}, uint64, errors.Error) {

	engine, err := buildESSearchClient(ctx, info, query)
	if err != nil {
		return nil, 0, err
	}

	engine = engine.Limit(uint64(query.Page.Offset()), uint64(query.Page.PageSize)).Fields("*")
	results := make([]map[string]interface{}, 0)
	total, gErr := engine.Search(ctx, &results)
	if gErr != nil {
		ctx.Log().ErrorJSON("es search error. cond: %s, err: %s", query, gErr)
		return nil, 0, ctx.Error().Errorf(code.ESSearchErrCode, gErr.Error())
	}
	return results, total, nil
}

func SearchESMetric(ctx context.Context, info define.MetricDataSource, metas []define.MetricDataSourceMeta, query *define.QueryReq) (*define.Chart, errors.Error) {

	for _, index := range query.Indexes {
		engine, aggNameFieldRela, err := buildESSearchMetricClient(ctx, info, query.GroupKeyId, index)
		if err != nil {
			return nil, err
		}
		defer logESEngine(ctx, engine)
		aggSeries, err := GetESAggBucketsV2Helper(ctx, engine, index.Output.Aggregator)
		if err != nil {
			ctx.Log().ErrorJSON("es search metric error. cond: %s, err: %s", query, err)
			return nil, err
		}
		// 枚举转换
		aggSeries, err = SearchESMetricEnum(ctx, metas, aggNameFieldRela, aggSeries)
		if err != nil {
			return nil, err
		}
		if query.Typ == define.ChartTypePie {
			return format.ESPieChartV2(ctx, aggSeries)
		} else {
			if index.Output.Aggregator.IsDataHistogram {
				return format.ESBarLineDateHistogramChartV2(ctx, int8(index.Output.Cycle), query.Typ, query.Date, aggSeries)
			}
			return format.ESBarLineChartV2(ctx, query.Typ, aggSeries)

		}

	}
	// 这里不存在这个情况， 输入参数，在上层已经校验过了
	return nil, nil
}

func SearchESMetricEnum(ctx context.Context, metas []define.MetricDataSourceMeta, aggNameFieldRela map[string]string,
	aggSeries map[string]map[string]map[string]float64) (map[string]map[string]map[string]float64, errors.Error) {
	// 找到需要处理枚举的字段
	enumRela := make(map[string]struct{}, len(aggNameFieldRela))
	for _, fieldName := range aggNameFieldRela {
		enumRela[fieldName] = struct{}{}
	}
	enumKv := make(map[string]map[string]interface{}, len(aggNameFieldRela))
	for _, meta := range metas {
		// 不需要转换的字段，放弃
		if _, ok := enumRela[meta.FieldName]; !ok {
			continue
		}
		enmuInfo, err := DataSourceEnumField(ctx, meta, nil)
		if err != nil {
			ctx.Log().ErrorJSON("metric top key value enum convert error. field: %s", meta)
			return nil, err
		}
		enumKv[meta.FieldName] = enmuInfo
	}

	for aggName, fieldName := range aggNameFieldRela {
		enumInfo := enumKv[fieldName]
		if len(enumInfo) == 0 {
			continue
		}

		metricVal := aggSeries[aggName]
		newMetricVal := make(map[string]map[string]float64, len(metricVal))
		for k, v := range metricVal {
			newK := k
			if newEnumV, ok := enumInfo[k]; ok {
				newK = fmt.Sprintf("%v", newEnumV)
			}
			newMetricVal[newK] = v
		}
		aggSeries[aggName] = newMetricVal
	}
	return aggSeries, nil

}

func buildESSearchClient(ctx context.Context, info define.MetricDataSource, query *define.QueryReq) (ges.Client, errors.Error) {
	engine := ges.ES()
	for _, index := range query.Indexes {
		engine, err := buildESSearchClientIndex(ctx, info, query.GroupKeyId, index.Filters)
		if err != nil {
			return nil, err
		}
		engine = engine.IndexName(info.DataName)
		for _, sort := range index.Output.Sorts {
			engine = engine.OrderBy(sort, true)
		}
		return engine, nil

	}

	return engine, nil
}

func buildESSearchMetricClient(ctx context.Context, info define.MetricDataSource, id uint64, index define.QueryIndexes) (ges.Client, map[string]string, errors.Error) {

	engine, err := buildESSearchClientIndex(ctx, info, id, index.Filters)
	if err != nil {
		return nil, nil, err
	}

	aggs, aggNameFieldRela, err := buildESSearchAggV2(ctx, index.Output)
	if err != nil {
		return nil, nil, err
	}
	engine = engine.Agg(aggs...)
	return engine, aggNameFieldRela, nil

}

func buildESSearchClientIndex(ctx context.Context, info define.MetricDataSource, id uint64, filters []define.QueryFilter) (ges.Client, errors.Error) {
	engine := ges.ES()
	handler, err := es.GetIndexHandle(ctx, info)
	if err != nil {
		return nil, err
	}
	for _, filter := range filters {
		condsItem, notCondsItem, err := es.QueryParams(ctx, handler, filter.Rules)
		if err != nil {
			ctx.Log().ErrorJSON("es query parameter error. Rules: %s, err: %s", filter.Rules, err.Error())
			return nil, err
		}
		engine = engine.Where(condsItem...).Not(notCondsItem...)
	}

	extra, err := handler.ExtraAdjust(ctx, info, id)
	if err != nil {
		ctx.Log().ErrorJSON("es extra filter error. err: %s", err.Error())
		return nil, err
	}
	return engine.IndexName(info.DataName).Where(extra), nil
}

func buildESSearchAggV2(ctx context.Context, output define.QueryOutput) ([]ges.Agg, map[string]string, errors.Error) {
	aggParam := output.Aggregator

	buildSubAgg := func(agg ges.Agg) ges.Agg {

		var subAgg ges.Agg = nil
		if aggParam.IsDataHistogram {
			subAgg = tools.DataHistogramCycle(output.Cycle, AggDateHistogramName, output.TimeField)
			if agg != nil {
				subAgg = subAgg.Aggs(agg)
			}
		} else if agg != nil {
			subAgg = agg
		}
		if aggParam.IsTopN {
			if aggParam.TopN.Num == 0 {
				aggParam.TopN.Num = 10
			}
			if subAgg == nil {
				subAgg = ges.AggDistinctName(AggDistinctName, aggParam.TopN.Field, aggParam.TopN.Num)
			} else {
				subAgg = ges.AggDistinctName(AggDistinctName, aggParam.TopN.Field, aggParam.TopN.Num).Aggs(subAgg)
			}
		}
		return subAgg
	}

	buildSubAggNoFilter := func(name string, agg ges.Agg) ges.Agg {
		return buildSubAgg(agg).Name(name)
	}

	aggNameFieldRela := make(map[string]string, 0)

	Aggs := make([]ges.Agg, 0, len(aggParam.XAxis.Filters))
	for _, f := range aggParam.Filters {
		//  判断f.Name 判断agg 聚合使用字段是否需要转换
		if aggParam.IsTopN {
			// 如果有top, field value 有可能需要做枚举转换
			aggNameFieldRela[f.Name] = aggParam.TopN.Field
		}

		hasFilter := false
		filter := ges.AggFilterName(f.Name)
		for _, term := range f.Terms {
			hasFilter = true
			filter = filter.Terms(term.Field, term.Value...)
		}
		for _, term := range f.Range {
			hasFilter = true
			rf := ges.NewRangeFilter()
			if term.From != "" {
				rf.Gte(term.From)
			}
			if term.To != "" {
				rf.Lte(term.To)
			}
			filter = filter.Range(term.Field, rf)
		}
		switch f.AggType {
		case define.QueryOutputAggregatorXAxisFiltersAggTypeSum:
			agg := ges.AggSum(AggFilterSubAgg, f.ExtraAggField)
			if f.ExtraAggFieldNested != "" {
				agg = ges.AggNested(AggFilterNestedName, f.ExtraAggFieldNested).Aggs(agg)
			}
			if hasFilter {
				Aggs = append(Aggs, ges.AggFilters(f.Name, buildSubAgg(agg), filter))
			} else {
				Aggs = append(Aggs, buildSubAggNoFilter(f.Name, agg))
			}
		case define.QueryOutputAggregatorXAxisFiltersAggTypeAvg:
			agg := ges.AggAvg(AggFilterSubAgg, f.ExtraAggField)
			if f.ExtraAggFieldNested != "" {
				agg = ges.AggNested(AggFilterNestedName, f.ExtraAggFieldNested).Aggs(agg)
			}
			if hasFilter {
				Aggs = append(Aggs, ges.AggFilters(f.Name, buildSubAgg(agg), filter))
			} else {
				Aggs = append(Aggs, buildSubAggNoFilter(f.Name, agg))
			}
		case define.QueryOutputAggregatorXAxisFiltersAggTypeCount, define.QueryOutputAggregatorXAxisFiltersAggTypeDefault:
			if hasFilter {
				subAgg := ges.AggFilters(f.Name, nil, filter)
				subAgg = subAgg.Aggs(buildSubAgg(nil))
				Aggs = append(Aggs, subAgg)
			} else {
				Aggs = append(Aggs, buildSubAggNoFilter(f.Name, nil))
			}

		}

	}
	// 有可能没有参数，只需要topN，或者只需要时间聚合，或者 topn + 时间聚合
	if len(Aggs) == 0 {
		subAgg := buildSubAgg(nil)
		if subAgg == nil {
			return Aggs, aggNameFieldRela, ctx.Error().Errorf(code.JSONDecodeErrCode, "Aggs is empty")
		}
		if aggParam.IsTopN {
			name, _ := subAgg.Result()
			aggNameFieldRela[name] = aggParam.TopN.Field
		}
		return []ges.Agg{subAgg}, aggNameFieldRela, nil
	}

	return Aggs, aggNameFieldRela, nil
}
