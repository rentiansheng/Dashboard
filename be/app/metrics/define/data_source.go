package define

import (
	"fmt"

	"github.com/rentiansheng/dashboard/middleware/errors"
)

/***************************
    @author: tiansheng.ren
    @date: 2025/6/9
    @desc:

***************************/

type MetricDataSource struct {
	ID             uint64 `json:"id"`
	DisplayName    string `json:"name"`
	Tips           string `json:"tips"`
	Remark         string `json:"remark"`
	DataName       string `json:"data_name"`
	Status         uint   `json:"status"`
	DataSourceType string `json:"data_type"`
	SortFields     string `json:"sort_fields"`
	Formatter      string `json:"formatter"`
	Template       string `json:"template"`
	GroupKey       string `json:"group_key"`

	EnableGroupKeyName bool `json:"enable_group_key_name"`

	Mtime uint32 `json:"mtime"`
	Ctime uint32 `json:"ctime"`
}

type MetricDataSourceMeta struct {
	ID             uint64                         `json:"id"`
	DataSourceID   uint64                         `json:"data_source_id"`
	FieldName      string                         `json:"name"`
	DisplayName    string                         `json:"display_name"`
	DataType       string                         `json:"data_type"`
	OutputDataType string                         `json:"output_data_type"`
	FieldTips      string                         `json:"field_tips"`
	Action         MetricDataSourceMetaEnumAction `json:"action"`
	Enum           MetricDataSourceMetaEnum       `json:"enum"`
	Formatter      string                         `json:"formatter"`
	Template       string                         `json:"template"`

	Mtime uint32 `json:"mtime"`
	Ctime uint32 `json:"ctime"`

	Nested string `json:"nested_path"`
}

type MetricDataSourceMetaEnumAction struct {
	Sort   bool `json:"sort"`
	Filter bool `json:"filter"`
	Key    bool `json:"key"`
	Time   bool `json:"time"`
	API    bool `json:"api"`
	Detail bool `json:"detail"`
}

type MetricDataSourceMetaEnum struct {
	API struct {
		Path    string `json:"path"`
		Dynamic bool   `json:"dynamic"`
	} `json:"api"`
	Values DataSourceFieldEnum `json:"values"`
}

type QueryReq struct {
	Date TimeRange `json:"date"`
	// 1柱状图,2折线图,3饼状图,100:详情
	Typ          string         `json:"type"`
	Indexes      []QueryIndexes `json:"indexes"`
	GroupKeyId   uint64         `json:"group_key_id" validate:"required"`
	Page         Page           `json:"page"`
	DataSourceID uint64         `json:"data_source_id"`
}

func (qr *QueryReq) Default() {
	(&(qr.Page)).Default()
	for idx, index := range qr.Indexes {
		if index.Output.Aggregator.AggType == QueryOutputAggTypeTopN ||
			index.Output.Aggregator.AggType == QueryOutputAggTypeHistogramTopN {
			if index.Output.Aggregator.TopN.Num == 0 {
				qr.Indexes[idx].Output.Aggregator.TopN.Num = 10
			}
		}
		if index.Output.Aggregator.AggType == QueryOutputAggTypeXAxisField {
			if index.Output.Aggregator.XAxis.Num == 0 {
				qr.Indexes[idx].Output.Aggregator.XAxis.Num = 10
			}
		}
	}
}

