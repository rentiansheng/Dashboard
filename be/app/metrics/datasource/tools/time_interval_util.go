package tools

import (
	"fmt"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"time"
)

const (
	firstDayOfMonth   = 1
	monthOfOneQuarter = 3
)

func GetDayInterval(now, dayToCount int64) (startTime, endTime uint64, err error) {
	timeStampToCount := now - dayToCount*OneDayTimeSecond
	timeIntervals, err := GetTimeIntervalList([]uint64{uint64(timeStampToCount)}, define.MetricCycleDay, Second, false, false)
	if err != nil {
		return 0, 0, err
	}
	if len(timeIntervals) != 1 {
		return 0, 0, fmt.Errorf("get time interval err, timeIntervals length: %v , expect length: 1", len(timeIntervals))
	}

	startTime = timeIntervals[0].Start
	endTime = timeIntervals[0].End
	return startTime, endTime, nil
}

func GetWeekInterval(now, weekToCount int64) (startTime, endTime uint64, err error) {
	nowT := time.Unix(now, 0)
	year := nowT.Year()
	month := nowT.Month()
	day := nowT.Day()
	// lastWeekT为上周某日的时间戳
	lastWeekT := time.Date(year, month, day-dayOfOneWeek*int(weekToCount), 0, 0, 0, 0, time.Now().Location()).Unix()

	timeIntervals, err := GetTimeIntervalList([]uint64{uint64(lastWeekT)}, define.MetricCycleWeek, Second, false, false)
	if err != nil {
		return 0, 0, err
	}
	if len(timeIntervals) != 1 {
		return 0, 0, fmt.Errorf("get time interval err, timeIntervals length: %v , expect length: 1", len(timeIntervals))
	}

	startTime = timeIntervals[0].Start
	endTime = timeIntervals[0].End
	return startTime, endTime, nil
}

func GetMonthInterval(now, monthToCount int64) (startTime, endTime uint64, err error) {
	nowT := time.Unix(now, 0)
	year := nowT.Year()
	month := nowT.Month()
	// monthT 为需要统计的季度某月1日的时间戳
	monthT := time.Date(year, month-time.Month(monthToCount), firstDayOfMonth, 0, 0, 0, 0,
		time.Now().Location()).Unix()

	timeIntervals, err := GetTimeIntervalList([]uint64{uint64(monthT)}, define.MetricCycleMonth, Second, false,
		false)
	if err != nil {
		return 0, 0, err
	}
	if len(timeIntervals) != 1 {
		return 0, 0, fmt.Errorf("get time interval err, timeIntervals length: %v , "+
			"expect length: 1", len(timeIntervals))
	}

	startTime = timeIntervals[0].Start
	endTime = timeIntervals[0].End
	return startTime, endTime, nil
}

func GetQuarterInterval(now, quarterToCount int64) (startTime, endTime uint64, err error) {
	nowT := time.Unix(int64(now), 0)
	year := nowT.Year()
	month := nowT.Month()
	// quarterT为上个季度某个月1日的时间戳
	quarterT := time.Date(year, month-monthOfOneQuarter*time.Month(quarterToCount), firstDayOfMonth, 0,
		0, 0, 0, time.Now().Location()).Unix()

	timeIntervals, err := GetTimeIntervalList([]uint64{uint64(quarterT)}, define.MetricCycleQuarter, Second,
		false, false)
	if err != nil {
		return 0, 0, err
	}
	if len(timeIntervals) != 1 {
		return 0, 0, fmt.Errorf("get time interval err, timeIntervals length: %v , "+
			"expect length: 1", len(timeIntervals))
	}

	startTime = timeIntervals[0].Start
	endTime = timeIntervals[0].End
	return startTime, endTime, nil
}

func GetYearInterval(now, yearToCount int64) (startTime, endTime uint64, err error) {
	nowT := time.Unix(now, 0)
	year := nowT.Year()
	month := nowT.Month()
	day := nowT.Day()
	// lastYearT为去年的时间戳
	lastYearT := time.Date(year-int(yearToCount), month, day, 0, 0, 0, 0, time.Now().Location()).Unix()

	timeIntervals, err := GetTimeIntervalList([]uint64{uint64(lastYearT)}, define.MetricCycleYear,
		Second, false, false)
	if err != nil {
		return 0, 0, err
	}
	if len(timeIntervals) != 1 {
		return 0, 0, fmt.Errorf("get time interval err, timeIntervals length: %v , "+
			"expect length: 1", len(timeIntervals))
	}

	startTime = timeIntervals[0].Start
	endTime = timeIntervals[0].End
	return startTime, endTime, nil
}

func GetWeekdaysTime(startTime, endTime string) (int64, error) {
	totalTime := int64(0)
	start, err := GetLocTimeByStrDate(startTime)
	if err != nil {
		return 0, err
	}
	end, err := GetLocTimeByStrDate(endTime)
	if err != nil {
		return 0, err
	}
	if start.Unix() >= end.Unix() {
		return 0, nil
	}
	startZeroTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endZeroTime := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	tomorrow := start.AddDate(0, 0, 1)
	tomorrowZeroTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0,
		tomorrow.Location())
	if startZeroTime.Unix() == endZeroTime.Unix() {
		if isWeekend(start) {
			return 0, nil
		} else {
			totalTime = end.Unix() - start.Unix()
			return totalTime, nil
		}
	}
	if !isWeekend(start) {
		totalTime += tomorrowZeroTime.Unix() - start.Unix()
	}
	for {
		if tomorrowZeroTime.Unix() >= endZeroTime.Unix() {
			break
		}
		if !isWeekend(tomorrowZeroTime) {
			totalTime += 24 * 60 * 60
		}
		tomorrowZeroTime = tomorrowZeroTime.Add(24 * time.Hour)
	}
	if !isWeekend(endZeroTime) {
		totalTime += end.Unix() - endZeroTime.Unix()
	}
	return totalTime, nil
}

// GetWeekendDaysListByTimestamp 获取时间戳区间内非工作日，
func GetWeekendDaysListByTimestamp(timeRange TimeInterval) []TimeInterval {
	result := make([]TimeInterval, 0)
	endTime := int64(timeRange.End)
	start, end := time.Unix(int64(timeRange.Start), 0), time.Unix(endTime, 0)

	if start.Unix() >= end.Unix() {
		return nil
	}
	dayStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	for dayStart.Unix() < endTime {
		switch dayStart.Weekday() {
		case time.Saturday:
			end := dayStart.Add(48 * time.Hour)
			result = append(result, TimeInterval{
				Start: uint64(dayStart.Unix()),
				// 周六开始，周日结束
				End: uint64(end.Unix()) - 1,
			})
			dayStart = end
		case time.Sunday:
			end := dayStart.Add(24 * time.Hour)
			result = append(result, TimeInterval{
				Start: uint64(dayStart.Unix()),
				// 周日开始，周一结束
				End: uint64(end.Unix()) - 1,
			})
			dayStart = end

		default:
			// 判断下一天是否是周末
			dayStart = dayStart.Add(24 * time.Hour)

		}

	}

	return result

}

func isWeekend(weekTime time.Time) bool {
	return weekTime.Weekday() == time.Saturday || weekTime.Weekday() == time.Sunday
}
