package adjust

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/ges"
)

// FieldAdjustFn  需要特殊处理的字段处理方法
type FieldAdjustFn func(ctx context.Context, rule define.QueryFilterRule) (define.QueryFilterRule, errors.Error)
type VirtualFieldExtraFilterFn func(ctx context.Context, rule define.QueryFilterRule) (ges.Filter, errors.Error)
