package query

import (
	"fmt"
	"time"

	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"github.com/rentiansheng/ges"
)

func GetESAggBucketsV2Helper(ctx context.Context, esClient ges.Client, aggParam define.QueryOutputAggregator) (map[string]map[string]map[string]float64, errors.Error) {

	esResults := make(map[string]interface{}, 0)
	if _, err := esClient.Search(ctx, &esResults); err != nil {
		return nil, ctx.Error().Errorf(code.ESDALSearchErrCode, err)
	}

	return GetESAggBucketsV2HelperSeriesMetric(ctx, esResults, aggParam)

}

func GetESAggBucketsV2HelperSeriesMetric(ctx context.Context, esResults map[string]interface{}, aggParam define.QueryOutputAggregator) (map[string]map[string]map[string]float64, errors.Error) {

	results := make(map[string]map[string]map[string]float64, 0)

	aggHasFilterCondRela := make(map[string]bool, len(aggParam.Filters))
	for _, agg := range aggParam.Filters {
		aggHasFilterCondRela[agg.Name] = len(agg.Terms) > 0 || len(agg.Range) > 0
	}
	for aggName, aggResultI := range esResults {

		aggNameHasFilter := aggHasFilterCondRela[aggName]
		if aggParam.IsTopN && !aggParam.IsDataHistogram {
			if aggNameHasFilter {
				aggResult := make(map[string]map[string]AggResultFiltersTopNSeries, 0)
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2AggResultFiltersTopNSeries(ctx, aggResult)

			} else {
				aggResult := DistinctBuckets{}
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2AggResultTopN(ctx, aggResult)
			}

		}
		if aggParam.IsDataHistogram && !aggParam.IsTopN {
			if aggNameHasFilter {
				aggResult := make(map[string]map[string]AggResultFiltersDateHistogramSeries, 0)
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2AggResultFiltersDateHistogramSeries(ctx, aggResult)
			} else {
				aggResult := DateHistogramBuckets{}
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2DateHistogramSeries(ctx, aggResult)
			}
		}

		if aggParam.IsDataHistogram && aggParam.IsTopN {
			if aggNameHasFilter {
				aggResult := make(map[string]map[string]AggResultFilterTopNDateHistogramSeries, 0)
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2AggResultFilterTopNDateHistogramSeries(ctx, aggResult)
			} else {
				aggResult := AggResultTopNDateHistogramSeries{}
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2DateHistogramTopNSeries(ctx, aggResult)
			}
		}
		// 出来 is_date_histogram 和 is_top_n 都是 false 的情况
		if !aggParam.IsDataHistogram && !aggParam.IsTopN {
			if aggNameHasFilter {
				aggResult := make(map[string]map[string]AggResultOnlyFilterSeries, 0)
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2AggResultOnlyFilterSeries(ctx, aggResult)
			} else {
				aggResult := AggResultOnlyExtra{}
				if err := ctx.Mapper("agg info convert", aggResultI, &aggResult); err != nil {
					return nil, err
				}
				results[aggName] = GetESAggBucketsV2AggResultOnlyExtraSeries(ctx, aggResult)
			}

		}
	}
	return results, nil
}

func GetESAggBucketsV2AggResultFiltersTopNSeries(ctx context.Context,
	input map[string]map[string]AggResultFiltersTopNSeries) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)
	for _, seriesItem := range input {
		for _, bucket := range seriesItem {

			if _, ok := results[define.AggFilterDefaultValueKey]; !ok {
				results[define.AggFilterDefaultValueKey] = make(map[string]float64, 0)
			}
			for _, item := range bucket.AresTop.Buckets {
				if item.NestedSubAgg != nil {
					if item.NestedSubAgg.SubAgg.Value != nil {
						results[define.AggFilterDefaultValueKey][fmt.Sprintf("%v", item.Key)] = *item.NestedSubAgg.SubAgg.Value
					}
				} else if item.SubAgg != nil {
					if item.SubAgg.Value != nil {
						results[define.AggFilterDefaultValueKey][fmt.Sprintf("%v", item.Key)] = *item.SubAgg.Value
					}
				} else {
					results[define.AggFilterDefaultValueKey][fmt.Sprintf("%v", item.Key)] = float64(item.DocCount)
				}
			}

		}
	}

	return results
}

