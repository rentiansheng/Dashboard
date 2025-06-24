package tools

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var timeZoneID = map[string]string{
	RegionSG:                  "Singapore",
	strings.ToUpper(RegionSG): "Singapore",
	RegionVN:                  "Asia/Ho_Chi_Minh",
	strings.ToUpper(RegionVN): "Asia/Ho_Chi_Minh",
	RegionID:                  "Asia/Jakarta",
	strings.ToUpper(RegionID): "Asia/Jakarta",
	RegioTH:                   "Asia/Bangkok",
	strings.ToUpper(RegioTH):  "Asia/Bangkok",
	RegionTW:                  "Asia/Taipei",
	strings.ToUpper(RegionTW): "Asia/Taipei",
	RegionMY:                  "Asia/Kuala_Lumpur",
	strings.ToUpper(RegionMY): "Asia/Kuala_Lumpur",
	RegionPH:                  "Asia/Manila",
	strings.ToUpper(RegionPH): "Asia/Manila",
	RegionCN:                  "Asia/Shanghai",
	strings.ToUpper(RegionCN): "Asia/Shanghai",
}
var (
	RegionSG = "sg"
	RegionVN = "vn"
	RegionID = "id"
	RegioTH  = "th"
	RegionTW = "tw"
	RegionMY = "my"
	RegionPH = "ph"
	RegionCN = "cn"

	OneDayTimeSecond  = int64(24 * 60 * 60)
	OneWeekTimeSecond = 7 * OneDayTimeSecond
	TwoWeekTimeSecond = 2 * OneWeekTimeSecond
)

func GetStartAndLastTimestampOfMonth(monthStr, region string) (int64, int64, error) {
	tz, err := getTimeZoneID(region)
	if err != nil {
		return 0, 0, err
	}
	loc, _ := time.LoadLocation(tz)
	t, _ := time.ParseInLocation("2006-01", monthStr, loc)
	last := t.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return t.Unix(), last.Unix(), nil
}

func TimeStringParse(timeString, format, region string) (int64, error) {
	if timeString == "" {
		return 0, nil
	}
	tz, err := getTimeZoneID(region)
	if err != nil {
		return 0, err
	}
	loc, _ := time.LoadLocation(tz)
	t, _ := time.ParseInLocation(format, timeString, loc)
	return t.Unix(), nil
}

func TimeStringFormat(timestamp int64, format, region string) (string, error) {
	if timestamp == 0 {
		return "", nil
	}
	t := time.Unix(timestamp, 0)
	tz, err := getTimeZoneID(region)
	if err != nil {
		return "", err
	}
	loc, _ := time.LoadLocation(tz)
	t = t.In(loc)
	tStr := t.Format(format)
	return tStr, nil
}

func getTimeZoneID(region string) (string, error) {
	if zone, ok := timeZoneID[region]; ok {
		return zone, nil
	}
	return "", fmt.Errorf("invalid region")
}

func NowTimestamp() int32 {
	return (int32)(time.Now().Unix())
}

// nowt : seconds
func GetTimestampMonthFirstDateByTime(nowt int64, region string) (time.Time, error) {
	tz, err := getTimeZoneID(region)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(nowt, 0)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	t = t.In(loc)
	t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return t, nil
}

func GetTimestampOfNextWeekMonday(region string) (int64, error) {
	tz, err := getTimeZoneID(region)
	if err != nil {
		return 0, err
	}
	t := time.Now()
	loc, _ := time.LoadLocation(tz)
	t = t.In(loc)
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, offset+7)
	return t.Unix(), nil
}

func GetTimestampOfNexDay(region string) (int64, error) {
	tz, err := getTimeZoneID(region)
	if err != nil {
		return 0, err
	}
	t := time.Now()
	loc, _ := time.LoadLocation(tz)
	t = t.In(loc)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, 1)
	return t.Unix(), nil
}

