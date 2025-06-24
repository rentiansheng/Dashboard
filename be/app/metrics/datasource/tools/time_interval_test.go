package tools

import (
	"testing"
	"time"

	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTimeIntervalList(t *testing.T) {
	// tsList为秒级时间戳
	tsList := []uint64{
		1639711326, // 2021-12-17 11:22:06
		1639754526, // 2021-12-17 23:22:06
		1639797726, // 2021-12-18 11:22:06
		1638933726, // 2021-12-08 11:22:06
		1631244126, // 2021-09-10 11:22:06
		1612927326, // 2021-02-10 11:22:06
		1581304926, // 2020-02-10 11:22:06
	}

	for i := range tsList {
		tsList[i] = tsList[i] * thousand
	}

	res, err := GetTimeIntervalList(tsList, define.MetricCycleWeek, Millisecond, true, true)
	_ = res
	assert.Nil(t, err)
}

func TestGetRangeTimeIntervalList(t *testing.T) {

	type in struct {
		tsRange   TimeInterval
		CycleType MetricCycleType
	}
	suits := []struct {
		in     in
		output []TimeInterval
	}{
		{
			in: in{
				tsRange: TimeInterval{
					// 2022-06-07 00:00:00
					Start: 1654531200,
					// 2022-06-10 23:59:59
					End: 1654876799,
				},
				CycleType: define.MetricCycleDay,
			},
			output: []TimeInterval{
				{Start: 1654531200, End: 1654617599},
				{Start: 1654617600, End: 1654703999},
				{Start: 1654704000, End: 1654790399},
				{Start: 1654790400, End: 1654876799},
			},
		},
		{
			in: in{
				tsRange: TimeInterval{
					//2022-01-20 00:00:00
					Start: 1642608000,
					// 2022-06-20 00:00:00
					End: 1655654400,
				},
				CycleType: define.MetricCycleQuarter,
			},
			output: []TimeInterval{
				// 1-3
				{
					// 2022-01-01 00:00:00
					Start: 1640966400,
					// 2022-03-31 23:59:59
					End: 1648742399,
				},
				// 4-6
				{
					// 2022-04-01 00:00:00
					Start: 1648742400,
					// 2022-06-30 23:59:59
					End: 1656604799,
				},
			},
		},
	}

	for idx, item := range suits {
		output, err := GetRangeTimeIntervalList(item.in.tsRange, item.in.CycleType)
		require.Equal(t, nil, err, "index; %d", idx)
		require.Equal(t, item.output, output, "index; %d", idx)

	}
}

func TestGetWeekdays(t *testing.T) {

	type args struct {
		startTime string
		endTime   string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"contain weekend", args{
			startTime: "2022-09-09 10:00:00",
			endTime:   "2022-09-10 10:00:00",
		}, 50400},
		{"contain weekend2", args{
			startTime: "2022-09-09 10:00:00",
			endTime:   "2022-09-12 10:00:00",
		}, 86400},
		{"not contain weekend", args{
			startTime: "2022-09-05 09:00:00",
			endTime:   "2022-09-09 10:00:00",
		}, 349200},
		{"is weekend", args{
			startTime: "2022-09-10 09:00:00",
			endTime:   "2022-09-11 10:00:00",
		}, 0},
		{"is today", args{
			startTime: "2022-09-09 10:00:00",
			endTime:   "2022-09-09 11:00:00",
		}, 3600},
		{"greater", args{
			startTime: "2022-09-10 10:00:00",
			endTime:   "2022-09-09 11:00:00",
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWeekdaysTime(tt.args.startTime, tt.args.endTime)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("GetWeekdaysTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHalfYear(t *testing.T) {
	type args struct {
		Year  int
		Month int
		Day   int
		Hour  int
	}
	suits := []struct {
		input args
		// output[1] 为了方便测试，是下一个周期的开始时间，结束时间每个月是不一样的
		output [2]args
	}{
		{
			input: args{
				Year:  2022,
				Month: 1,
				Day:   1,
				Hour:  10,
			},
			output: [2]args{
				{Year: 2022, Month: 1, Day: 1},
				{Year: 2022, Month: 7, Day: 1},
			},
		},
		{
			input: args{
				Year:  2022,
				Month: 2,
				Day:   15,
				Hour:  10,
			},
			output: [2]args{
				{Year: 2022, Month: 1, Day: 1},
				{Year: 2022, Month: 7, Day: 1},
			},
		},
		{
			input: args{
				Year:  2022,
				Month: 6,
				Day:   10,
				Hour:  10,
			},
			output: [2]args{
				{Year: 2022, Month: 1, Day: 1},
				{Year: 2022, Month: 7, Day: 1},
			},
		},
		{
			input: args{
				Year:  2022,
				Month: 7,
				Day:   1,
				Hour:  10,
			},
			output: [2]args{
				{Year: 2022, Month: 7, Day: 1},
				{Year: 2023, Month: 1, Day: 1},
			},
		},
		{
			input: args{
				Year:  2022,
				Month: 8,
				Day:   10,
				Hour:  10,
			},
			output: [2]args{
				{Year: 2022, Month: 7, Day: 1},
				{Year: 2023, Month: 1, Day: 1},
			},
		},
	}
	for idx, suit := range suits {
		ts := time.Date(suit.input.Year, time.Month(suit.input.Month), suit.input.Day, suit.input.Hour, 0, 0, 0, time.UTC)
		dateRange, err := getHalfYearInterval(uint64(ts.Unix()))
		require.Equal(t, nil, err)
		startTime := time.Unix(int64(dateRange.Start), 0)
		endTime := time.Unix(int64(dateRange.End+1), 0)
		require.Equal(t, suit.output[0].Year, startTime.Year(), "index", idx)
		require.Equal(t, suit.output[0].Month, int(startTime.Month()), "index", idx)
		require.Equal(t, suit.output[0].Day, startTime.Day(), "index", idx)
		require.Equal(t, suit.output[1].Year, endTime.Year(), "index", idx)
		require.Equal(t, suit.output[1].Month, int(endTime.Month()), "index", idx)
		require.Equal(t, suit.output[1].Day, endTime.Day(), "index", idx)
	}
}
