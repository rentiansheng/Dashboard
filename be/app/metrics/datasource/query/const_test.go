package query

import (
	"encoding/json"
	"testing"
)

func TestDecode(t *testing.T) {

}

func TestAggResultFilterTopNDateHistogramSeriesDecode(t *testing.T) {
	esResult := `{"series1":{"buckets":{"series1":{"doc_count":225794,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":3,"doc_count":8,"__dashboard_data__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":8}]}}]}}}}}`
	res := make(map[string]map[string]map[string]AggResultFilterTopNDateHistogramSeries, 0)
	if err := json.Unmarshal([]byte(esResult), &res); err != nil {
		t.Error(err)
	}

}

func TestAggResultFiltersDateHistogramSeriesEncode(t *testing.T) {
	esResult := `{"series1":{"buckets":{"series1":{"doc_count":225794,"__dashboard_date__":{"buckets":[{"key_as_string":"2022-01-01","key":1640966400000,"doc_count":200},{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":225594}]}}}}}`
	res := make(map[string]map[string]map[string]AggResultFiltersDateHistogramSeries, 0)
	if err := json.Unmarshal([]byte(esResult), &res); err != nil {
		t.Error(err)
	}

}

func TestAggResultFiltersTopNSeriesDecode(t *testing.T) {
	esResult := `{"series1":{"buckets":{"series1":{"doc_count":225794,"__dashboard_top__":{"doc_count_error_upper_bound":2112,"sum_other_doc_count":168514,"buckets":[{"key":"x","doc_count":3567}]}}}}}`
	res := make(map[string]map[string]map[string]AggResultFiltersTopNSeries, 0)
	if err := json.Unmarshal([]byte(esResult), &res); err != nil {
		t.Error(err)
	}
	t.Log(res)
}