func GetESAggBucketsV2AggResultTopN(ctx context.Context, input DistinctBuckets) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)

	for _, b := range input.Buckets {

		if b.NestedSubAgg != nil {
			if b.NestedSubAgg.SubAgg.Value != nil {
				results[fmt.Sprintf("%v", b.Key)] = map[string]float64{define.AggFilterDefaultValueKey: *b.NestedSubAgg.SubAgg.Value}

			}
		} else if b.SubAgg != nil {
			if b.SubAgg.Value != nil {
				results[fmt.Sprintf("%v", b.Key)] = map[string]float64{define.AggFilterDefaultValueKey: *b.SubAgg.Value}

			}
		} else {
			results[fmt.Sprintf("%v", b.Key)] = map[string]float64{define.AggFilterDefaultValueKey: float64(b.DocCount)}
		}
	}

	return results
}

func GetESAggBucketsV2AggResultFiltersDateHistogramSeries(ctx context.Context,
	input map[string]map[string]AggResultFiltersDateHistogramSeries) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)
	for _, seriesItem := range input {
		for _, bucket := range seriesItem {

			if _, ok := results[define.AggFilterDefaultValueKey]; !ok {
				results[define.AggFilterDefaultValueKey] = make(map[string]float64, 0)
			}
			for _, item := range bucket.AresDate.Buckets {
				strKey := time.Unix(item.Key/1000, 0).Format("2006-01-02")
				if item.NestedSubAgg != nil {
					if item.NestedSubAgg.SubAgg.Value != nil {
						results[define.AggFilterDefaultValueKey][strKey] = *item.NestedSubAgg.SubAgg.Value
					}
				} else if item.SubAgg != nil {
					if item.SubAgg.Value != nil {
						results[define.AggFilterDefaultValueKey][strKey] = *item.SubAgg.Value
					}
				} else {
					results[define.AggFilterDefaultValueKey][strKey] = float64(item.DocCount)
				}
			}
		}
	}

	return results
}

func GetESAggBucketsV2DateHistogramSeries(ctx context.Context, input DateHistogramBuckets) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)
	for _, b := range input.Buckets {

		if _, ok := results[define.AggFilterDefaultValueKey]; !ok {
			results[define.AggFilterDefaultValueKey] = make(map[string]float64, 0)
		}
		strKey := time.Unix(b.Key/1000, 0).Format("2006-01-02")
		if b.NestedSubAgg != nil {
			if b.NestedSubAgg.SubAgg.Value != nil {
				results[define.AggFilterDefaultValueKey][strKey] = *b.NestedSubAgg.SubAgg.Value
			}
		} else if b.SubAgg != nil {
			if b.SubAgg.Value != nil {
				results[define.AggFilterDefaultValueKey][strKey] = *b.SubAgg.Value

			}
		} else {
			results[define.AggFilterDefaultValueKey][strKey] = float64(b.DocCount)
		}
	}

	return results
}

func GetESAggBucketsV2AggResultFilterTopNDateHistogramSeries(ctx context.Context,
	input map[string]map[string]AggResultFilterTopNDateHistogramSeries) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)

	for _, inputSeriesAggValueBuckets := range input {
		for _, series := range inputSeriesAggValueBuckets {

			for _, topBucket := range series.AresTop.Buckets {
				topKey := fmt.Sprintf("%v", topBucket.Key)
				if _, ok := results[topKey]; !ok {
					results[topKey] = make(map[string]float64, 0)
				}

				for _, buckets := range topBucket.AresData.Buckets {
					strKey := time.Unix(buckets.Key/1000, 0).Format("2006-01-02")
					if buckets.NestedSubAgg != nil {
						if buckets.NestedSubAgg.SubAgg.Value != nil {
							results[topKey][fmt.Sprintf("%v", strKey)] = *buckets.NestedSubAgg.SubAgg.Value
						}
					} else if buckets.SubAgg != nil {
						if buckets.SubAgg.Value != nil {
							results[topKey][fmt.Sprintf("%v", strKey)] = *buckets.SubAgg.Value
						}
					} else {
						results[topKey][fmt.Sprintf("%v", strKey)] = float64(buckets.DocCount)
					}
				}
			}
		}

	}

	return results
}

