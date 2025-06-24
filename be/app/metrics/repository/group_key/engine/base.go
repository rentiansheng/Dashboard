package engine

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
)

type Engine interface {
	Roots(ctx context.Context) ([]define.GroupKey, errors.Error)
	Tree(ctx context.Context, id uint64) ([]*define.GroupKeyTree, errors.Error)

	ChildrenIds(ctx context.Context, id uint64) ([]uint64, errors.Error)
	AllChildrenIds(ctx context.Context, id uint64) ([]uint64, errors.Error)
	AllChildrenRelationIds(ctx context.Context, id uint64) ([]uint64, errors.Error)
	AllChildrenRelationNames(ctx context.Context, id uint64) ([]string, errors.Error)

	// TODO: add management methods

	// AddRoot(ctx context.Context, node define.GroupKey) (define.GroupKey, errors.Error)
	// AddNode(ctx context.Context, node define.GroupKey) (define.GroupKey, errors.Error)
	// ChangeNode(ctx context.Context, node define.GroupKey) (define.GroupKey, errors.Error)
	// ChangeStatus(ctx context.Context, id uint64, status uint8) (define.GroupKey, errors.Error)
	// AddDataGroup(ctx context.Context, id uint64, status uint8) (define.DataGroup errors.Error)
	// ChangeDataGroup(ctx context.Context, id uint64, status uint8) (define.DataGroup, errors.Error)
	// ChangeDataGroupStatus(ctx context.Context, id uint64, status uint8) (define.DataGroup, errors.Error)
	// AddRelation(ctx context.Context, groupId , dataGroupId uint64) (errors.Error)
	// DelRelation(ctx context.Context, groupId , dataGroupId uint64) (errors.Error)
}