func GetTimestampOfNextWeekMondayBytime(nowt int64, region string) (*time.Time, error) {
	t := time.Unix(nowt, 0)
	tz, err := getTimeZoneID(region)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation(tz)
	t = t.In(loc)
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, offset+7)
	return &t, nil
}

func GetLastWeekRange(nowt int64, region string) (start int64, end int64, err error) {
	return GetLastWeeksRange(1, nowt, region)
}

func GetLastWeeksRange(week int, nowt int64, region string) (start int64, end int64, err error) {
	nextWeekMonDay, err := GetTimestampOfNextWeekMondayBytime(nowt, region)
	if err != nil {
		return 0, 0, err
	}
	thisWeekMonDay := nextWeekMonDay.AddDate(0, 0, -7)
	lastWeekMonDay := thisWeekMonDay.AddDate(0, 0, -7*week)
	return lastWeekMonDay.Unix(), thisWeekMonDay.Unix() - 1, nil
}

func GetLastMonthRange(nowt int64, region string) (start int64, end int64, err error) {
	return GetLastMonthsRange(1, nowt, region)
}

func GetLastMonthsRange(months int, nowt int64, region string) (start int64, end int64, err error) {
	thisMonthFirstDay, err := GetTimestampMonthFirstDateByTime(nowt, region)
	if err != nil {
		return 0, 0, err
	}
	lastMonthsFirstDay := thisMonthFirstDay.AddDate(0, -1*months, 0)
	return lastMonthsFirstDay.Unix(), thisMonthFirstDay.Unix() - 1, nil
}

func GetDayRange(nowt int64, region string) (start int64, end int64, err error) {
	tz, err := getTimeZoneID(region)
	if err != nil {
		return 0, 0, err
	}
	t := time.Unix(nowt, 0)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return 0, 0, err
	}
	t = t.In(loc)
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	start = today.Unix()
	end = start + OneDayTimeSecond - 1
	return start, end, nil
}

func LastWeekRangeByEffectiveTime(ts, effectiveTime int64, region string) (start, end int64, err error) {
	start, end, err = GetLastWeekRange(ts, region)
	if err != nil {
		return 0, 0, err
	}
	if end+1 < effectiveTime {
		return 0, 0, nil
	}
	return start, end, err
}

func LastBiweekRangeByEffectiveTime(ts, effectiveTime int64, region string) (start, end int64, err error) {
	start, end, err = GetLastWeeksRange(2, ts, region)
	if err != nil {
		return 0, 0, err
	}
	if end+1 < effectiveTime {
		return 0, 0, nil
	}
	if (end-effectiveTime)%TwoWeekTimeSecond != 0 {
		start -= OneWeekTimeSecond
		end -= OneWeekTimeSecond
	}
	if end+1 < effectiveTime {
		return 0, 0, nil
	}
	return start, end, err
}

func LastMonthRangeByEffectiveTime(ts, effectiveTime int64, region string) (start, end int64, err error) {
	start, end, err = GetLastMonthRange(ts, region)
	if err != nil {
		return 0, 0, err
	}
	if end+1 < effectiveTime {
		return 0, 0, nil
	}
	return start, end, err
}
func LastDayRangeByEffectiveTime(ts, effectiveTime int64, region string) (start, end int64, err error) {
	start, end, err = GetDayRange(ts-OneDayTimeSecond, region)
	if err != nil {
		return 0, 0, err
	}
	if end+1 < effectiveTime {
		return 0, 0, nil
	}
	return start, end, err
}

func FormatToYYYYMMDD(t time.Time) string {
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

// GetTimeOfThisWeekMondayByUnix 根据nowt来获取所属星期的星期一
func GetTimeOfThisWeekMondayByUnix(nowt int64, region string) (time.Time, error) {
	t := time.Unix(nowt, 0)
	tz, err := getTimeZoneID(region)
	if err != nil {
		return t, err
	}

	loc, _ := time.LoadLocation(tz)
	t = t.In(loc)
	offset := int(time.Monday - t.Weekday())

	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, offset)
	return t, nil
}

