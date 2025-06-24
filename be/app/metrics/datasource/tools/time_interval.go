package tools

import (
	"fmt"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"sort"
	"time"
)

// 指标周期（1年，2季，3月，4周，5日，6时）
const (
	thousand          uint64 = 1000
	weekDayOfMonday          = 1
	weekDayOfThursday        = 4
	dayOfOneWeek             = 7

	endMonthOfQuarterOne   = 3
	endMonthOfQuarterTwo   = 6
	endMonthOfQuarterThree = 9

	quarterOne   = 1
	quarterTwo   = 2
	quarterThree = 3
	quarterFour  = 4

	halfYearOne = 6

	endTimeFromSecondToMillisecond uint64 = 999
)

type MetricCycleType = define.MetricCycleType

type TimeStampType string

const (
	Second      TimeStampType = "second"
	Millisecond TimeStampType = "millisecond"
)

type TimeInterval = define.TimeRange

type IntervalHandler func(ts uint64) (TimeInterval, error)

var intervalHandlerMap map[MetricCycleType]IntervalHandler

func init() {
	intervalHandlerMap = map[MetricCycleType]IntervalHandler{
		define.MetricCycleYear:     getYearInterval,
		define.MetricCycleQuarter:  getQuarterInterval,
		define.MetricCycleMonth:    getMonthInterval,
		define.MetricCycleDay:      getDayInterval,
		define.MetricCycleHour:     getHourInterval,
		define.MetricCycleWeek:     getNaturalWeekInterval,
		define.MetricCycleHalfYear: getHalfYearInterval,
	}
}

// GetTimeInterval 输入时间戳，输出时间戳所在周期的开始时间和结束时间
func GetTimeInterval(timestamp uint64, metricCycle MetricCycleType, timeStampType TimeStampType) (TimeInterval, error) {
	var res TimeInterval
	resList, err := GetTimeIntervalList([]uint64{timestamp}, metricCycle, timeStampType, true, true)
	if err != nil {
		return res, err
	}
	if len(resList) < 1 {
		return res, fmt.Errorf("unexpected err: result list of GetTimeIntervalList is empty. "+
			"timestamp: %v, metric cycle: %v, timestamp type: %v", timestamp, metricCycle, timeStampType)
	}
	return resList[0], nil
}

// GetTimeIntervalRangeBorder 根据输入开始结束时间和周期，得到基于周期开始时间，
//
//	通俗的讲就是start,end 是一个时间段， 在start，end在metricCycle 回划分未多个连续的区间，
//	但是start不是区间内第一个周期的开始时间， end 不是区间内最后一个周期结束时间。
//	所以需要调整到start是区间内第一个周期的开始时间，end是区间内最后一个周期结束时间
func GetTimeIntervalRangeBorder(ti TimeInterval, metricCycle MetricCycleType) (TimeInterval, error) {
	newTi := TimeInterval{}
	// 根据任务周期，对前端传过来的值，进行周期适配，保证开始和结束周期是周期开始和结束
	newStart, err := GetTimeInterval(ti.Start, metricCycle, Second)
	if err != nil {
		return newTi, err
	}
	newEnd, err := GetTimeInterval(ti.End, metricCycle, Second)
	if err != nil {
		return newTi, err
	}
	newTi = TimeInterval{
		Start: newStart.Start,
		End:   newEnd.End,
	}
	return newTi, nil
}

// GetTimeIntervalList 输入时间戳列表以及周期类型，输出对应周期的起始时间及结束时间
// timeStampType 输入及输出时间戳的类型, 不输入则默认为秒级时间戳
// needSort 输出的周期是否排序
// removeDuplicate 是否去除重复的周期
func GetTimeIntervalList(timeList []uint64, metricCycle MetricCycleType, timeStampType TimeStampType, needSort, removeDuplicate bool) ([]TimeInterval, error) {
	if timeStampType == Millisecond { // 转化为秒级时间戳处理
		for i := range timeList {
			timeList[i] = timeList[i] / thousand
		}
	}

	handler, ok := intervalHandlerMap[metricCycle]
	if !ok {
		return nil, fmt.Errorf("Unknown MetricCycle type, support year(1), quarter(2), month(3), week(4), day(5). ")
	}
	return getIntervalList(timeList, handler, timeStampType, needSort, removeDuplicate)
}