func GetESAggBucketsV2DateHistogramTopNSeries(ctx context.Context, input AggResultTopNDateHistogramSeries) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)

	for _, bucket := range input.Buckets {
		if _, ok := results[fmt.Sprintf("%v", bucket.Key)]; !ok {
			results[fmt.Sprintf("%v", bucket.Key)] = make(map[string]float64, 0)
		}
		for _, b := range bucket.AresData.Buckets {
			strKey := time.Unix(b.Key/1000, 0).Format("2006-01-02")
			if b.NestedSubAgg != nil {
				if b.NestedSubAgg.SubAgg.Value != nil {
					results[fmt.Sprintf("%v", bucket.Key)][strKey] = *(b.NestedSubAgg.SubAgg.Value)
				}
			} else if b.SubAgg != nil {
				if b.SubAgg.Value != nil {
					results[fmt.Sprintf("%v", bucket.Key)][strKey] = *b.SubAgg.Value
				}
			} else {
				results[fmt.Sprintf("%v", bucket.Key)][strKey] = float64(b.DocCount)
			}

		}
	}

	return results
}

func GetESAggBucketsV2AggResultOnlyFilterSeries(ctx context.Context, input map[string]map[string]AggResultOnlyFilterSeries) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)

	for _, seriesItem := range input {
		for name, bucket := range seriesItem {

			if _, ok := results[define.AggFilterDefaultValueKey]; !ok {
				results[define.AggFilterDefaultValueKey] = make(map[string]float64, 0)
			}
			if bucket.NestedSubAgg != nil {
				if bucket.NestedSubAgg.SubAgg.Value != nil {
					results[define.AggFilterDefaultValueKey][name] = *bucket.NestedSubAgg.SubAgg.Value
				}
			} else if bucket.SubAgg != nil {
				if bucket.SubAgg.Value != nil {
					results[define.AggFilterDefaultValueKey][name] = *bucket.SubAgg.Value
				}
			} else {
				results[define.AggFilterDefaultValueKey][name] = float64(bucket.DocCount)
			}
		}
	}

	return results
}

func GetESAggBucketsV2AggResultOnlyExtraSeries(ctx context.Context, input AggResultOnlyExtra) map[string]map[string]float64 {
	results := make(map[string]map[string]float64, 0)

	if _, ok := results[define.AggFilterDefaultValueKey]; !ok {
		results[define.AggFilterDefaultValueKey] = make(map[string]float64, 0)
	}
	if input.NestedSubAgg != nil {
		if input.NestedSubAgg.SubAgg.Value != nil {
			results[define.AggFilterDefaultValueKey][define.AggFilterDefaultValueKey] = *input.NestedSubAgg.SubAgg.Value
		}
	} else if input.SubAgg != nil {
		if input.SubAgg.Value != nil {
			results[define.AggFilterDefaultValueKey][define.AggFilterDefaultValueKey] = *input.SubAgg.Value
		}
	} else {
		results[define.AggFilterDefaultValueKey][define.AggFilterDefaultValueKey] = float64(input.DocCount)
	}

	return results
}

// ***** 		以下是待删除 		***** //

func GetESAggBuckets(ctx context.Context, aggType uint8) interface{} {
	switch aggType {
	case define.QueryOutputAggTypeDateHistogram:
		return make(map[string]ges.DateHistogramBuckets)
	case define.QueryOutputAggTypeTopN:
		return make(map[string]ges.DistinctBuckets)
	case define.QueryOutputAggTypeTerms:
		return make(map[string]ges.DocCountBuckets)
	case define.QueryOutputAggTypeHistogramTerms:
		return make(map[string]termsHistogramBucket, 0)
	case define.QueryOutputAggTypeXAxisField:
		return make(map[string]xAxisBucket, 0)
	default:

		return make(map[string]interface{}, 0)
	}
}

func getBuckets(bucketArrI []interface{}) []define.ESAggBuckets {
	results := make([]define.ESAggBuckets, 0)
	for _, bucket := range bucketArrI {
		item := define.ESAggBuckets{
			Values:    map[string]interface{}{},
			SubBucket: map[string][]define.ESAggBuckets{},
		}
		kv, ok := bucket.(map[string]interface{})
		if !ok {
			continue
		}
		for key, val := range kv {
			switch key {
			case "key":
				item.Key = val
			case "key_as_string":
				item.KeyAsStr = val.(string)
			case "doc_count":

				item.DocCount = val

			case "doc_count_error_upper_bound":
			case "sum_other_doc_count":
			default:
				if subKv, ok := val.(map[string]interface{}); ok {
					if extraKvV, ok := subKv["value"]; ok {
						item.Values[key] = extraKvV
					} else if nextBuckets, ok := subKv["buckets"].([]interface{}); ok {

						item.SubBucket[key] = getBuckets(nextBuckets)
					}

				}
			}
		}
		results = append(results, item)

	}
	return results
}