// GetTimeByTimeString 根据 格式化的时间 获取time.time
func GetTimeByTimeString(timeString, format, region string) (time.Time, error) {
	if timeString == "" {
		return time.Time{}, nil
	}
	tz, err := getTimeZoneID(region)
	if err != nil {
		return time.Time{}, err
	}
	loc, _ := time.LoadLocation(tz)
	t, _ := time.ParseInLocation(format, timeString, loc)
	return t, nil
}

// ValidateTimeString 校验时间格式是否合法
func ValidateTimeString(timeString, format, region string) error {

	ts, err := TimeStringParse(timeString, format, region)
	if err != nil {
		return fmt.Errorf("failed to TimeStringParse, time:%s, format:%s, region:%s, err:%v", timeString, format, region, err)
	}

	newTimeString, err := TimeStringFormat(ts, format, region)
	if err != nil {
		return fmt.Errorf("failed to TimeStringFormat, ts:%d, format:%s, region:%s, err:%v", ts, format, region, err)
	}

	if timeString != newTimeString {
		return fmt.Errorf("time:%s is invalid", timeString)
	}

	return nil
}

// GetTimeByTimestamp 根据时间戳获取 time.Time
func GetTimeByTimestamp(ts int64, region string) (time.Time, error) {
	t := time.Unix(ts, 0)

	tz, err := getTimeZoneID(region)
	if err != nil {
		return time.Time{}, err
	}

	// time zone
	loc, _ := time.LoadLocation(tz)
	t = t.In(loc)

	return t, nil
}

// ConvertTimeStr convert date str to yyyy-mm-dd
func ConvertTimeStr(in, inFormat, outFormat string) (string, error) {
	// already satisfy output format
	ts, err := time.Parse(outFormat, in)
	if err == nil {
		return in, nil
	}
	ts, err = time.Parse(inFormat, in)
	if err != nil {
		return "", fmt.Errorf("failed to parse, time:%s, format:%s, err:%v", in, inFormat, err)
	}
	res := ts.Format(outFormat)
	return res, nil

}

func WeekByTimeStamp(timeUnix int64) (uWeek, uYear uint64, err error) {
	timeInput := time.Unix(timeUnix, 0)
	year := timeInput.Year()

	firstThursday := firstThursdayOfYear(year)
	// 说明timeInput所在的周属于上一年
	if timeInput.Unix() < firstThursday.Unix() {
		year -= 1
		firstThursday = firstThursdayOfYear(year)
	}
	timeDuration := timeInput.Unix() - firstThursday.Unix()
	week := int(timeDuration/OneWeekTimeSecond) + 1
	return uint64(week), uint64(year), nil
}

// firstThursdayOfYear 找到该年的第一个星期四
func firstThursdayOfYear(year int) time.Time {
	firstThursday := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location())
	dayToAdd := 4 - int(firstThursday.Weekday())
	if dayToAdd < 0 {
		dayToAdd += 7
	}
	timeToAdd := time.Duration(dayToAdd*24) * time.Hour
	firstThursday = firstThursday.Add(timeToAdd)

	return firstThursday
}

func QuarterByTimeStamp(timeUnix int64) (quarter, year int) {
	timeInput := time.Unix(timeUnix, 0)
	year = timeInput.Year()

	month := timeInput.Month()
	quarter = int(month/3) + 1
	return quarter, year
}

