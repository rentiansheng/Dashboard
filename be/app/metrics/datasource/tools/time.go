package tools

import (
	"time"

	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
)

func BuildXAxisKey(ctx context.Context, start, end int64) string {
	tm := time.Unix(start, 0)
	strStart := tm.Format("2006-01-02")
	tm = time.Unix(end, 0)
	strEnd := tm.Format("2006-01-02")
	return strStart + "\n" + strEnd
}

func BuildTimeRangeCycleList(ctx context.Context, cycle int8, timeRange define.TimeRange) ([]string, []uint64, errors.Error) {
	cycleRangeBorder := TimeInterval{Start: timeRange.Start, End: timeRange.End}
	cycleRangeList, err := GetRangeTimeIntervalList(cycleRangeBorder, MetricCycleType(cycle))
	if err != nil {
		ctx.Log().ErrorJSON("convert time range list error. cycle: %s, timeRange: %s, err: %s", cycle, timeRange, err.Error())
		return nil, nil, ctx.Error().Errorf(code.ConvertToErrCode, err, "timeRange", "cycleRangeList", err.Error())
	}
	keyList := make([]string, 0)
	existMap := make(map[uint64]struct{}, 0)
	times := make([]uint64, 0)
	ctx.Log().InfoJSON("cycle range list: %s", cycleRangeList)
	for _, cycleRange := range cycleRangeList {
		if _, ok := existMap[cycleRange.Start]; ok {
			continue
		}
		times = append(times, cycleRange.Start)
		keyList = append(keyList, BuildXAxisKey(ctx, int64(cycleRange.Start), int64(cycleRange.End)))
	}
	return keyList, times, nil
}

func BuildTimeRangeCycleBorder(ctx context.Context, cycle MetricCycleType, timeRange define.TimeRange) (define.TimeRange, errors.Error) {
	tmpDateRange, err := GetTimeIntervalRangeBorder(TimeInterval{
		Start: timeRange.Start,
		End:   timeRange.End,
	}, cycle)
	if err != nil {
		ctx.Log().ErrorJSON("start cycle error. data: %s, err: %s", timeRange, err)
		return define.TimeRange{}, ctx.Error().Errorf(code.ConvertToErrCode, "time cycle boder", err)
	}

	return define.TimeRange{
		Start: tmpDateRange.Start,
		End:   timeRange.End,
	}, nil
}

func SecondToMillisecond(ts uint64) uint64 {
	return Sec2MillisTimestampUint64(ts)
}

func Sec2MillisTimestampUint64(timestamp uint64) uint64 {
	return timestamp * thousand
}
