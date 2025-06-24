package define

/***************************
    @author: tiansheng.ren
    @date: 7/24/23
    @desc:

***************************/

type ESAggBuckets struct {
	Key       interface{}               `json:"key"`
	KeyAsStr  string                    `json:"key_as_string"`
	DocCount  interface{}               `json:"doc_count"`
	Values    map[string]interface{}    `json:"values"`
	SubBucket map[string][]ESAggBuckets `json:"sub_bucket"`
}

// ****** 以下是es返回的结果的结构体字段预定义，可以是使用struct 来decode 结果 ****** //

const AggDistinctName = "__dashboard_top__"
const AggDateHistogramName = "__dashboard_date__"
const AggFilterSubAgg = "__dashboard_sub_agg__"
const AggFilterDefaultValueKey = "__dashboard_default_value__"
const AggFilterNestedName = "__dashboard_nested__"

// ****** 以上是es返回的结果的结构体字段预定义，可以是使用struct 来decode 结果 ****** //
