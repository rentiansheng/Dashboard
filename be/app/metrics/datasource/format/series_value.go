package format

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"math"
)

func SeriesValueFormat(chartObj *define.Chart) {
	if chartObj == nil {
		return
	}
	for idx, series := range chartObj.Series.Items {
		for sIdx, s := range series.Data {
			if s == nil {
				continue
			}
			newS := *s * math.Pow10(int(series.Power))
			// 保留两位小数
			chartObj.Series.Items[idx].Data[sIdx] = SeriesValueTruncate(newS)
		}
	}
}

func SeriesValueAdd(a, b *float64) float64 {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return *b
	}
	if b == nil {
		return *a
	}
	return *a + *b
}

func SeriesValueAddPtr(a, b *float64) *float64 {

	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return SeriesValue(*a + *b)
}

func SeriesValue(f float64) *float64 {
	return &f
}

func SeriesValueTruncatePtr(f *float64) *float64 {
	if f == nil {
		return nil
	}
	return SeriesValue(math.Trunc(*f*100) / 100)
}

func SeriesValueTruncate(f float64) *float64 {

	return SeriesValue(math.Trunc(f*100) / 100)
}
