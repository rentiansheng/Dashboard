package tools

import (
	"github.com/rentiansheng/ges"
)

const (
	weekInterval       = "week"
	defaultDateFormat  = "yyyy-MM-dd"
	weekIntervalOffset = "+3d"
	defaultTimeZone    = "+08:00"
)

func DataHistogramCycle(cycle MetricCycleType, name, field string) ges.Agg {
	return ges.AggDataHistogramName(name, field, cycle.ESInterval(), defaultDateFormat, "", defaultTimeZone)
}

func DataHistogramSum(cycle MetricCycleType, name, field, sumName, sumField string) ges.Agg {
	return ges.AggDataHistogramSub(name, field, cycle.ESInterval(), defaultDateFormat, "", defaultTimeZone, ges.AggSum(sumName, sumField))
}

func DataHistogramAvg(cycle MetricCycleType, name, field, avgName, sumField string) ges.Agg {
	return ges.AggDataHistogramSub(name, field, cycle.ESInterval(), defaultDateFormat, "", defaultTimeZone, ges.AggAvg(avgName, sumField))
}
