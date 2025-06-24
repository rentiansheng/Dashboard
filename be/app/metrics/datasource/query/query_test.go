package query

import (
	"encoding/json"
	"fmt"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"strings"
	"testing"
)

func TestBuildESSearchAgg(t *testing.T) {
	testSuits := []struct {
		name       string
		inputStr   string
		output     string
		nestedRela map[string]string
	}{
		{
			name:     "date histogram only",
			inputStr: `{"is_date_histogram": true}`,
			output:   `{"__dashboard_date__":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}}`,
		},
		{
			name:     "date histogram and filter",
			inputStr: `{"is_date_histogram":true,"filters":[{"name":"series1","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "date histogram and filter and sum",
			inputStr: `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "date histogram and filter and avg",
			inputStr: `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "date histogram and filter and sum and nested",
			inputStr: `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "date histogram and filter and avg and nested",
			inputStr: `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":1,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "top key only",
			inputStr: `{ "is_top_n": true, "top_n":{ "field": "jira_status"}}`,
			output:   `{"__dashboard_top__":{"terms":{"field":"jira_status","size":10}}}`,
		},
		{
			name:     "top key and filter",
			inputStr: `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "top key and filter and sum",
			inputStr: `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "top key and filter and avg",
			inputStr: `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "top key and filter and sum and nested",
			inputStr: `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "top key and filter and avg and nested",
			inputStr: `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "only filter",
			inputStr: `{"filters":[{"name":"series1","agg_type":0,"terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "only filter and sum",
			inputStr: `{"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "filter and avg",
			inputStr: `{"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "filter and sum and nested",
			inputStr: `{"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "filter and avg and nested",
			inputStr: `{"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "date histogram and top key",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"}}`,
			output:   `{"__dashboard_top__":{"aggs":{"__dashboard_date__":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"priority","size":10}}}`,
		},
		{
			name:     "date histogram and top key and filter",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_date__":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		//
		{
			name:     "date histogram and top key and filter and sum",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "date histogram and top key and filter and avg",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
		},
		{
			name:     "date histogram and top key and filter and sum and nested",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "date histogram and top key and filter and avg and nested",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_top__":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_nested__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"sub_task_type_story_points.program_development"}}},"nested":{"path":"sub_task_type_story_points"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"priority","size":10}}},"filters":{"filters":{"series1":{"terms":{"jira_status":["Done"]}}}}}}`,
			nestedRela: map[string]string{
				"sub_task_type_story_points.program_development": "sub_task_type_story_points",
			},
		},
		{
			name:     "date histogram and only count",
			inputStr: `{"is_top_n": false,"is_date_histogram":true,"filters":[{"name":"series1"}]}`,
			output:   `{"series1":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}}`,
		},
		{
			name:     "date histogram and only sum",
			inputStr: `{"is_top_n": false,"is_date_histogram":true,"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day"}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}}`,
		},
		{
			name:     "data histogram and only avg",
			inputStr: `{"is_top_n": false,"is_date_histogram":true,"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day"}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}}`,
		},
		//{
		//	name: "date histogram and only sum and nested",
		//},
		//{
		//	name: "data histogram and only avg and nested",
		//},
		{
			name:     "top k and only count",
			inputStr: `{"is_top_n": true,"is_date_histogram":false, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":0}]}`,
			output:   `{"series1":{"terms":{"field":"jira_status","size":10}}}`,
		},
		{
			name:     "top k and only sum",
			inputStr: `{"is_top_n": true,"is_date_histogram":false, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day"}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"terms":{"field":"jira_status","size":10}}}`,
		},
		{
			name:     "top k and only avg",
			inputStr: `{"is_top_n": true,"is_date_histogram":false, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day"}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"terms":{"field":"jira_status","size":10}}}`,
		},
		//{
		//	name: "top k and only sum and nested",
		//},
		//{
		//	name: "top k and only avg and nested",
		//},
		{
			name:     "date histogram and top k and only count",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":0}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"jira_status","size":10}}}`,
		},
		{
			name:     "date histogram and top k and only sum",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day"}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_sub_agg__":{"sum":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"jira_status","size":10}}}`,
		},
		{
			name:     "date histogram and top k and only avg",
			inputStr: `{"is_top_n": true,"is_date_histogram":true, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day"}]}`,
			output:   `{"series1":{"aggs":{"__dashboard_date__":{"aggs":{"__dashboard_sub_agg__":{"avg":{"field":"total_cost_day"}}},"date_histogram":{"field":"ctime","calendar_interval":"year","format":"yyyy-MM-dd","time_zone":"+08:00"}}},"terms":{"field":"jira_status","size":10}}}`,
		},
		//{
		//	name: "date histogram and top k and only sum and nested",
		//},
		//{
		//	name: "date histogram and top k and only avg and nested",
		//},
	}

	allAggKv := make(map[string]interface{}, 0)

	for idx, test := range testSuits {

		input := define.QueryOutput{
			Cycle:     1,
			TimeField: "ctime",
		}
		if err := json.Unmarshal([]byte(test.inputStr), &input.Aggregator); err != nil {
			t.Errorf("test %v name: %v failed, err: %v", idx, test.name, err)
			continue
		}
		for idx, filter := range input.Aggregator.Filters {
			if nested, ok := test.nestedRela[filter.ExtraAggField]; ok {
				input.Aggregator.Filters[idx].ExtraAggFieldNested = nested
			}
		}

		aggs, _, meErr := buildESSearchAggV2(nil, input)
		if meErr != nil {
			t.Errorf("test %v name: %v failed, err: %v", idx, test.name, meErr)
			continue
		}
		aggKv := make(map[string]interface{}, len(aggs))
		for _, agg := range aggs {
			k, v := agg.Result()
			aggKv[k] = v
			key := fmt.Sprintf("%v_%v", 100+idx, k)
			allAggKv[key] = v
		}

		aggBytes, err := json.Marshal(aggKv)
		if err != nil {
			t.Errorf("test %v failed, err: %v", idx, err)
			continue
		}
		if string(aggBytes) != test.output {
			t.Errorf("test %v name %v, failed, expect: %v, got: %v", idx, test.name, test.output, string(aggBytes))
			continue
		}

	}

	aggBytes, err := json.Marshal(allAggKv)
	if err != nil {
		t.Errorf("test failed, err: %v", err)
	}
	fmt.Println(string(aggBytes))

}

func TestESAggConvertToMetric(t *testing.T) {

	tests := []struct {
		name      string
		inputStr  string
		aggResult string
		metricStr string
	}{
		{
			name:      "date histogram only",
			inputStr:  `{"is_date_histogram": true}`,
			aggResult: `{"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":21219}]}}`,
			metricStr: `{"__dashboard_date__":{"__dashboard_default_value__":{"2023-01-01":21219}}}`,
		},
		{
			name:      "date histogram and filter",
			inputStr:  `{"is_date_histogram":true,"filters":[{"name":"series1","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":10568}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2023-01-01":10568}}}`,
		},
		{
			name:      "date histogram and filter and sum",
			inputStr:  `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":10568,"__dashboard_sub_agg__":{"value":278094.4500413425}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2023-01-01":278094.4500413425}}}`,
		},
		{
			name:      "date histogram and filter and avg",
			inputStr:  `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":10568,"__dashboard_sub_agg__":{"value":26.314766279460876}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2023-01-01":26.314766279460876}}}`,
		},
		{
			name:      "date histogram and filter and sum and nested",
			inputStr:  `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":10568,"__dashboard_nested__":{"doc_count":8407,"__dashboard_sub_agg__":{"value":0.023200904006437328}}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2023-01-01":0.023200904006437328}}}`,
		},
		{
			name:      "date histogram and filter and avg and nested",
			inputStr:  `{"is_date_histogram":true,"filters":[{"name":"series1","agg_type":1,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":10568}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2023-01-01":10568}}}`,
		},
		{
			name:      "top key only",
			inputStr:  `{ "is_top_n": true, "top_n":{ "field": "jira_status"}}`,
			aggResult: `{"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":203,"buckets":[{"key":"Done","doc_count":10568},{"key":"Waiting","doc_count":5045},{"key":"Closed","doc_count":3083},{"key":"Developing","doc_count":658},{"key":"Staging","doc_count":558},{"key":"PRD","doc_count":386},{"key":"UAT","doc_count":233},{"key":"Icebox","doc_count":189},{"key":"Testing","doc_count":167},{"key":"Delivering","doc_count":129}]}}`,
			metricStr: `{"__dashboard_top__":{"Closed":{"__dashboard_default_value__":3083},"Delivering":{"__dashboard_default_value__":129},"Developing":{"__dashboard_default_value__":658},"Done":{"__dashboard_default_value__":10568},"Icebox":{"__dashboard_default_value__":189},"PRD":{"__dashboard_default_value__":386},"Staging":{"__dashboard_default_value__":558},"Testing":{"__dashboard_default_value__":167},"UAT":{"__dashboard_default_value__":233},"Waiting":{"__dashboard_default_value__":5045}}}`,
		},
		{
			name:      "top key and filter",
			inputStr:  `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135},{"key":"High","doc_count":646},{"key":"Medium","doc_count":313},{"key":"Highest","doc_count":250},{"key":"","doc_count":29},{"key":"Lowest","doc_count":3},{"key":"Blocker","doc_count":2}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"":29,"Blocker":2,"High":646,"Highest":250,"Low":4135,"Lowest":3,"Medium":313}}}`,
		},
		{
			name:      "top key and filter and sum",
			inputStr:  `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_sub_agg__":{"value":91525.1500494685}},{"key":"High","doc_count":646,"__dashboard_sub_agg__":{"value":25510.579972507432}},{"key":"Medium","doc_count":313,"__dashboard_sub_agg__":{"value":9151.430022606626}},{"key":"Highest","doc_count":250,"__dashboard_sub_agg__":{"value":8571.340042842552}},{"key":"","doc_count":29,"__dashboard_sub_agg__":{"value":680.4299903009087}},{"key":"Lowest","doc_count":3,"__dashboard_sub_agg__":{"value":34.249999046325684}},{"key":"Blocker","doc_count":2,"__dashboard_sub_agg__":{"value":3.7699999809265137}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"":680.4299903009087,"Blocker":3.7699999809265137,"High":25510.579972507432,"Highest":8571.340042842552,"Low":91525.1500494685,"Lowest":34.249999046325684,"Medium":9151.430022606626}}}`,
		},
		{
			name:      "top key and filter and avg",
			inputStr:  `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_sub_agg__":{"value":22.13425636021004}},{"key":"High","doc_count":646,"__dashboard_sub_agg__":{"value":39.49006187694649}},{"key":"Medium","doc_count":313,"__dashboard_sub_agg__":{"value":29.237795599382192}},{"key":"Highest","doc_count":250,"__dashboard_sub_agg__":{"value":34.28536017137021}},{"key":"","doc_count":29,"__dashboard_sub_agg__":{"value":23.463103113824438}},{"key":"Lowest","doc_count":3,"__dashboard_sub_agg__":{"value":11.416666348775228}},{"key":"Blocker","doc_count":2,"__dashboard_sub_agg__":{"value":1.8849999904632568}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"":23.463103113824438,"Blocker":1.8849999904632568,"High":39.49006187694649,"Highest":34.28536017137021,"Low":22.13425636021004,"Lowest":11.416666348775228,"Medium":29.237795599382192}}}`,
		},
		{
			name:      "top key and filter and sum and nested",
			inputStr:  `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_nested__":{"doc_count":3176,"__dashboard_sub_agg__":{"value":158.24999995529652}}},{"key":"High","doc_count":646,"__dashboard_nested__":{"doc_count":572,"__dashboard_sub_agg__":{"value":11.5}}},{"key":"Medium","doc_count":313,"__dashboard_nested__":{"doc_count":261,"__dashboard_sub_agg__":{"value":0}}},{"key":"Highest","doc_count":250,"__dashboard_nested__":{"doc_count":228,"__dashboard_sub_agg__":{"value":12.100000023841858}}},{"key":"","doc_count":29,"__dashboard_nested__":{"doc_count":28,"__dashboard_sub_agg__":{"value":2.5}}},{"key":"Lowest","doc_count":3,"__dashboard_nested__":{"doc_count":3,"__dashboard_sub_agg__":{"value":0}}},{"key":"Blocker","doc_count":2,"__dashboard_nested__":{"doc_count":0,"__dashboard_sub_agg__":{"value":0}}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"":2.5,"Blocker":0,"High":11.5,"Highest":12.100000023841858,"Low":158.24999995529652,"Lowest":0,"Medium":0}}}`,
		},
		{
			name:      "top key and filter and avg and nested",
			inputStr:  `{ "is_top_n": true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_nested__":{"doc_count":3176,"__dashboard_sub_agg__":{"value":0.04982682618239815}}},{"key":"High","doc_count":646,"__dashboard_nested__":{"doc_count":572,"__dashboard_sub_agg__":{"value":0.020104895104895104}}},{"key":"Medium","doc_count":313,"__dashboard_nested__":{"doc_count":261,"__dashboard_sub_agg__":{"value":0}}},{"key":"Highest","doc_count":250,"__dashboard_nested__":{"doc_count":228,"__dashboard_sub_agg__":{"value":0.053070175543166044}}},{"key":"","doc_count":29,"__dashboard_nested__":{"doc_count":28,"__dashboard_sub_agg__":{"value":0.08928571428571429}}},{"key":"Lowest","doc_count":3,"__dashboard_nested__":{"doc_count":3,"__dashboard_sub_agg__":{"value":0}}},{"key":"Blocker","doc_count":2,"__dashboard_nested__":{"doc_count":0,"__dashboard_sub_agg__":{"value":null}}}]}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"":0.08928571428571429,"High":0.020104895104895104,"Highest":0.053070175543166044,"Low":0.04982682618239815,"Lowest":0,"Medium":0}}}`,
		},
		{
			name:      "only filter",
			inputStr:  `{"filters":[{"name":"series1","agg_type":0,"terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"series1":10568}}}`,
		},
		{
			name:      "only filter and sum",
			inputStr:  `{"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_sub_agg__":{"value":278094.4500413425}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"series1":278094.4500413425}}}`,
		},
		{
			name:      "filter and avg",
			inputStr:  `{"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_sub_agg__":{"value":26.314766279460876}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"series1":26.314766279460876}}}`,
		},
		{
			name:      "filter and sum and nested",
			inputStr:  `{"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_nested__":{"doc_count":8407,"__dashboard_sub_agg__":{"value":195.0499999821186}}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"series1":195.0499999821186}}}`,
		},
		{
			name:      "filter and avg and nested",
			inputStr:  `{"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_nested__":{"doc_count":8407,"__dashboard_sub_agg__":{"value":0.023200904006437328}}}}}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"series1":0.023200904006437328}}}`,
		},
		{
			name:      "date histogram and top key",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"}}`,
			aggResult: `{"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":8007,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":8007}]}},{"key":"High","doc_count":1274,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":1274}]}},{"key":"Medium","doc_count":684,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":684}]}},{"key":"Highest","doc_count":511,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":511}]}},{"key":"","doc_count":100,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":100}]}},{"key":"Lowest","doc_count":20,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":20}]}},{"key":"Blocker","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2}]}},{"key":"NONE","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2}]}},{"key":"Minor","doc_count":1,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":1}]}}]}}`,
			metricStr: `{"__dashboard_top__":{"":{"2023-01-01":100},"Blocker":{"2023-01-01":2},"High":{"2023-01-01":1274},"Highest":{"2023-01-01":511},"Low":{"2023-01-01":8007},"Lowest":{"2023-01-01":20},"Medium":{"2023-01-01":684},"Minor":{"2023-01-01":1},"NONE":{"2023-01-01":2}}}`,
		},
		{
			name:      "date histogram and top key and filter",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":4135}]}},{"key":"High","doc_count":646,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":646}]}},{"key":"Medium","doc_count":313,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":313}]}},{"key":"Highest","doc_count":250,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":250}]}},{"key":"","doc_count":29,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":29}]}},{"key":"Lowest","doc_count":3,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":3}]}},{"key":"Blocker","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2}]}}]}}}}}`,
			metricStr: `{"series1":{"":{"2023-01-01":29},"Blocker":{"2023-01-01":2},"High":{"2023-01-01":646},"Highest":{"2023-01-01":250},"Low":{"2023-01-01":4135},"Lowest":{"2023-01-01":3},"Medium":{"2023-01-01":313}}}`,
		},
		{
			name:      "date histogram and top key and filter and sum",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":4135,"__dashboard_sub_agg__":{"value":91525.1500494685}}]}},{"key":"High","doc_count":646,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":646,"__dashboard_sub_agg__":{"value":25510.579972507432}}]}},{"key":"Medium","doc_count":313,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":313,"__dashboard_sub_agg__":{"value":9151.430022606626}}]}},{"key":"Highest","doc_count":250,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":250,"__dashboard_sub_agg__":{"value":8571.340042842552}}]}},{"key":"","doc_count":29,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":29,"__dashboard_sub_agg__":{"value":680.4299903009087}}]}},{"key":"Lowest","doc_count":3,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":3,"__dashboard_sub_agg__":{"value":34.249999046325684}}]}},{"key":"Blocker","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2,"__dashboard_sub_agg__":{"value":3.7699999809265137}}]}}]}}}}}`,
			metricStr: `{"series1":{"":{"2023-01-01":680.4299903009087},"Blocker":{"2023-01-01":3.7699999809265137},"High":{"2023-01-01":25510.579972507432},"Highest":{"2023-01-01":8571.340042842552},"Low":{"2023-01-01":91525.1500494685},"Lowest":{"2023-01-01":34.249999046325684},"Medium":{"2023-01-01":9151.430022606626}}}`,
		},
		{
			name:      "date histogram and top key and filter and avg",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":4135,"__dashboard_sub_agg__":{"value":22.13425636021004}}]}},{"key":"High","doc_count":646,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":646,"__dashboard_sub_agg__":{"value":39.49006187694649}}]}},{"key":"Medium","doc_count":313,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":313,"__dashboard_sub_agg__":{"value":29.237795599382192}}]}},{"key":"Highest","doc_count":250,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":250,"__dashboard_sub_agg__":{"value":34.28536017137021}}]}},{"key":"","doc_count":29,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":29,"__dashboard_sub_agg__":{"value":23.463103113824438}}]}},{"key":"Lowest","doc_count":3,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":3,"__dashboard_sub_agg__":{"value":11.416666348775228}}]}},{"key":"Blocker","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2,"__dashboard_sub_agg__":{"value":1.8849999904632568}}]}}]}}}}}`,
			metricStr: `{"series1":{"":{"2023-01-01":23.463103113824438},"Blocker":{"2023-01-01":1.8849999904632568},"High":{"2023-01-01":39.49006187694649},"Highest":{"2023-01-01":34.28536017137021},"Low":{"2023-01-01":22.13425636021004},"Lowest":{"2023-01-01":11.416666348775228},"Medium":{"2023-01-01":29.237795599382192}}}`,
		},
		{
			name:      "date histogram and top key and filter and sum and nested",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":4135,"__dashboard_nested__":{"doc_count":3176,"__dashboard_sub_agg__":{"value":158.24999995529652}}}]}},{"key":"High","doc_count":646,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":646,"__dashboard_nested__":{"doc_count":572,"__dashboard_sub_agg__":{"value":11.5}}}]}},{"key":"Medium","doc_count":313,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":313,"__dashboard_nested__":{"doc_count":261,"__dashboard_sub_agg__":{"value":0}}}]}},{"key":"Highest","doc_count":250,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":250,"__dashboard_nested__":{"doc_count":228,"__dashboard_sub_agg__":{"value":12.100000023841858}}}]}},{"key":"","doc_count":29,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":29,"__dashboard_nested__":{"doc_count":28,"__dashboard_sub_agg__":{"value":2.5}}}]}},{"key":"Lowest","doc_count":3,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":3,"__dashboard_nested__":{"doc_count":3,"__dashboard_sub_agg__":{"value":0}}}]}},{"key":"Blocker","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2,"__dashboard_nested__":{"doc_count":0,"__dashboard_sub_agg__":{"value":0}}}]}}]}}}}}`,
			metricStr: `{"series1":{"":{"2023-01-01":2.5},"Blocker":{"2023-01-01":0},"High":{"2023-01-01":11.5},"Highest":{"2023-01-01":12.100000023841858},"Low":{"2023-01-01":158.24999995529652},"Lowest":{"2023-01-01":0},"Medium":{"2023-01-01":0}}}`,
		},
		{
			name:      "date histogram and top key and filter and avg and nested",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{ "field": "priority"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"sub_task_type_story_points.program_development","terms":[{"field":"jira_status","value":["Done"]}]}]}`,
			aggResult: `{"series1":{"buckets":{"series1":{"doc_count":10568,"__dashboard_top__":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Low","doc_count":4135,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":4135,"__dashboard_nested__":{"doc_count":3176,"__dashboard_sub_agg__":{"value":0.04982682618239815}}}]}},{"key":"High","doc_count":646,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":646,"__dashboard_nested__":{"doc_count":572,"__dashboard_sub_agg__":{"value":0.020104895104895104}}}]}},{"key":"Medium","doc_count":313,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":313,"__dashboard_nested__":{"doc_count":261,"__dashboard_sub_agg__":{"value":0}}}]}},{"key":"Highest","doc_count":250,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":250,"__dashboard_nested__":{"doc_count":228,"__dashboard_sub_agg__":{"value":0.053070175543166044}}}]}},{"key":"","doc_count":29,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":29,"__dashboard_nested__":{"doc_count":28,"__dashboard_sub_agg__":{"value":0.08928571428571429}}}]}},{"key":"Lowest","doc_count":3,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":3,"__dashboard_nested__":{"doc_count":3,"__dashboard_sub_agg__":{"value":0}}}]}},{"key":"Blocker","doc_count":2,"__dashboard_date__":{"buckets":[{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":2,"__dashboard_nested__":{"doc_count":0,"__dashboard_sub_agg__":{"value":null}}}]}}]}}}}}`,
			metricStr: `{"series1":{"":{"2023-01-01":0.08928571428571429},"Blocker":{},"High":{"2023-01-01":0.020104895104895104},"Highest":{"2023-01-01":0.053070175543166044},"Low":{"2023-01-01":0.04982682618239815},"Lowest":{"2023-01-01":0},"Medium":{"2023-01-01":0}}}`,
		},

		{
			name:      "date histogram and only count",
			inputStr:  `{"is_top_n": false,"is_date_histogram":true,"filters":[{"name":"series1"}]}`,
			aggResult: `{"series1":{"buckets":[{"key_as_string":"2021-01-01","key":1609430400000,"doc_count":1},{"key_as_string":"2022-01-01","key":1640966400000,"doc_count":15},{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":1373},{"key_as_string":"2024-01-01","key":1704038400000,"doc_count":24}]}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2021-01-01":1,"2022-01-01":15,"2023-01-01":1373,"2024-01-01":24}}}`,
		},
		{
			name:      "date histogram and only sum",
			inputStr:  `{"is_top_n": false,"is_date_histogram":true,"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day"}]}`,
			aggResult: `{"series1":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":3,"__dashboard_sub_agg__":{"value":141}},{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":70,"__dashboard_sub_agg__":{"value":2407.8200164437294}},{"key_as_string":"2020-01-01","key":1577808000000,"doc_count":460,"__dashboard_sub_agg__":{"value":22119.51997485198}},{"key_as_string":"2021-01-01","key":1609430400000,"doc_count":2814,"__dashboard_sub_agg__":{"value":95344.96002810635}},{"key_as_string":"2022-01-01","key":1640966400000,"doc_count":9437,"__dashboard_sub_agg__":{"value":265553.3400535807}},{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":9938,"__dashboard_sub_agg__":{"value":180554.7500198651}}]}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2018-01-01":141,"2019-01-01":2407.8200164437294,"2020-01-01":22119.51997485198,"2021-01-01":95344.96002810635,"2022-01-01":265553.3400535807,"2023-01-01":180554.7500198651}}}`,
		},
		{
			name:      "data histogram and only avg",
			inputStr:  `{"is_top_n": false,"is_date_histogram":true,"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day"}]}`,
			aggResult: `{"series1":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":3,"__dashboard_sub_agg__":{"value":47}},{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":70,"__dashboard_sub_agg__":{"value":34.39742880633899}},{"key_as_string":"2020-01-01","key":1577808000000,"doc_count":460,"__dashboard_sub_agg__":{"value":48.08591298880865}},{"key_as_string":"2021-01-01","key":1609430400000,"doc_count":2814,"__dashboard_sub_agg__":{"value":33.88235964040737}},{"key_as_string":"2022-01-01","key":1640966400000,"doc_count":9437,"__dashboard_sub_agg__":{"value":28.148541451513747}},{"key_as_string":"2023-01-01","key":1672502400000,"doc_count":9938,"__dashboard_sub_agg__":{"value":18.16811732942897}}]}}`,
			metricStr: `{"series1":{"__dashboard_default_value__":{"2018-01-01":47,"2019-01-01":34.39742880633899,"2020-01-01":48.08591298880865,"2021-01-01":33.88235964040737,"2022-01-01":28.148541451513747,"2023-01-01":18.16811732942897}}}`,
		},
		//{
		//	name: "date histogram and only sum and nested",
		//},
		//{
		//	name: "data histogram and only avg and nested",
		//},
		{
			name:      "top k and only count",
			inputStr:  `{"is_top_n": true,"is_date_histogram":false, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":0}]}`,
			aggResult: `{"series1":{"doc_count_error_upper_bound":0,"sum_other_doc_count":142,"buckets":[{"key":"Done","doc_count":13486},{"key":"Closed","doc_count":4896},{"key":"Waiting","doc_count":2581},{"key":"Developing","doc_count":372},{"key":"Icebox","doc_count":371},{"key":"Staging","doc_count":354},{"key":"PRD","doc_count":218},{"key":"Delivering","doc_count":148},{"key":"Testing","doc_count":83},{"key":"Designing","doc_count":71}]}}`,
			metricStr: `{"series1":{"Closed":{"__dashboard_default_value__":4896},"Delivering":{"__dashboard_default_value__":148},"Designing":{"__dashboard_default_value__":71},"Developing":{"__dashboard_default_value__":372},"Done":{"__dashboard_default_value__":13486},"Icebox":{"__dashboard_default_value__":371},"PRD":{"__dashboard_default_value__":218},"Staging":{"__dashboard_default_value__":354},"Testing":{"__dashboard_default_value__":83},"Waiting":{"__dashboard_default_value__":2581}}}`,
		},
		{
			name:      "top k and only sum",
			inputStr:  `{"is_top_n": true,"is_date_histogram":false, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day"}]}`,
			aggResult: `{"series1":{"doc_count_error_upper_bound":0,"sum_other_doc_count":142,"buckets":[{"key":"Done","doc_count":13486,"__dashboard_sub_agg__":{"value":563394.0300785992}},{"key":"Closed","doc_count":4896,"__dashboard_sub_agg__":{"value":2056.000007737428}},{"key":"Waiting","doc_count":2581,"__dashboard_sub_agg__":{"value":81.70999908447266}},{"key":"Developing","doc_count":372,"__dashboard_sub_agg__":{"value":51.679999351501465}},{"key":"Icebox","doc_count":371,"__dashboard_sub_agg__":{"value":305.2900085449219}},{"key":"Staging","doc_count":354,"__dashboard_sub_agg__":{"value":12.770000457763672}},{"key":"PRD","doc_count":218,"__dashboard_sub_agg__":{"value":0}},{"key":"Delivering","doc_count":148,"__dashboard_sub_agg__":{"value":202.9699993133545}},{"key":"Testing","doc_count":83,"__dashboard_sub_agg__":{"value":16.02999973297119}},{"key":"Designing","doc_count":71,"__dashboard_sub_agg__":{"value":0.9100000262260437}}]}}`,
			metricStr: `{"series1":{"Closed":{"__dashboard_default_value__":2056.000007737428},"Delivering":{"__dashboard_default_value__":202.9699993133545},"Designing":{"__dashboard_default_value__":0.9100000262260437},"Developing":{"__dashboard_default_value__":51.679999351501465},"Done":{"__dashboard_default_value__":563394.0300785992},"Icebox":{"__dashboard_default_value__":305.2900085449219},"PRD":{"__dashboard_default_value__":0},"Staging":{"__dashboard_default_value__":12.770000457763672},"Testing":{"__dashboard_default_value__":16.02999973297119},"Waiting":{"__dashboard_default_value__":81.70999908447266}}} `,
		},
		{
			name:      "top k and only avg",
			inputStr:  `{"is_top_n": true,"is_date_histogram":false, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day"}]}`,
			aggResult: `{"series1":{"doc_count_error_upper_bound":0,"sum_other_doc_count":142,"buckets":[{"key":"Done","doc_count":13486,"__dashboard_sub_agg__":{"value":41.77621459873938}},{"key":"Closed","doc_count":4896,"__dashboard_sub_agg__":{"value":0.41993464210323284}},{"key":"Waiting","doc_count":2581,"__dashboard_sub_agg__":{"value":0.031695112135171705}},{"key":"Developing","doc_count":372,"__dashboard_sub_agg__":{"value":0.13892472943952006}},{"key":"Icebox","doc_count":371,"__dashboard_sub_agg__":{"value":0.8228841200671748}},{"key":"Staging","doc_count":354,"__dashboard_sub_agg__":{"value":0.036073447620801335}},{"key":"PRD","doc_count":218,"__dashboard_sub_agg__":{"value":0}},{"key":"Delivering","doc_count":148,"__dashboard_sub_agg__":{"value":1.3714189142794222}},{"key":"Testing","doc_count":83,"__dashboard_sub_agg__":{"value":0.19313252690326738}},{"key":"Designing","doc_count":71,"__dashboard_sub_agg__":{"value":0.012816901777831602}}]}}`,
			metricStr: `{"series1":{"Closed":{"__dashboard_default_value__":0.41993464210323284},"Delivering":{"__dashboard_default_value__":1.3714189142794222},"Designing":{"__dashboard_default_value__":0.012816901777831602},"Developing":{"__dashboard_default_value__":0.13892472943952006},"Done":{"__dashboard_default_value__":41.77621459873938},"Icebox":{"__dashboard_default_value__":0.8228841200671748},"PRD":{"__dashboard_default_value__":0},"Staging":{"__dashboard_default_value__":0.036073447620801335},"Testing":{"__dashboard_default_value__":0.19313252690326738},"Waiting":{"__dashboard_default_value__":0.031695112135171705}}}`,
		},
		//{
		//	name: "top k and only sum and nested",
		//},
		//{
		//	name: "top k and only avg and nested",
		//},
		{
			name:      "date histogram and top k and only count",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":0}]}`,
			aggResult: `{"series1":{"doc_count_error_upper_bound":0,"sum_other_doc_count":142,"buckets":[{"key":"Done","doc_count":13486,"__dashboard_date__":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":2},{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":33},{"key_as_string":"2020-01-01","key":1577808000000,"doc_count":348}]}},{"key":"Closed","doc_count":4896,"__dashboard_date__":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":1}]}},{"key":"Waiting","doc_count":2581,"__dashboard_date__":{"buckets":[{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":6}]}}]}}`,
			metricStr: `{"series1":{"Closed":{"2018-01-01":1},"Done":{"2018-01-01":2,"2019-01-01":33,"2020-01-01":348},"Waiting":{"2019-01-01":6}}}`,
		},
		{
			name:      "date histogram and top k and only sum",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":3,"extra_agg_field":"total_cost_day"}]}`,
			aggResult: `{"series1":{"doc_count_error_upper_bound":0,"sum_other_doc_count":142,"buckets":[{"key":"Done","doc_count":13486,"__dashboard_date__":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":2,"__dashboard_sub_agg__":{"value":141}}]}},{"key":"Closed","doc_count":4896,"__dashboard_date__":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":1,"__dashboard_sub_agg__":{"value":0}},{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":27,"__dashboard_sub_agg__":{"value":141.16000366210938}}]}},{"key":"Waiting","doc_count":2581,"__dashboard_date__":{"buckets":[{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":6,"__dashboard_sub_agg__":{"value":0}},{"key_as_string":"2020-01-01","key":1577808000000,"doc_count":0,"__dashboard_sub_agg__":{"value":0}}]}}]}}`,
			metricStr: `{"series1":{"Closed":{"2018-01-01":0,"2019-01-01":141.16000366210938},"Done":{"2018-01-01":141},"Waiting":{"2019-01-01":0,"2020-01-01":0}}}`,
		},
		{
			name:      "date histogram and top k and only avg",
			inputStr:  `{"is_top_n": true,"is_date_histogram":true, "top_n":{"field":"jira_status"},"filters":[{"name":"series1","agg_type":2,"extra_agg_field":"total_cost_day"}]}`,
			aggResult: ` {"series1":{"doc_count_error_upper_bound":27,"sum_other_doc_count":4332,"buckets":[{"key":"Done","doc_count":13491,"__dashboard_date__":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":2,"__dashboard_sub_agg__":{"value":141}},{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":33,"__dashboard_sub_agg__":{"value":2266.66001278162}}]}},{"key":"Closed","doc_count":4899,"__dashboard_date__":{"buckets":[{"key_as_string":"2018-01-01","key":1514736000000,"doc_count":1,"__dashboard_sub_agg__":{"value":0}},{"key_as_string":"2019-01-01","key":1546272000000,"doc_count":27,"__dashboard_sub_agg__":{"value":141.16000366210938}}]}}]}}`,
			metricStr: `{"series1":{"Closed":{"2018-01-01":0,"2019-01-01":141.16000366210938},"Done":{"2018-01-01":141,"2019-01-01":2266.66001278162}}}`,
		},
		//{
		//	name: "date histogram and top k and only sum and nested",
		//},
		//{
		//	name: "date histogram and top k and only avg and nested",
		//},
	}

	ctx := context.TODO()
	for idx, test := range tests {
		aggParam := define.QueryOutputAggregator{}
		if err := json.Unmarshal([]byte(test.inputStr), &aggParam); err != nil {
			t.Errorf("test %v name: %v failed, err: %v", idx, test.name, err)
			continue
		}
		esResult := make(map[string]interface{}, 0)
		if err := json.Unmarshal([]byte(test.aggResult), &esResult); err != nil {
			t.Errorf("test %v name: %v failed, err: %v", idx, test.name, err)
			continue
		}
		metrics, err := GetESAggBucketsV2HelperSeriesMetric(ctx, esResult, aggParam)
		if err != nil {
			t.Errorf("test %v name: %v failed, err: %v", idx, test.name, err)
			continue
		}

		byteMetric, osErr := json.Marshal(metrics)
		if osErr != nil {
			t.Errorf("test %v name: %v failed, err: %v", idx, test.name, osErr)
			continue
		}
		if strings.TrimSpace(test.metricStr) != string(byteMetric) {
			t.Errorf("test %v name: %v result metric no match. result: %s ", idx, test.name, string(byteMetric))
			continue
		}

		/*fmt.Println("{")
		fmt.Printf("name: \"%s\",\n", test.name)
		fmt.Printf("inputStr: `%s`,\n", test.inputStr)
		fmt.Printf("aggResult: `%s`,\n", test.aggResult)
		fmt.Printf("metricStr: `%s`,\n", string(byteMetric))
		fmt.Println("},")*/

	}
}