func getMetricBuckets(keys []interface{}, subBucket map[string][]define.ESAggBuckets) map[string][]esAggMetricBucket {
	results := make(map[string][]esAggMetricBucket, 0)

	for name, buckets := range subBucket {
		for _, bucket := range buckets {
			item := esAggMetricBucket{}
			itemKeys := []interface{}{}
			itemKeys = append(itemKeys, keys...)
			if bucket.Key != nil {
				itemKeys = append(keys, bucket.Key)
			}
			if len(bucket.SubBucket) > 0 {
				subMetricBuckets := getMetricBuckets(itemKeys, bucket.SubBucket)
				for name, subs := range subMetricBuckets {
					results[name] = append(results[name], subs...)
				}
			} else {
				item.Keys = itemKeys
				item.Bucket = bucket
				results[name] = append(results[name], item)
			}
		}

	}
	return results
}

func GetMetricBucket(bucketValue interface{}) map[string][]esAggMetricBucket {
	esBucketRela := map[string][]define.ESAggBuckets{}
	if val, ok := bucketValue.(map[string]ges.DateHistogramBuckets); ok {
		for aggName, buckets := range val {
			for _, b := range buckets.Buckets {
				esBucketRela[aggName] = append(esBucketRela[aggName], define.ESAggBuckets{
					Key:      b.Key,
					KeyAsStr: b.StrKey,
					DocCount: b.Count,
				})
			}
		}

	} else if val, ok := bucketValue.(map[string]ges.DistinctBuckets); ok {
		for aggName, buckets := range val {
			for _, b := range buckets.Buckets {
				esBucketRela[aggName] = append(esBucketRela[aggName], define.ESAggBuckets{
					Key:      b.Key,
					DocCount: b.DocCount,
				})
			}
		}
	} else if val, ok := bucketValue.(map[string]ges.DocCountBuckets); ok {
		for aggName, buckets := range val {
			for bName, b := range buckets.Buckets {
				esBucketRela[aggName] = append(esBucketRela[aggName], define.ESAggBuckets{
					Key:      bName,
					DocCount: b.DocCount,
				})
			}
		}

	} else if val, ok := bucketValue.(map[string]termsHistogramBucket); ok {
		for _, bucketsL1 := range val {
			for aggName, bucketsL2 := range bucketsL1.Buckets {
				if subAgg, ok := bucketsL2.(map[string]interface{}); ok {
					for _, subAggValue := range subAgg {
						if subAggValueKey, ok := subAggValue.(map[string]interface{}); ok {
							if nodeBuckets, ok := subAggValueKey["buckets"].([]interface{}); ok {
								esBucketRela[aggName] = getBuckets(nodeBuckets)
							}
						}

					}
				}

			}
		}
	} else if val, ok := bucketValue.(map[string]xAxisBucket); ok {
		// {"aggregations":{"project":{"doc_count_error_upper_bound":0,"sum_other_doc_count":534,"buckets":[{"key":"SPXFM","doc_count":3283,"dev2":{"value":31125.918976545843},"dev":{"buckets":{"dev":{"doc_count":1917,"dev":{"value":59066625}}}}}]}}}
		// L1 为 aggregations.project 值
		for _, L1 := range val {
			//  aggregations.project.buckets.doc_count

			//  aggregations.project.buckets 值
			for _, l1Bucket := range L1.Buckets {
				xAxisName := fmt.Sprintf("%v", l1Bucket["key"])

				//  有是agg name ，也可能是es 聚合内置字段eg(key, doc_count)
				for l2k, l2v := range l1Bucket {
					switch l2k {
					case "key":

					case "doc_count":

					default:
						//  aggregations.project.buckets.[dev2|dev] 等子聚合
						if l3, ok := l2v.(map[string]interface{}); ok {
							aggName := l2k

							//  aggregations.project.buckets.dev2.value
							if l2vVal, ok := l3["value"]; ok {
								if l2vVal == nil {
									l2vVal = (*float64)(nil)
								}
								esBucketRela[aggName] = append(esBucketRela[aggName], define.ESAggBuckets{
									Key:      xAxisName,
									DocCount: l2vVal,
								})
							} else if l3vBucket, ok := l3["buckets"].(map[string]interface{}); ok { //aggregations.project.buckets.dev.buckets
								//aggregations.project.buckets.dev.buckets.dev
								if l4, ok := l3vBucket[aggName].(map[string]interface{}); ok {
									//  aggregations.project.buckets.dev.buckets.dev.dev
									if l5, ok := l4[aggName].(map[string]interface{}); ok {
										if l5Val, ok := l5["value"]; ok {
											esBucketRela[aggName] = append(esBucketRela[aggName], define.ESAggBuckets{
												Key:      xAxisName,
												DocCount: l5Val,
											})
										}

									}
								}
							}

						}

					}
				}

			}
		}

	} else if val, ok := bucketValue.(map[string]interface{}); ok {
		esBucketRela = getMetricBucket(val)
	}
	return getMetricBuckets([]interface{}{}, esBucketRela)
}

