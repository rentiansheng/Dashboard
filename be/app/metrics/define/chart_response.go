package define

import (
	"encoding/json"
	"fmt"
)

const (
	ChartTitleAlignCenter = "center"
	ChartTypeStackedLine  = "line"
	ChartTypeBar          = "bar"
	ChartTypePie          = "pie"
	// ChartTypeSTackedBar 主要是data source 指标输出使用
	ChartTypeSTackedBar = "stacked bar"
	OutputDetail        = "detail"
)

// ChartDesc 返回值描述chart 信息
type ChartDesc struct {
	Title string `json:"title"`
	//TitleAlign string `json:"title_align"`
	Tips  string `json:"tips"`
	Cycle int8   `json:"cycle"`
}

type ChartXAxis struct {
	Typ    AxisType `json:"type"`
	Hidden bool     `json:"hidden"`
	//   数组表示所在位置的index，
	Data []string `json:"data"`
}

type ChartYAxis struct {
	Typ AxisType `json:"type"`
	//   数组表示所在位置的index，
	Data     []string `json:"data,omitempty"`
	Title    string   `json:"Title,omitempty"`
	Position string   `json:"position,omitempty"`
	Format   string   `json:"formatter,omitempty"`
}

type ChartLegend struct {
	Position ChartLegendPositionType `json:"position"`
	Hidden   bool                    `json:"hidden"`
}

type ChartSeriesItem struct {
	Name   string     `json:"name"`
	Color  string     `json:"color"`
	Type   string     `json:"type"`
	Format string     `json:"formatter,omitempty"`
	Power  int8       `json:"power"`
	Data   []*float64 `json:"data"`
	From   string     `json:"-"  copy:"from_name"`
}

type ChartSeries struct {
	Stack bool              `json:"stack"`
	Items []ChartSeriesItem `json:"items"`
	Extra ChartSeriesExtra  `json:"-"`
}

type ChartSeriesExtra struct {
	Data       map[string][]*float64 `json:"data"`
	XAxisCount int                   `json:"x_axis_count"`
}

type Chart struct {
	ChartDesc   ChartDesc    `json:"chart"`
	XAxis       ChartXAxis   `json:"x_axis,omitempty"`
	YAxis       []ChartYAxis `json:"y_axis,omitempty"`
	Legend      ChartLegend  `json:"legend"`
	Series      ChartSeries  `json:"series"`
	ChartDetail bool         `json:"chart_detail"`
	FilterCfg   struct {
		Series    bool `json:"series"`
		Condition bool `json:"condition"`
	} `json:"filter_config"`
}

type AxisType string

const (
	// AxisTypeValue value: 数值轴
	AxisTypeValue AxisType = "value"
	// AxisTypeTime time: 时间轴，
	AxisTypeTime AxisType = "time"
	// AxisTypeCategory category： 类目，来源于data
	AxisTypeCategory AxisType = "category"
)

type ChartLegendPositionType string

const (
	// 位置信息 top, left, right, bottom

	ChartLegendSPositionTop    ChartLegendPositionType = "top"
	ChartLegendSPositionLeft   ChartLegendPositionType = "left"
	ChartLegendSPositionRight  ChartLegendPositionType = "right"
	ChartLegendSPositionBottom ChartLegendPositionType = "bottom"
)

// MetricDetailMeta 指标详情数据描述
type MetricDetailMeta struct {
	// 详情中字段名， 展示表格中显示数据key
	Name string `json:"name"`
	// 显示的名字
	DisplayName string `json:"display_name"`
	// 用来做提示
	Tips string `json:"tips"`
	// 数据类型, json,link,string,int,float,bool
	Typ       string `json:"type"`
	Width     uint32 `json:"width"`
	Formatter string `json:"formatter"`
	Template  string `json:"template"`
	Nested    string `json:"nested_path"`
	Hide      bool   `json:"hide"`
	// 导出的时候，是http 导出的是名字不是link
	ExportLinkName bool `json:"export_link_name"`

	Action MetricDataSourceMetaEnumAction `json:"action"`
	Enum   MetricDataSourceMetaEnum       `json:"enum"`
}

func (m MetricDetailMeta) ToMetricDetailMetaResp() MetricDetailMetaResp {
	return MetricDetailMetaResp{
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Tips:        m.Tips,
		Typ:         m.Typ,
		Width:       m.Width,
		Hide:        m.Hide,
	}
}

// MetricDetailMetaResp 指标详情数据描述
type MetricDetailMetaResp struct {
	// 详情中字段名， 展示表格中显示数据key
	Name string `json:"name"`
	// 显示的名字
	DisplayName string `json:"display_name"`
	// 用来做提示
	Tips string `json:"tips"`
	// 数据类型, json,link,string,int,float,bool,es time
	Typ   string `json:"type"`
	Width uint32 `json:"width"`
	Hide  bool   `json:"hide"`
}

const (
	MetricDetailMetaTypeLink   = "link"
	MetricDetailMetaTypeESTime = "es time"
	MetricDetailMetaTypeTime   = "time"
	MetricDetailMetaTypeStr    = "string"
	MetricDetailMetaTypeBool   = "bool"
	MetricDetailMetaTypeDate   = "seconds to date"
)

const (
	MetricDetailIDFieldName        = "_series_name"
	MetricDetailIDFieldDisplayName = "series name"
)

func (m MetricDetailMeta) IsLink() bool {
	return m.Typ == MetricDetailMetaTypeLink
}

func (m MetricDetailMeta) IsESTime() bool {
	return m.Typ == MetricDetailMetaTypeESTime
}

func (m MetricDetailMeta) IsSeconds2Date() bool {
	return m.Typ == MetricDetailMetaTypeDate
}

type LinkObj struct {
	// Link url
	Link string `json:"link"`
	// Name display value
	Name string `json:"name"`
}

func (m MetricDetailMeta) LinkValue(rawValue, displayValue interface{}) LinkObj {
	data, err := m.unmarshalToMap(rawValue)
	if err == nil {
		return data
	}
	data, err = m.unmarshalToMap(displayValue)
	if err == nil {
		return data
	}
	data = LinkObj{
		Link: fmt.Sprintf("%s", displayValue),
		Name: fmt.Sprintf("%s", rawValue),
	}
	return data
}

func (m MetricDetailMeta) unmarshalToMap(rawValue interface{}) (link LinkObj, err error) {
	bodys := []byte{}
	switch v := rawValue.(type) {
	case []byte:
		bodys = v
	case string:
		bodys = []byte(v)
	default:
		bodys, err = json.Marshal(rawValue)
		if err != nil {
			return link, err
		}
	}

	err = json.Unmarshal(bodys, &link)
	if err != nil {
		return link, err
	}

	return link, nil
}