/*
通过时间戳计算所属周的时间范围
execType: 0:计算当前周期，1:计算上一次周期, 其他:计算当前周期
*/
func GetWeekTimeByTimeStamp(timeUnix int64, execType int) (int64, int64, error) {
	var (
		startTime, endTime time.Time
		err                error
	)
	timeFormat := "2006-01-02"
	timeTime := time.Unix(timeUnix, 0)
	diffDay := int(timeTime.Weekday() - 4)
	var thursDayTime string
	if diffDay >= 0 {
		thursDayTime = timeTime.AddDate(0, 0, -diffDay).Format(timeFormat)
	} else {
		thursDayTime = timeTime.AddDate(0, 0, -(diffDay + 7)).Format(timeFormat)
	}
	switch execType {
	case 0:
		startTime, err = time.ParseInLocation(timeFormat, thursDayTime, time.Local)
		if err != nil {
			return 0, 0, err
		}
		endTime = startTime.AddDate(0, 0, 7)
	case 1:
		endTime, err = time.ParseInLocation(timeFormat, thursDayTime, time.Local)
		if err != nil {
			return 0, 0, err
		}
		startTime = endTime.AddDate(0, 0, -7)
	default:
		now, err := time.ParseInLocation(timeFormat, thursDayTime, time.Local)
		if err != nil {
			return 0, 0, err
		}
		endTime = now.AddDate(0, 0, -(7 * (execType - 1)))
		startTime = now.AddDate(0, 0, -(7 * execType))
	}
	return startTime.Unix(), endTime.Unix() - 1, nil
}

/*
通过时间戳计算所属日的时间范围
execType: 0:计算当前周期，1:计算上一次周期, 其他:计算当前周期
*/
func GetDayTimeByTimeStamp(timeUnix int64, execType int) (int64, int64, error) {
	var startDay, endDay string
	timeFormat := "2006-01-02"
	timeTime := time.Unix(timeUnix, 0)
	switch execType {
	case 0:
		startDay = timeTime.Format(timeFormat)
		endDay = timeTime.AddDate(0, 0, 1).Format(timeFormat)
	case 1:
		startDay = timeTime.AddDate(0, 0, -1).Format(timeFormat)
		endDay = timeTime.Format(timeFormat)
	default:
		startDay = timeTime.AddDate(0, 0, -execType).Format(timeFormat)
		endDay = timeTime.Format(timeFormat)
	}
	startTime, err := TimeStringParse(startDay, timeFormat, RegionSG)
	if err != nil {
		return 0, 0, err
	}
	endTime, err := TimeStringParse(endDay, timeFormat, RegionSG)
	if err != nil {
		return 0, 0, err
	}
	return startTime, endTime - 1, nil
}

// GetMonthTimeByTimeStamp 通过时间戳计算所属月的时间范围
// execType: 0:计算当前周期，1:计算上一次周期, 其他:不支持
func GetMonthTimeByTimeStamp(timeUnix int64, execType int) (int64, int64, error) {
	var startMonth, endMonth string
	timeLayout := "2006-01"
	timeTime := time.Unix(timeUnix, 0)
	switch execType {
	case 0:
		startMonth = timeTime.Format(timeLayout)
		endMonth = timeTime.AddDate(0, 1, 0).Format(timeLayout)
	case 1:
		startMonth = timeTime.AddDate(0, -1, 0).Format(timeLayout)
		endMonth = timeTime.Format(timeLayout)
	// AddDate跨月时会计算错误，所以不支持其他类型
	default:
		return 0, 0, fmt.Errorf("not support this execType, please input 0 or 1")
	}
	startMonthStamp, err := time.ParseInLocation(timeLayout, startMonth, time.Local)
	if err != nil {
		return 0, 0, err
	}
	endMonthStamp, err := time.ParseInLocation(timeLayout, endMonth, time.Local)
	if err != nil {
		return 0, 0, err
	}
	return startMonthStamp.Unix(), endMonthStamp.Unix() - 1, nil
}

