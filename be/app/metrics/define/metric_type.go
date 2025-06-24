package define

/***************************
    @author: tiansheng.ren
    @date: 2025/6/9
    @desc:

***************************/

type MetricCycleType uint8

const (
	MetricCycleYear     MetricCycleType = 1
	MetricCycleQuarter                  = 2
	MetricCycleMonth                    = 3
	MetricCycleWeek                     = 4
	MetricCycleDay                      = 5
	MetricCycleHour                     = 6
	MetricCycleHalfYear                 = 8
)

func (m MetricCycleType) ESInterval() string {
	switch m {
	case MetricCycleYear:
		return "year"
	case MetricCycleQuarter:
		return "quarter"
	case MetricCycleMonth:
		return "month"
	case MetricCycleWeek:
		return "week"
	case MetricCycleDay:
		return "day"
	case MetricCycleHour:
		return "hour"
	case MetricCycleHalfYear:
		return "6M"
	}
	return "unknown"
}

func (m MetricCycleType) String() string {
	switch m {
	case MetricCycleYear:
		return "year"
	case MetricCycleQuarter:
		return "quarter"
	case MetricCycleMonth:
		return "month"
	case MetricCycleWeek:
		return "week"
	case MetricCycleDay:
		return "day"
	case MetricCycleHour:
		return "hour"
	case MetricCycleHalfYear:
		return "half year"
	}
	return "unknown"
}