// TimeIntervalListByBeginEnd return TimeInterval list by start time and end time.
//
//	type TimeInterval struct {
//		Begin uint64
//		End   uint64
//	}
//
// For TimeInterval TI, it will be returned when TI meets one of following conditions:
// 1. TI.Start <= begin && TI.End >= begin
// 2. TI.Start >= begin && TI.Start <= end
func TimeIntervalListByBeginEnd(begin, end uint64, metricCycle MetricCycleType, timeStampType TimeStampType) ([]TimeInterval, error) {
	if timeStampType == Millisecond { // convert to second timestamp
		begin = begin / thousand
		end = end / thousand
	}

	res := make([]TimeInterval, 0)
	for t := begin; t <= end; {
		list, err := GetTimeIntervalList([]uint64{t}, metricCycle, Second, true, true)
		if err != nil {
			return nil, err
		}

		if len(list) != 1 {
			return nil, fmt.Errorf("TimeIntervalListByBeginEnd err, list returned by GetTimeIntervalList is empty")
		}

		res = append(res, list[0])
		t = list[0].End + 1
	}
	return res, nil
}

func TimeIntervalListByBeginCycleNums(begin uint64, cycleNums int, metricCycle MetricCycleType, timeStampType TimeStampType) ([]TimeInterval, error) {
	if cycleNums <= 0 {
		return nil, fmt.Errorf("get time interval list err, invalid cycle num: %v", cycleNums)
	}

	if timeStampType == Millisecond { // convert to second timestamp
		begin = begin / thousand
	}

	res := make([]TimeInterval, 0)
	t := begin
	for len(res) < cycleNums {
		list, err := GetTimeIntervalList([]uint64{t}, metricCycle, Second, true, true)
		if err != nil {
			return nil, err
		}

		if len(list) != 1 {
			return nil, fmt.Errorf("TimeIntervalListByBeginEnd err, list returned by GetTimeIntervalList is empty")
		}

		res = append(res, list[0])
		t = list[0].End + 1
	}

	return res, nil
}

// QuarterIntervalByYearAndQ 通过年和季度获取StartTime和EndTime
func QuarterIntervalByYearAndQ(year, quarter int) (TimeInterval, error) {
	return quarterIntervalByYearAndQ(year, quarter)
}

// HalfYearIntervalByYearAndH 获取年和1,2 half query 的时间区间
func HalfYearIntervalByYearAndH(year, seqNo int) (TimeInterval, error) {
	return halfYearIntervalByYearAndH(year, seqNo)
}

// sortTimeInterval 处理排序及去重
func sortTimeInterval(tiList []TimeInterval, timeStampType TimeStampType, needSort, removeDuplicate bool) []TimeInterval {
	if len(tiList) == 0 {
		return tiList
	}
	// 排序
	if needSort {
		sort.Slice(tiList, func(i, j int) bool {
			if tiList[i].Start < tiList[j].Start {
				return true
			}
			return false
		})
	}

	// 去重
	if removeDuplicate {
		temp := make([]TimeInterval, 0, len(tiList))
		resSet := make(map[uint64]struct{}, len(tiList))
		for _, ti := range tiList {
			_, ok := resSet[ti.Start]
			if !ok {
				temp = append(temp, ti)
				resSet[ti.Start] = struct{}{}
			}
		}
		tiList = temp
	}

	if timeStampType == Millisecond {
		for i := range tiList {
			tiList[i].Start = thousand * tiList[i].Start
			tiList[i].End = thousand*tiList[i].End + endTimeFromSecondToMillisecond
		}
	}
	return tiList
}

func getIntervalList(timeList []uint64, handler IntervalHandler, timeStampType TimeStampType, needSort, removeDuplicate bool) ([]TimeInterval, error) {
	res := make([]TimeInterval, 0, len(timeList))
	for _, timeStamp := range timeList {
		interval, err := handler(timeStamp)
		if err != nil {
			return nil, err
		}
		res = append(res, interval)
	}
	return sortTimeInterval(res, timeStampType, needSort, removeDuplicate), nil
}

func getYearInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	begin := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	end := time.Date(year+1, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(begin),
		End:   uint64(end),
	}, nil
}

func getQuarterInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	month := timeInput.Month()
	var quarter int
	switch {
	case month <= endMonthOfQuarterOne:
		quarter = quarterOne
	case month > endMonthOfQuarterOne && month <= endMonthOfQuarterTwo:
		quarter = quarterTwo
	case month > endMonthOfQuarterTwo && month <= endMonthOfQuarterThree:
		quarter = quarterThree
	default:
		quarter = quarterFour
	}
	return quarterIntervalByYearAndQ(year, quarter)
}

func getHalfYearInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	month := timeInput.Month()
	// 	 上半年
	endYear := year
	startMonth, endMonth := time.Month(1), time.Month(7)
	if month > halfYearOne {
		// 下半年
		endYear = year + 1
		startMonth, endMonth = 7, 1

	}

	// 天开始时间
	beginTime := time.Date(year, startMonth, 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	// 天结束时间
	endTime := time.Date(endYear, endMonth, 1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1

	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil

}

func quarterIntervalByYearAndQ(year, quarter int) (TimeInterval, error) {
	var res TimeInterval
	if year < 0 {
		return res, fmt.Errorf("invalid year: %v", year)
	}
	if quarter < 1 || quarter > 4 {
		return res, fmt.Errorf("invalid quarter: %v", quarter)
	}
	// 季度开始的时间
	beginMonth := quarter*3 - 2
	beginTime := time.Date(year, time.Month(beginMonth), 1, 0, 0, 0, 0, time.Now().Location()).Unix()

	// 季度结束的时间
	endMonth := quarter * 3
	endTime := time.Date(year, time.Month(endMonth+1), 1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

// halfYearIntervalByYearAndH 获取年和1,2 half query 的时间区间
func halfYearIntervalByYearAndH(year, seqNo int) (TimeInterval, error) {
	var res TimeInterval
	if year < 0 {
		return res, fmt.Errorf("invalid year: %v", year)
	}
	if seqNo < 1 || seqNo > 2 {
		return res, fmt.Errorf("invalid half year seqNo: %v", seqNo)
	}
	// endMoth 是下一个周期的开始月份，-1 是为了得到上一个周期的结束月份，结束周期是不固定的。 每个月的天数不一样
	startMoth, endMoth := time.Month(1), time.Month(7)
	endYear := year
	if seqNo == 2 {
		startMoth, endMoth = 7, 1
		year += 1
	}

	beginTime := time.Date(year, startMoth, 1, 0, 0, 0, 0, time.Now().Location()).Unix()

	endTime := time.Date(endYear, endMoth, 1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

func QuarterIntervalByYearAndQV2(year, quarter int) (TimeInterval, error) {
	var res TimeInterval
	if year < 0 {
		return res, fmt.Errorf("invalid year: %v", year)
	}
	if quarter < 1 || quarter > 4 {
		return res, fmt.Errorf("invalid quarter: %v", year)
	}
	// 季度开始的时间
	beginMonth := quarter*3 - 2
	beginTime := time.Date(year, time.Month(beginMonth), 1, 0, 0, 0, 0, time.Now().Location()).Unix()

	// 季度结束的时间
	endMonth := quarter * 3
	endTime := time.Date(year, time.Month(endMonth+1), 1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

func getMonthInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	month := timeInput.Month()

	// 月开始的时间
	beginTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location()).Unix()

	// 月结束的时间
	endTime := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

func getDayInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	month := timeInput.Month()
	day := timeInput.Day()

	// 天开始时间
	beginTime := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Unix()

	// 天结束时间
	endTime := time.Date(year, month, day+1, 0, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

func getHourInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	month := timeInput.Month()
	day := timeInput.Day()
	hour := timeInput.Hour()

	// 小时开始时间
	beginTime := time.Date(year, month, day, hour, 0, 0, 0, time.Now().Location()).Unix()

	// 小时结束时间
	endTime := time.Date(year, month, day, hour+1, 0, 0, 0, time.Now().Location()).Unix() - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

func getNaturalWeekInterval(ts uint64) (TimeInterval, error) {
	timeInput := time.Unix(int64(ts), 0)
	year := timeInput.Year()
	month := timeInput.Month()
	day := timeInput.Day()

	// 找到上一个周一
	lastMonday := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	weekDay := int(lastMonday.Weekday())
	dayToSub := weekDay - weekDayOfMonday
	if dayToSub < 0 {
		dayToSub += dayOfOneWeek
	}
	timeToSub := time.Duration(dayToSub*24) * time.Hour
	lastMonday = lastMonday.Add(-timeToSub)

	beginTime := lastMonday.Unix()
	endTime := beginTime + OneWeekTimeSecond - 1
	return TimeInterval{
		Start: uint64(beginTime),
		End:   uint64(endTime),
	}, nil
}

// GetRangeTimeIntervalList 根据输入的时间范围，返回所有周期开始和结束范围
func GetRangeTimeIntervalList(tsRange TimeInterval, cycleType MetricCycleType) ([]TimeInterval, error) {
	res := make([]TimeInterval, 0, 12)
	for newStart := tsRange.Start; newStart < tsRange.End; {
		item, err := GetTimeInterval(newStart, cycleType, Second)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
		// 转移到下一个周期
		newStart = item.End + 1
	}

	return res, nil

}