/*
通过时间戳计算所属季度的时间范围
execType: 0:计算当前周期，1:计算上一次周期, 其他:计算当前周期
*/
func GetQuarterTimeByTimeStamp(timeUnix int64, execType int) (int64, int64, error) {
	var startQuarterStartMonth, endQuarterStartMonth string
	timeLayout := "2006-01"
	timeTime := time.Unix(timeUnix, 0)
	month := timeTime.Month()
	diffMonth := month % 3
	switch execType {
	case 0:
		switch diffMonth {
		case 0:
			startQuarterStartMonth = timeTime.AddDate(0, -2, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.Format(timeLayout)
		case 1:
			startQuarterStartMonth = timeTime.Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, 2, 0).Format(timeLayout)
		case 2:
			startQuarterStartMonth = timeTime.AddDate(0, -1, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, 1, 0).Format(timeLayout)
		}
	case 1:
		diffMonth := month % 3
		switch diffMonth {
		case 0:
			startQuarterStartMonth = timeTime.AddDate(0, -5, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, -2, 0).Format(timeLayout)
		case 1:
			startQuarterStartMonth = timeTime.AddDate(0, -3, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, 0, 0).Format(timeLayout)
		case 2:
			startQuarterStartMonth = timeTime.AddDate(0, -4, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, -1, 0).Format(timeLayout)
		}
	default:
		switch diffMonth {
		case 0:
			startQuarterStartMonth = timeTime.AddDate(0, -2, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.Format(timeLayout)
		case 1:
			startQuarterStartMonth = timeTime.Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, 2, 0).Format(timeLayout)
		case 2:
			startQuarterStartMonth = timeTime.AddDate(0, -1, 0).Format(timeLayout)
			endQuarterStartMonth = timeTime.AddDate(0, 1, 0).Format(timeLayout)
		}
	}
	startQuarterStamp, err := time.ParseInLocation(timeLayout, startQuarterStartMonth, time.Local)
	if err != nil {
		return 0, 0, err
	}
	endQuarterStamp, err := time.ParseInLocation(timeLayout, endQuarterStartMonth, time.Local)
	if err != nil {
		return 0, 0, err
	}
	return startQuarterStamp.Unix(), endQuarterStamp.Unix() - 1, nil
}

/*
通过时间戳计算所属年的时间范围
execType: 0:计算当前周期，1:计算上一次周期, 其他:计算当前周期
*/
func GetYearTimeByTimeStamp(timeUnix int64, execType int) (int64, int64, error) {
	var startYear, endYear string
	timeLayout := "2006"
	timeTime := time.Unix(timeUnix, 0)
	switch execType {
	case 0:
		startYear = timeTime.Format(timeLayout)
		endYear = timeTime.AddDate(1, 0, 0).Format(timeLayout)
	case 1:
		startYear = timeTime.AddDate(-1, 0, 0).Format(timeLayout)
		endYear = timeTime.Format(timeLayout)
	default:
		startYear = timeTime.AddDate(-execType, 0, 0).Format(timeLayout)
		endYear = timeTime.Format(timeLayout)
	}
	endYearStamp, err := time.ParseInLocation(timeLayout, endYear, time.Local)
	if err != nil {
		return 0, 0, err
	}
	startYearStamp, err := time.ParseInLocation(timeLayout, startYear, time.Local)
	if err != nil {
		return 0, 0, err
	}
	return startYearStamp.Unix(), endYearStamp.Unix() - 1, nil
}

// GetDayStartTimestamp 通过一天当中任何时间戳获取当日最开始时的时间戳
func GetDayStartTimestamp(nowt int64, region string) (start int64, err error) {
	tz, err := getTimeZoneID(region)
	if err != nil {
		return 0, err
	}
	t := time.Unix(nowt, 0)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return 0, err
	}
	t = t.In(loc)
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return today.Unix(), nil
}

// GetStartAndEndOfPastNMonthFromNow 获取过去n个月开始到现在的时间戳范围: 前n个月 ～ now
func GetStartAndEndOfPastNMonthFromNow(monthNum int, region string) (int64, int64, error) {
	now := time.Now().Unix()
	start, _, err := GetLastMonthsRange(monthNum, now, region)
	if err != nil {
		return 0, 0, err
	}
	return start, now, nil
}

