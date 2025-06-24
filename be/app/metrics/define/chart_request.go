package define

type ChartReq struct {
	Name       string    `json:"metric_name" form:"metric_name" validate:"required,min=5,max=128"`
	MetricKeys []string  `json:"metric_keys" form:"metric_keys"`
	Time       TimeRange `json:"time" form:"time"`
	Cycle      *int8     `json:"cycle" form:"cycle"`
	Horizontal bool      `json:"horizontal" form:"horizontal"`

	SeriesHook ChartReqSeriesHook `json:"series_hook"`

	DisableSnapshot *bool `json:"disable_snapshot" form:"disable_snapshot"`
}

type ChartReqSeriesHook struct {

	// 需要保留的series， 为空保留全部
	Series []string `json:"series"`
	// 计算公式
	Formulas []ChartReqSeriesHookFormula `json:"formula"`
}

type ChartReqSeriesHookFormula struct {
	// 1 sum, 2 公式计算
	Typ int `json:"type"`
	// 计算公式, 1 的时候是 需要求和字段，为空全部字段，2. 是计算公式 eg:val["series_name_1"]/val["series_name_2"]
	Formula string `json:"formula"`
	// 指标名字
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (r ChartReq) DisableChartSnapshot() bool {
	return r.DisableSnapshot != nil && *r.DisableSnapshot == true
}

type ChartListResp struct {
	List  []ChartListRespListItem `json:"list"`
	Total int64                   `json:"total"`
}

type ChartDetailMetaReq struct {
	Name        string   `json:"metric_name" form:"metric_name" validate:"required,min=5,max=128"`
	Cycle       *int8    `json:"cycle"  form:"cycle"`
	SeriesNames []string `json:"series_names" form:"series_names"`
}

type ChartDetailMetaFilterReq struct {
	MetricName string `json:"metric_name" form:"metric_name"`
}

type ChartListRespListItem struct {
	ChartChartName     string `json:"metric_chart_name"`
	ChartDisplayName   string `json:"metric_chart_display_name"`
	ChartChartNameTips string `json:"metric_chart_name_tips"`
	ClassifyName       string `json:"classify_name"`
	GroupName          string `json:"group_name"`
	// deprecated compatible code . not use cycle polymerization
	Cycle     int8   `json:"cycle"`
	Cycles    []int8 `json:"cycles"`
	ChartType string `json:"chart_type"`
}

type ListReq struct {
	Page
	ClassifyName string `json:"classify_name" form:"classify_name"`
	DisplayAll   bool   `json:"display_all" form:"display_all"`
}

type CreateTaskReq struct {
	// 任务的名字，同时也是指标名字
	Name string `json:"name" validate:"required"`
	// 指标周期（1年，2季，3月，4周，5日，6时）（接口必须）
	Cycle int8 `json:"cycle"  validate:"required"`
	// 周期执行方式，1 周期结束后执行，2周期中每天计算一次
	CycleMode int8 `json:"cycle_mode"  validate:"required"`
	// 开始处理任务的时间， 有start+cycle 可以选出结束时间
	Start       uint64 `json:"start"  validate:"required"`
	Collect     string `json:"collect"  validate:"required"`
	Filters     string `json:"filters"  validate:"required" `
	Aggregators string `json:"aggregators"  validate:"required" `
	Output      string `json:"output"  validate:"required"`
	// 任务状态， 1 正常，可以允许， 2. 暂停，不被执行 3. 待删除
	Status         int8           `json:"status"  validate:"required" `
	Power          map[string]int `json:"power"  validate:"required" `
	LastFinishTime int64          `json:"last_finish_time"  validate:"required"`
}

type UpdateTaskReq struct {
	ID uint64 `json:"id" validate:"required"`
	CreateTaskReq
}

type ChartDetailPageInfoReq struct {
	Name       string                        `json:"metric_name"  query:"metric_name"  validate:"required,min=5,max=128"`
	MetricKeys []string                      `json:"metric_keys" query:"metric_keys"`
	XAxisName  string                        `json:"x_axis_name" query:"x_axis_name" validate:"required"`
	Cycle      *int8                         `json:"cycle" query:"cycle"  form:"cycle"`
	Filter     *ChartDetailPageInfoReqFilter `json:"filter"`
	Page
}

type ChartDetailExportReq struct {
	*ChartDetailPageInfoReq
	ExportFields []string `json:"export_fields"`
}

type ChartDetailPageInfoReqFilter struct {
	SeriesName string            `json:"series_name"`
	Conditions []QueryFilterRule `json:"conditions"`
	SortField  string            `json:"sort_field"`
}

type ChartListRespV2 struct {
	List  []ChartListRespListItemV2 `json:"list"`
	Total int64                     `json:"total"`
}

type ChartListRespListItemV2 struct {
	ChartChartName     string `json:"metric_chart_name"`
	ChartChartNameTips string `json:"metric_chart_name_tips"`
	Tips               string `json:"tips"`
	ClassifyName       string `json:"classify_name"`
	GroupName          string `json:"group_name"`
	Cycles             []int8 `json:"cycles"`
	ChartType          string `json:"chart_type"`
}

type ChartDetailReq struct {
	Id uint64 `json:"id" validate:"required"`
}
