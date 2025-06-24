package query

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
)

type esAggMetricBucket struct {
	Keys []interface{} `json:"keys"`
	// 这里是bucket 最后一级别
	Bucket define.ESAggBuckets `json:"buckets"`
}

const aggDefaultName = "_default_"

type termsHistogramBucket struct {
	Buckets map[string]interface{} `json:"buckets"`
}

type xAxisBucket struct {
	DocCountErrorUpperBound int                      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                      `json:"sum_other_doc_count"`
	Buckets                 []map[string]interface{} `json:"buckets"`
}

// ****** 以下是es返回的结果的结构体字段预定义，可以是使用struct 来decode 结果 ****** //
const AggDistinctName = define.AggDistinctName
const AggDateHistogramName = define.AggDateHistogramName
const AggFilterSubAgg = define.AggFilterSubAgg
const AggFilterNestedName = define.AggFilterNestedName

// ****** 以上是es返回的结果的结构体字段预定义，可以是使用struct 来decode 结果 ****** //

// ****** 以下是es返回的结果的结构体，根据参数定义 ****** //

type DateHistogramBuckets struct {
	Buckets []DateHistogramBucket `json:"buckets"`
}

type DistinctBuckets struct {
	DocCountErrorUpperBound int              `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int              `json:"sum_other_doc_count"`
	Buckets                 []DistinctBucket `json:"buckets"`
}

type AggResultTopNDateHistogramSeries struct {
	Buckets []struct {
		AresData struct {
			Buckets []DateHistogramBucket `json:"buckets"`
		} `json:"__dashboard_date__"`
		DocCount int         `json:"doc_count"`
		Key      interface{} `json:"key"`
	} `json:"buckets"`
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
}

type AggResultFiltersTopNSeries struct {
	DocCount int `json:"doc_count"`
	AresTop  struct {
		DocCountErrorUpperBound int              `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int              `json:"sum_other_doc_count"`
		Buckets                 []DistinctBucket `json:"buckets"`
	} `json:"__dashboard_top__"`
}

type AggResultFiltersTopNSeriesBak struct {
	DocCount int `json:"doc_count"`
	AresTop  struct {
		DocCountErrorUpperBound int              `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int              `json:"sum_other_doc_count"`
		Buckets                 []DistinctBucket `json:"buckets"`
	} `json:"__dashboard_top__"`
}

type AggResultFiltersDateHistogramSeries struct {
	DocCount int `json:"doc_count"`
	AresDate struct {
		Buckets []DateHistogramBucket `json:"buckets"`
	} `json:"__dashboard_date__"`
}

type AggResultFilterTopNDateHistogramSeries struct {
	DocCount int `json:"doc_count"`
	AresTop  struct {
		DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int `json:"sum_other_doc_count"`
		Buckets                 []struct {
			Key      interface{} `json:"key"`
			DocCount int         `json:"doc_count"`
			AresData struct {
				Buckets []DateHistogramBucket `json:"buckets"`
			} `json:"__dashboard_date__"`
		} `json:"buckets"`
	} `json:"__dashboard_top__"`
}

type DateHistogramBucket struct {
	DocCount    int    `json:"doc_count"`
	Key         int64  `json:"key"`
	KeyAsString string `json:"key_as_string"`
	// filter 有 sum,avg 时候才有
	SubAgg *AggResultFilterSubAggValue `json:"__dashboard_sub_agg__"`
	//   filter 有 sum,avg  且字段是nested 有效
	NestedSubAgg *AggResultFilterNestedSubAggValue `json:"__dashboard_nested__"`
}

type DistinctBucket struct {
	Key      interface{} `json:"key"`
	DocCount int         `json:"doc_count"`
	// filter 有 sum,avg 时候才有
	SubAgg *AggResultFilterSubAggValue `json:"__dashboard_sub_agg__"`
	//   filter 有 sum,avg  且字段是nested 有效
	NestedSubAgg *AggResultFilterNestedSubAggValue `json:"__dashboard_nested__"`
}

// filter 有 sum,avg 时候才有
type AggResultFilterSubAggValue struct {
	Value *float64 `json:"value"`
}

type AggResultOnlyFilterSeries AggResultOnlyExtra

type AggResultOnlyExtra struct {
	DocCount int `json:"doc_count"`
	// filter 有 sum,avg 时候才有
	SubAgg *AggResultFilterSubAggValue `json:"__dashboard_sub_agg__"`
	//   filter 有 sum,avg  且字段是nested 有效
	NestedSubAgg *AggResultFilterNestedSubAggValue `json:"__dashboard_nested__"`
}

// filter 有 sum,avg  且字段是nested 有效
type AggResultFilterNestedSubAggValue struct {
	DocCount int                        `json:"doc_count"`
	SubAgg   AggResultFilterSubAggValue `json:"__dashboard_sub_agg__"`
}

// ****** 以上是es返回的结果的结构体，根据参数定义 ****** //