// GetTimestampOfPastNQuarterFromNow 获取过去n个季度开始到现在的时间戳: 前n个季度 ～ now
func GetTimestampOfPastNQuarterFromNow(quarterNum int, region string) (int64, int64, error) {
	now := time.Now()
	thisMonthFirstDay, err := GetTimestampMonthFirstDateByTime(now.Unix(), region)
	if err != nil {
		return 0, 0, err
	}
	pathMonthsFirstDay := thisMonthFirstDay.AddDate(0, -1*quarterNum*3, 0)
	return pathMonthsFirstDay.Unix(), now.Unix(), nil
}

func GetMonthByTimestamp(timestamp int64) (int8, int16, error) {
	timeInput := time.Unix(timestamp, 0)
	year := timeInput.Year()
	month := timeInput.Month()
	return int8(month), int16(year), nil
}

func GetQuarterByMonthNum(month int8) int8 {
	return int8(int((month-1)/3) + 1)
}

func GetTimestampOfPastNYearFromNow(yearNum int) (int64, int64) {
	now := time.Now()
	first := time.Date(now.Year()-yearNum, time.January, 1, 0, 0, 0, 0, now.Location())
	return first.Unix(), now.Unix()
}

func GetContinuousMonthNums(startMonth, endMonth int8, startYear, endYear int16) ([]int8, error) {
	result := make([]int8, 0)
	start := startMonth
	end := endMonth + 12*(int8)(endYear-startYear)
	if end < start {
		return nil, errors.New("start and end quarter set error")
	}
	for start <= end {
		temp := start % 12
		if temp == 0 {
			temp = 12
		}
		result = append(result, temp)
		start += 1
	}
	return result, nil
}

func GetContinuousQuarterNums(startQuarter, endQuarter int8, startYear, endYear int16) ([]int8, error) {
	result := make([]int8, 0)
	start := startQuarter
	end := endQuarter + 4*(int8)(endYear-startYear)

	if end < start {
		return nil, errors.New("start and end quarter set error")
	}

	for start <= end {
		temp := start % 4
		if temp == 0 {
			temp = 4
		}
		result = append(result, temp)
		start += 1
	}
	return result, nil
}

// GetLastTime 获取过去几年/季度/月/周的时间轴(包括跨年)
func GetLastTime(timeUnit string, amount int64, year int) ([]string, []string) {
	var (
		timeCount          []int
		timeCountLastYear  []int
		currentTimeNum     uint64
		totalCountLastYear uint64
	)

	switch timeUnit {
	case "week":
		currentTimeNum, _, _ = WeekByTimeStamp(time.Now().Unix())
		lastYearDay := time.Date(year, 12, 31, 12, 00, 00, 0, time.UTC)
		totalCountLastYear, _, _ = WeekByTimeStamp(lastYearDay.Unix())
	case "month":
		currentTimeNum = uint64(time.Now().Month())
		totalCountLastYear = 12
	case "quarter":
		timeNum, _ := QuarterByTimeStamp(time.Now().Unix())
		currentTimeNum = uint64(timeNum)
		totalCountLastYear = 4
	case "year":
		for i := year; amount > 0; i-- {
			timeCount = append(timeCount, i)

			amount--
		}
		sort.Ints(timeCount)
		timeCountStr := MapToSlice(timeCount, func(s int) string {
			return strconv.Itoa(s)
		})
		return timeCountStr, nil
	}
	if amount <= int64(currentTimeNum) {
		// 如果没有跨年
		for i := currentTimeNum; amount > 0; i-- {
			timeCount = append(timeCount, int(i))
			amount--
		}
	} else {
		for i := currentTimeNum; i > 0; i-- {
			timeCount = append(timeCount, int(i))
			amount--
		}
		for i := totalCountLastYear; amount > 0; i-- {
			timeCountLastYear = append(timeCountLastYear, int(i))
			amount--
		}
	}
	sort.Ints(timeCount)
	sort.Ints(timeCountLastYear)
	timeCountStr := MapToSlice(timeCount, func(s int) string {
		return strconv.Itoa(s)
	})
	timeCountLastYearStr := MapToSlice(timeCountLastYear, func(s int) string {
		return strconv.Itoa(s)
	})
	return timeCountStr, timeCountLastYearStr
}

