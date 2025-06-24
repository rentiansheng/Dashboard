package format

import (
	"github.com/rentiansheng/dashboard/app/metrics/datasource/tools"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"strings"
	"time"
)

func chartSeriesStyle(seriesStyle string) (string, bool) {
	switch seriesStyle {
	case define.ChartTypeSTackedBar:
		return define.ChartTypeBar, true
	case define.ChartTypeStackedLine:
		return define.ChartTypeStackedLine, true
	case define.ChartTypePie:
		return define.ChartTypePie, true // 饼图
	}
	return seriesStyle, false
}

func ESPieChartV2(ctx context.Context, metrics map[string]map[string]map[string]float64) (*define.Chart, errors.Error) {
	chart := &define.Chart{}
	chart.Legend.Position = define.ChartLegendSPositionLeft
	chart.YAxis = []define.ChartYAxis{{Typ: define.AxisTypeValue, Position: "top", Format: ""}}
	var series []define.ChartSeriesItem
	for seriesName, terms := range metrics {
		for termName, values := range terms {
			for name, value := range values {
				if seriesName == name {
					// 不展示名字重复的值
					name = ""
				}
				series = append(series, define.ChartSeriesItem{
					Name:   buildSeriesName(seriesName, termName, name),
					Color:  "",
					Type:   define.ChartTypePie,
					Format: "",
					Power:  0,
					Data:   []*float64{SeriesValue(value)},
				})
			}
		}
	}
	chart.Series = define.ChartSeries{
		Stack: false,
		Items: series,
	}
	return chart, nil
}

func ESBarLineChartV2(ctx context.Context, seriesStyle string, metrics map[string]map[string]map[string]float64) (*define.Chart, errors.Error) {
	chart := &define.Chart{}
	chart.XAxis.Typ = define.AxisTypeCategory
	chart.Legend.Position = define.ChartLegendSPositionLeft
	chart.YAxis = []define.ChartYAxis{{Typ: define.AxisTypeValue, Position: "top", Format: ""}}

	// 生成x轴
	xAxis := make([]string, 0)
	xAxisRela := make(map[string]int, 0)
	for _, values := range metrics {
		for _, values := range values {
			for name := range values {
				if _, ok := xAxisRela[name]; !ok {
					xAxisRela[name] = len(xAxis)
					xAxis = append(xAxis, name)
				}
			}
		}
	}
	chart.XAxis.Data = xAxis
	var seriesArr []define.ChartSeriesItem

	for seriesName, values := range metrics {
		for termsName, values := range values {
			series := define.ChartSeriesItem{
				Type: seriesStyle,
				Name: buildSeriesName(seriesName, termsName, ""),
				Data: make([]*float64, len(xAxis)),
			}
			for name, value := range values {
				idx := xAxisRela[name]
				series.Data[idx] = SeriesValue(value)
			}
			seriesArr = append(seriesArr, series)
		}
	}

	chart.Series = define.ChartSeries{
		Stack: false,
		Items: seriesArr,
	}
	return chart, nil
}

func ESBarLineDateHistogramChartV2(ctx context.Context, cycle int8, seriesStyle string, date define.TimeRange, metrics map[string]map[string]map[string]float64) (*define.Chart, errors.Error) {
	chart := &define.Chart{}
	chart.XAxis.Typ = define.AxisTypeCategory
	chart.Legend.Position = define.ChartLegendSPositionLeft
	chart.YAxis = []define.ChartYAxis{{Typ: define.AxisTypeValue, Position: "top", Format: ""}}
	xaxis, times, err := tools.BuildTimeRangeCycleList(ctx, cycle, date)
	if err != nil {
		return nil, err
	}
	chart.XAxis.Data = xaxis

	chart.XAxis.Typ = define.AxisTypeCategory
	chart.Legend.Position = define.ChartLegendSPositionLeft
	chart.YAxis = []define.ChartYAxis{{Typ: define.AxisTypeValue, Position: "top", Format: ""}}
	xAxisCount := len(times)
	// 根据时间戳进行排序，决定用来显示x轴的位置
	timeIdxRela := make(map[string]int, xAxisCount)
	for idx, ts := range times {
		strDate := time.Unix(int64(ts), 0).Format("2006-01-02")
		timeIdxRela[strDate] = idx
	}
	seriesRela := make(map[string][]*float64, 0)
	for seriesName, seriesValues := range metrics {

		for termName, values := range seriesValues {
			// 不能省略，这里保证数据个数与x轴节点数相同
			// 将存储的字段，转换为展示的值
			if _, ok := seriesRela[seriesName]; !ok {
				seriesRela[buildSeriesName(seriesName, termName, "")] = make([]*float64, xAxisCount)
			}
			for name, value := range values {
				strDate := name
				vIdx, ok := timeIdxRela[strDate]
				if !ok {
					ctx.Log().InfoJSON("not found time in series index. values: %s, value index: %d, times: %s",
						seriesName, vIdx, times)
					continue
				}
				seriesRela[buildSeriesName(seriesName, termName, "")][vIdx] = SeriesValue(value)
			}
		}
	}

	var series []define.ChartSeriesItem
	extraSeries := make(map[string][]*float64, 0)
	for name, values := range seriesRela {
		extraSeries[name] = values
		series = append(series, define.ChartSeriesItem{
			Name:   name,
			Color:  "",
			Type:   seriesStyle,
			Format: "",
			Power:  0,
			Data:   values,
		})
	}
	chart.Series = define.ChartSeries{
		Stack: true,
		Items: series,
		Extra: define.ChartSeriesExtra{
			Data:       extraSeries,
			XAxisCount: len(times),
		},
	}

	return chart, nil

}

func buildSeriesName(seriesName, termName, name string) (rawSeriesName string) {
	nameParts := make([]string, 0, 3)
	if len(name) != 0 && name != define.AggFilterDefaultValueKey {
		nameParts = append(nameParts, name)
	}
	if termName != define.AggFilterDefaultValueKey {
		nameParts = append([]string{termName}, nameParts...)
	}
	if seriesName != define.AggDistinctName && seriesName != define.AggDateHistogramName {
		nameParts = append([]string{seriesName}, nameParts...)
	}
	rawSeriesName = strings.Join(nameParts, "_")
	if len(rawSeriesName) == 0 {
		rawSeriesName = ""
	}
	return rawSeriesName
}
