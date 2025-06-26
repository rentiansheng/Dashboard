package adjust

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/app/metrics/repository"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"

	"github.com/rentiansheng/ges"
)

const ()

type ExtraFilterFn func(ctx context.Context, dataSource define.MetricDataSource, id uint64) (ges.Filter, errors.Error)

func ExtraGroupKeyFilterFn(ctx context.Context, dataSource define.MetricDataSource, id uint64) (ges.Filter, errors.Error) {
	const skipGroupKey = "-"
	if dataSource.GroupKey == skipGroupKey {
		// 如果没有group key，则不需要额外的过滤
		return nil, nil
	}
	if dataSource.EnableGroupKeyName {
		// 需要用names
		names, err := repository.GroupKey().AllChildrenRelationNames(ctx, id)
		if err != nil {
			return nil, err
		}
		return ges.Terms(dataSource.GroupKey, names), nil
	}
	ids, err := repository.GroupKey().AllChildrenRelationIds(ctx, id)
	if err != nil {
		return nil, err
	}
	return ges.Terms(dataSource.GroupKey, ids), nil
}