func GetContinuesUnits(start, end int64, timeUnit string) []string {
	result := make([]string, 0)
	if start > end {
		return result
	}
	startDate := time.Unix(start, 0)
	endDate := time.Unix(end, 0)
	switch timeUnit {
	case "year":
		format := "2006"
		aimEndYear := endDate.Format(format)
		var endYear string
		for index := 0; endYear != aimEndYear; index++ {
			endYear = startDate.AddDate(index, 0, 0).Format(format)
			result = append(result, endYear)
		}
		return result
	case "month":
		format := "06M01"
		aimEndMonth := endDate.Format(format)
		var endMonth string
		for index := 0; endMonth != aimEndMonth; index++ {
			endMonth = startDate.AddDate(0, index, 0).Format(format)
			result = append(result, endMonth)
		}
		return result
	case "week":
		aimEndWeek := GetWeek(end)
		var endWeek string
		var index int64
		for index = 0; endWeek != aimEndWeek; index++ {
			temp := start + OneWeekTimeSecond*index
			endWeek = GetWeek(temp)
			result = append(result, endWeek)
		}
		return result
	}
	return result
}

// GetWeek 该方法取到的周与数据库中取到的周不一样
func GetWeek(timestamp int64) string {
	datetime := time.Unix(timestamp, 0).Format(time.ANSIC)
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(time.ANSIC, datetime, loc)
	year, week := tmp.ISOWeek()

	return fmt.Sprintf("%dW%02d", year, week-1)
}

func GetDiffDays(t1, t2 time.Time) string {
	return strconv.Itoa(int(t1.Sub(t2).Hours() / 24))
}

func TransferStringToDuration(s string) time.Duration {
	var D time.Duration
	if s != "" {
		D, _ = time.ParseDuration(s)
	}
	D, _ = time.ParseDuration("0s")
	return D
}

func GetWeekTimeRanges(start, end int64) [][]int64 {
	ranges := make([][]int64, 0)

	nextMonday, _ := GetTimestampOfNextWeekMondayBytime(start, RegionSG)
	tempMonday := nextMonday.Unix() - OneWeekTimeSecond
	tempSunday := nextMonday.Unix() - 1
	for tempSunday < end {
		ranges = append(ranges, []int64{tempMonday, tempSunday})
		tempMonday = tempSunday + 1
		tempSunday = tempSunday + OneWeekTimeSecond
	}
	return ranges
}

// GetLocTimeByStrDate 获取本地时间根据字符串 eg:str=2006-01-02 15:04:05
func GetLocTimeByStrDate(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}
	const timeParesLayout = "2006-01-02 15:04:05"
	return GetTimeByTimeString(str, timeParesLayout, RegionCN)
}

