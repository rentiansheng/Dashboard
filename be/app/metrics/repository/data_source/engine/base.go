package engine

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
)

type Engine interface {
	GetDataSourceList(ctx context.Context) ([]define.MetricDataSource, int64, errors.Error)
	GetDataSourceMeta(ctx context.Context, dataSourceID uint64) ([]define.MetricDataSourceMeta, errors.Error)
	GetDataSourceByID(ctx context.Context, id uint64) (*define.MetricDataSource, errors.Error)
	GetDataSourceByName(ctx context.Context, name string) (define.MetricDataSource, errors.Error)
}