func (qr QueryReq) Validate() error {
	for _, index := range qr.Indexes {
		if err := index.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type QueryIndexes struct {
	Name    string        `json:"name"`
	Output  QueryOutput   `json:"output"`
	Filters []QueryFilter `json:"filters"`
}

func (qi QueryIndexes) Validate() error {
	return qi.Output.Validate()
}

type QueryOutput struct {
	TimeField  string                `json:"time_field"`
	Sorts      []string              `json:"sorts"`
	Cycle      MetricCycleType       `json:"cycle"`
	Aggregator QueryOutputAggregator `json:"aggregator"`
	Fields     []string              `json:"fields"`

	// 柱状图时候，为空表示使用时间模式下，非空用字段的值作为x轴
	XAxisField string `json:"x_axis_field"`
}

func (qo QueryOutput) Validate() error {

	return qo.Aggregator.Validate()
}

type QueryOutputAggregator struct {
	V2              bool `json:"v2"`
	IsDataHistogram bool `json:"is_date_histogram"`
	IsTopN          bool `json:"is_top_n"`
	// only for top n is true
	TopN    QueryOutputAggregatorTop            `json:"top_n"`
	Filters []QueryOutputAggregatorXAxisFilters `json:"filters"`

	// stupid code. 后续太难维护了。

	// 0:无,1:top_n,2:terms, 3. date_histogram and topn, 4. date_histogram and terms, 5: x_axis use field
	AggType uint8                          `json:"agg_type"`
	Terms   []QueryOutputAggregatorTerm    `json:"terms"`
	XAxis   QueryOutputAggregatorTermXAxis `json:"x_axis"`
	Name    string                         `json:"name"`
}

func (qoa QueryOutputAggregator) IsDateHistogram() bool {
	return qoa.AggType == QueryOutputAggTypeDateHistogram || qoa.AggType == QueryOutputAggTypeHistogramTopN
}

func (qoa QueryOutputAggregator) Validate() error {
	switch qoa.AggType {
	case QueryOutputAggTypeDateHistogram:
		return nil
	case QueryOutputAggTypeTopN, QueryOutputAggTypeHistogramTopN:
		return qoa.TopN.Validate()
	case QueryOutputAggTypeTerms, QueryOutputAggTypeHistogramTerms:
		for _, term := range qoa.Terms {
			if err := term.Validate(); err != nil {
				return err
			}
		}
	case QueryOutputAggTypeXAxisField:
		if err := qoa.XAxis.Validate(); err != nil {
			return err
		}

	default:
		return nil
	}
	return nil
}

const (
	// QueryOutputAggTypeDateHistogram 时间模式， 默认值
	QueryOutputAggTypeDateHistogram uint8 = 0
	// QueryOutputAggTypeTopN top n 模式
	QueryOutputAggTypeTopN uint8 = 1
	// QueryOutputAggTypeTerms terms 模式,
	QueryOutputAggTypeTerms uint8 = 2
	// QueryOutputAggTypeHistogramTopN date_histogram and topn
	QueryOutputAggTypeHistogramTopN uint8 = 3
	// QueryOutputAggTypeHistogramTerms date_histogram and terms
	QueryOutputAggTypeHistogramTerms uint8 = 4
	// QueryOutputAggTypeXAxisField x轴字段模式
	QueryOutputAggTypeXAxisField uint8 = 5
)

type QueryOutputAggregatorTop struct {
	Num   int64  `json:"num"`
	Field string `json:"field"`
}

func (qo QueryOutputAggregatorTop) Validate() error {
	if qo.Field == "" {
		return fmt.Errorf("output.aggregator.top_n.field is empty")
	}
	return nil
}

type QueryOutputAggregatorTerm struct {
	Field string        `json:"field"`
	Value []interface{} `json:"value"`
	Name  string        `json:"name"`
}

func (qo QueryOutputAggregatorTerm) Validate() error {
	if qo.Field == "" {
		return fmt.Errorf("output.aggregator.terms.field is empty")
	}
	if len(qo.Value) == 0 {
		return fmt.Errorf("output.aggregator.terms.value is empty")
	}
	return nil
}

type QueryOutputAggregatorXAxisFilters struct {
	Name  string                                   `json:"name"`
	Terms []QueryOutputAggregatorXAxisFiltersTerm  `json:"terms"`
	Range []QueryOutputAggregatorXAxisFiltersRange `json:"range"`
	// 0: count, 1: avg, 2: sum
	AggType             QueryOutputAggregatorXAxisFiltersAggType `json:"agg_type"`
	ExtraAggFieldNested string                                   `json:"-"`
	// 非必填，只有avg, sum 时候需要
	ExtraAggField string `json:"extra_agg_field"`
}

type QueryOutputAggregatorXAxisFiltersTerm struct {
	Field string        `json:"field"`
	Value []interface{} `json:"value"`

	// 通过字段定义获取到的, 暂不支持
	Nesteds []string `json:"-"`
}

type QueryOutputAggregatorXAxisFiltersRange struct {
	Field string      `json:"field"`
	From  interface{} `json:"from"`
	To    interface{} `json:"to"`

	// 通过字段定义获取到的，暂不支持
	Nesteds []string `json:"-"`
}

type QueryOutputAggregatorXAxisFiltersAggType uint8

const (
	QueryOutputAggregatorXAxisFiltersAggTypeDefault QueryOutputAggregatorXAxisFiltersAggType = 0
	QueryOutputAggregatorXAxisFiltersAggTypeCount   QueryOutputAggregatorXAxisFiltersAggType = 1
	QueryOutputAggregatorXAxisFiltersAggTypeAvg     QueryOutputAggregatorXAxisFiltersAggType = 2
	QueryOutputAggregatorXAxisFiltersAggTypeSum     QueryOutputAggregatorXAxisFiltersAggType = 3
)

type QueryOutputAggregatorTermXAxis struct {
	Num     int64                               `json:"num"`
	Field   string                              `json:"field"`
	Filters []QueryOutputAggregatorXAxisFilters `json:"filters"`
}

func (q QueryOutputAggregatorTermXAxis) Validate() error {
	if q.Field == "" {
		return fmt.Errorf("output.aggregator.x_axis.field is empty")
	}
	if q.Num <= 0 {
		return fmt.Errorf("output.aggregator.x_axis.num is empty")
	}
	if len(q.Filters) == 0 {
		return fmt.Errorf("output.aggregator.x_axis.filters is empty")
	}

	for idx, f := range q.Filters {
		if f.Name == "" {
			return fmt.Errorf("output.aggregator.x_axis.filters[%d].name is empty", idx)
		}
		if f.AggType != 0 {
			if f.ExtraAggField == "" {
				return fmt.Errorf("output.aggregator.x_axis.terms[%d].extra_agg_field is empty", idx)
			}
		}

	}
	return nil
}

type QueryFilter struct {
	Name     string            `json:"name"`
	Operator string            `json:"join_operator"`
	Rules    []QueryFilterRule `json:"rules"`
}

type QueryFilterRule struct {
	Field    string                  `json:"field"`
	Operator QueryFilterRuleOperator `json:"operator"`
	Values   []interface{}           `json:"values"`
	Pause    bool                    `json:"pause"`
	Not      bool                    `json:"not"`

	// nested 字段
	Nesteds []string `json:"-"`
}

type QueryFilterRuleOperator string

const (
	QueryFilterRuleEQ      QueryFilterRuleOperator = "eq"
	QueryFilterRuleIn      QueryFilterRuleOperator = "in"
	QueryFilterRuleBetween QueryFilterRuleOperator = "between"
	QueryFilterRuleLt      QueryFilterRuleOperator = "lt"
	QueryFilterRuleLte     QueryFilterRuleOperator = "lte"
	QueryFilterRuleGt      QueryFilterRuleOperator = "gt"
	QueryFilterRuleGte     QueryFilterRuleOperator = "gte"
	// QueryFilterRuleLike 系统的like 是前缀匹配, QueryFilterRulePrefixLike 别名
	QueryFilterRuleLike QueryFilterRuleOperator = "like"
	// QueryFilterRulePrefix  系统的like 是前缀匹配 别名
	QueryFilterRulePrefix QueryFilterRuleOperator = "prefix"
	// QueryFilterRuleSuffix 前缀模糊查询
	QueryFilterRuleSuffix QueryFilterRuleOperator = "suffix"
)

const (
	DataSourceTypeES = "es"
)

type DataSourceFieldEnumReq struct {
	DepartmentId uint64        `json:"department_id"`
	DataSourceId uint64        `json:"data_source_id"`
	Date         TimeRange     `json:"date"`
	FieldName    string        `json:"field_name"`
	FieldValue   []interface{} `json:"field_value"`
	// FieldValueByKey 用key 作为enum 查询条件，未支持
	FieldValueByKey bool `json:"field_value_by_key"`
	// reserved
	Relations []DataSourceFieldEEnumRelations `json:"relations"`
}

type DataSourceFieldEEnumRelations struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type DataSourceFieldEnum struct {
	Typ    string      `json:"type,omitempty"`
	Values interface{} `json:"values,omitempty"`
}

const (
	DataSourceFieldEnumTypeArray = "array"
	DataSourceFieldEnumTypeKv    = "kv"
)

type DataSourceDetailMeta struct {
	Fields  []MetricDetailMeta                 `json:"fields"`
	Enum    map[string]map[string]interface{}  `json:"enums"`
	EnumAPI map[string]MetricDetailMetaEnumAPI `json:"enum_api"`
	// 拷贝枚举转换后的值 map[to]from
	CopyFieldValue map[string]string `json:"copy_field_value"`
	// 拷贝枚举转换前的值 map[to]from
	CopyRawFieldValue map[string]string `json:"copy_raw_field_value"`
	// 基本上专用方案
	VirtualFieldFn map[string]VirtualField `json:"virtual_field"`
}

type MetricDetailMetaEnumAPI struct {
	Name string `json:"name"`
}

type VirtualField struct {
	// 用来做数据格转换，避免不同数据字段名字不同,  map[row field]fn handle field
	ValueFields map[string]string `json:"value_fields"`
	FuncName    string            `json:"fn"`
}

type DataSourceMetaReq struct {
	ID uint64 `json:"data_source_id"`
}

type DataSourceDataResp struct {
	List  []MetricDataSource `json:"list"`
	Total int64              `json:"total"`
}

type DataSourceFieldMetaResp struct {
	Fields []MetricDataSourceMeta `json:"fields"`
}

type Base struct {
	RetCode int32  `json:"retcode"`
	Message string `json:"message"`
}

func (b Base) Error() errors.Error {
	if b.RetCode == 0 {
		return nil
	}
	return errors.NewError(b.RetCode, b.Message)
}
