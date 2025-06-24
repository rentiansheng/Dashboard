package service

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/app/metrics/repository"
	"github.com/rentiansheng/dashboard/middleware/context"
	. "github.com/rentiansheng/dashboard/middleware/errors"
)

type GroupKeyService interface {
	Roots(ctx context.Context) ([]define.GroupKey, Error)
	Tree(ctx context.Context, id uint64) ([]*define.GroupKeyTree, Error)
}

func NewGroupKeyService() GroupKeyService {
	return &GroupKeyServiceImpl{}
}

type GroupKeyServiceImpl struct {
}

func (g *GroupKeyServiceImpl) Roots(ctx context.Context) ([]define.GroupKey, Error) {
	rows, err := repository.GroupKey().Roots(ctx)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (g *GroupKeyServiceImpl) Tree(ctx context.Context, id uint64) ([]*define.GroupKeyTree, Error) {
	rows, err := repository.GroupKey().Tree(ctx, id)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