func getMetricBucket(bucketValue map[string]interface{}) map[string][]define.ESAggBuckets {

	esBucketRela := map[string][]define.ESAggBuckets{}
	for name, aggItem := range bucketValue {
		if aggKv, ok := aggItem.(map[string]interface{}); ok {
			if _, ok := aggKv["buckets"]; ok {
				switch buckets := aggKv["buckets"].(type) {
				case []interface{}:
					// 判断是否需要更改
					esBucketRela[name] = getBuckets(buckets)
				}
			}

		}
	}

	return esBucketRela
}

// MetricBucketEnum 根据metric的meta信息，对bucket进行枚举,调整到对应的展示值
func MetricBucketEnum(ctx context.Context, dataSourceId uint64, metas []define.MetricDataSourceMeta, aggNameFieldRela map[string]string,
	metricBuckets map[string][]esAggMetricBucket, aggType uint8) (map[string][]define.ESAggBuckets, errors.Error) {
	metaRela := make(map[string]define.MetricDataSourceMeta, 0)
	for _, meta := range metas {
		if meta.DataSourceID == dataSourceId {
			metaRela[meta.FieldName] = meta
		}
	}
	enumFields := make([]define.MetricDataSourceMeta, 0)

	for _, field := range aggNameFieldRela {
		if meta, ok := metaRela[field]; ok {
			enumFields = append(enumFields, meta)
		}
	}
	enumRela := make(map[string]map[string]interface{}, 0)
	for _, enumField := range enumFields {
		enumFieldRela, err := DataSourceEnumField(ctx, enumField, nil)
		if err != nil {
			return nil, err
		}
		enumRela[enumField.FieldName] = enumFieldRela
	}

	newMetricBuckets := make(map[string][]define.ESAggBuckets, 0)

	for aggName, bucket := range metricBuckets {
		// aggName  是存储层面的值， top n + histogram top n 需要转换成对应的值
		newAggName := aggName
		for _, item := range bucket {
			switch aggType {
			case define.QueryOutputAggTypeXAxisField:
				// 判断agg字段是否需要做枚举转换
				newAggName = fmt.Sprintf("%v", item.Keys[0])
			case define.QueryOutputAggTypeHistogramTopN:
				// {"aggregations":{"NAME":{"doc_count_error_upper_bound":691,"sum_other_doc_count":21993,"buckets":[{"key":"SPXFM","doc_count":12540,"NAM1E":{"buckets":{"xxxxxx":{"doc_count":0,"xxxx":{"value":0}}}},"NAME":{"buckets":{"xxxxxx":{"doc_count":12540,"xxxx":{"value":60665.859998246655}}}}}]}}}
				newAggName = fmt.Sprintf("%v", item.Keys[len(item.Keys)-2])

			}
			// 判断agg字段是否需要做枚举转换
			if enumRelaField, ok := enumRela[aggNameFieldRela[aggName]]; ok {
				switch aggType {
				case define.QueryOutputAggTypeTopN:
					// bucket key 需要做枚举转换
					val := item.Bucket.Key
					if enumVal, ok := enumRelaField[fmt.Sprintf("%v", val)]; ok {
						item.Bucket.Key = enumVal
					}
				case define.QueryOutputAggTypeHistogramTopN:
					// bucket 倒数key 需要做枚举转换，是枚举字段的top，
					val := item.Keys[len(item.Keys)-2]
					if enumVal, ok := enumRela[fmt.Sprintf("%v", val)]; ok {
						newAggName = fmt.Sprintf("%v", enumVal)
					}
				case define.QueryOutputAggTypeXAxisField:
					// bucket key 需要做枚举转换
					val := item.Keys[0]
					if enumVal, ok := enumRela[fmt.Sprintf("%v", val)]; ok {
						item.Bucket.Key = enumVal
					}
					// 不改变aggName， 只是key 需要装换
					newAggName = aggName
				}
			}
			newMetricBuckets[newAggName] = append(newMetricBuckets[newAggName], item.Bucket)
		}
	}

	return newMetricBuckets, nil

}