// GetLocTimestampByStrDate 获取本地时间戳根据字符串 eg:str=2006-01-02 15:04:05
func GetLocTimestampByStrDate(str string) (int64, error) {
	t, err := GetLocTimeByStrDate(str)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// GetLocTimestampByStrDateUint64 获取本地时间戳根据字符串 eg:str=2006-01-02 15:04:05
func GetLocTimestampByStrDateUint64(str string) (uint64, error) {
	t, err := GetLocTimestampByStrDate(str)
	if err != nil {
		return 0, err
	}
	return uint64(t), nil
}

const (
	yearSeconds  = 365 * 24 * 3600
	monthSeconds = 30 * 24 * 3600
	daySeconds   = 24 * 3600
	hourSeconds  = 3600
	minSeconds   = 60
	secSeconds   = 1

	yearStr   = "year"
	monthStr  = "mth"
	dayStr    = "d"
	hourStr   = "hour"
	minStr    = "min"
	secondStr = "sec"

	dateIntervalLength = 2
)

var timeSecondList = []uint64{yearSeconds, monthSeconds, daySeconds, hourSeconds, minSeconds, secSeconds}
var timeStrList = []string{yearStr, monthStr, dayStr, hourStr, minStr, secondStr}

// Second2DateInterval input seconds, output date interval:
// input: 65, output: 1min 5s
// input: 3726, output: 1hour 2min
func Second2DateInterval(sec uint64) string {
	res := make([]string, 0)
	for i := range timeSecondList {
		intervalSec := timeSecondList[i]
		intervalStr := timeStrList[i]

		val := sec / intervalSec
		sec -= val * intervalSec
		if val > 0 || ((i == len(timeSecondList)-1) && len(res) == 0) {
			res = append(res, strconv.Itoa(int(val))+intervalStr)
		} else if len(res) > 0 {
			break
		}

		if len(res) >= dateIntervalLength {
			break
		}
	}

	return strings.Join(res, " ")
}

const jiraTimeLayout = "2006-01-02T15:04:05"

// ConvertJiraTime input 2022-10-20T14:17:09.000+0800, output: second timestamp
func ConvertJiraTime(jiraTime string) (int64, error) {
	if jiraTime == "" {
		return 0, nil
	}

	jiraTime = strings.Split(jiraTime, ".")[0]
	timeObj, err := time.Parse(jiraTimeLayout, jiraTime)
	if err != nil {
		return 0, err
	}

	return timeObj.Unix() - 3600*8, nil
}

// GetLastFriday 获取上周五某个小时的时间
func GetLastFriday(now time.Time, targetHour int) time.Time {
	weekday := now.Weekday()

	// 计算上周五的时间
	daysToLastFriday := (int(weekday) - int(time.Friday) + 7) % 7
	if daysToLastFriday == 0 {
		daysToLastFriday = 7
	}
	lastFriday := now.AddDate(0, 0, -daysToLastFriday)

	lastFridayAtTargetHour := time.Date(
		lastFriday.Year(), lastFriday.Month(), lastFriday.Day(),
		targetHour, 0, 0, 0, // Nanoseconds set to zero
		lastFriday.Location(),
	)
	return lastFridayAtTargetHour
}

// GetThisFriday 获取这周五某个小时的时间
func GetThisFriday(now time.Time, targetHour int) time.Time {
	weekday := now.Weekday()
	// Adjust weekday so that Monday = 0, ..., Sunday = 6
	adjustedWeekday := (int(weekday) + 6) % 7
	// Calculate the number of days to this week's Friday (always return this Friday)
	daysToThisFriday := (4 - adjustedWeekday + 7) % 7
	thisFriday := now.AddDate(0, 0, daysToThisFriday)

	// Set the time to the target hour
	thisFridayAtTargetHour := time.Date(
		thisFriday.Year(), thisFriday.Month(), thisFriday.Day(),
		targetHour, 0, 0, 0,
		thisFriday.Location(),
	)
	return thisFridayAtTargetHour
}

func SleepRandSeconds(min int, max int) {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	if max < min {
		// 交换
		max, min = min, max
	}
	randomSeconds := rand.Intn(max-min+1) + min
	// 使用time.Sleep进行等待
	time.Sleep(time.Duration(randomSeconds) * time.Second)
}

func CurDayStartTimeStamp() int64 {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return newTime.UnixMilli()
}

func CurPrevious24HStamp() int64 {
	t := time.Now().Add(-24 * time.Hour) // 获取24小时前的时间
	return t.UnixMilli()
}

func TimeStampToDate(ts int64) string {
	// 将毫秒时间戳转换为秒和纳秒
	t := time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond))
	return t.Format("2006-01-02T15:04:05")
}

func CheckIfExceededHours(timestampMillis int64, hours int) bool {
	// 将毫秒时间戳转换为时间.Time
	timestamp := time.Unix(0, timestampMillis*int64(time.Millisecond))

	// 获取当前时间
	now := time.Now()

	// 计算时间差
	duration := now.Sub(timestamp)

	// 检查是否超过指定的小时数
	return duration.Hours() > float64(hours)
}

func MapToSlice(arr []int, fn func(s int) string) []string {
	var newArray []string
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}
