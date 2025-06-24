package mysql

import (
	"sort"

	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/app/metrics/repository/group_key/engine"
	"github.com/rentiansheng/dashboard/app/metrics/repository/group_key/impl/mysql/model"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"gorm.io/gorm"
)

type GroupKeyImpl struct {
	db *gorm.DB
}

func (g *GroupKeyImpl) Roots(ctx context.Context) ([]define.GroupKey, errors.Error) {
	rows := make([]model.GroupKeyTab, 0)
	if err := g.db.Where("parent_id = 0 and group_key_status = ?", model.StatusNormal).
		Find(&rows).Error; err != nil {
		ctx.Log().Errorf("get group key roots error: %v", err)
		return nil, ctx.Error().Error(code.DBExecuteErrCode, err)
	}

	res := make([]define.GroupKey, 0, len(rows))
	if err := ctx.AllMapper("mapper group key", &rows, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (g *GroupKeyImpl) Tree(ctx context.Context, rootId uint64) ([]*define.GroupKeyTree, errors.Error) {

	rows := make([]model.GroupKeyTab, 0)
	if err := g.db.Where("group_key_status = ?", model.StatusNormal).
		Find(&rows).Error; err != nil {
		ctx.Log().Errorf("get group key tree id: %v, error: %v", err)
		return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}

	return g.buildTree(ctx, rows)
}

func (g *GroupKeyImpl) buildTree(ctx context.Context, nodes []model.GroupKeyTab) ([]*define.GroupKeyTree, errors.Error) {
	res := make([]*define.GroupKeyTree, 0)
	rela := make(map[uint64]*define.GroupKeyTree)

	for _, row := range nodes {
		node := &define.GroupKeyTree{
			ID:          row.ID,
			DisplayName: row.DisplayName,
			Order:       row.OrderIdx,
			Children:    make([]*define.GroupKeyTree, 0),
		}
		if row.ParentID == 0 {
			res = append(res, node)
		}
		rela[row.ID] = node
	}
	for _, row := range nodes {
		if row.ParentID == 0 {
			continue
		}
		parentNode, ok := rela[row.ParentID]
		if !ok {

			continue
		}
		parentNode.Children = append(parentNode.Children, rela[row.ID])

	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Order < res[j].Order
	})

	for _, node := range res {
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Order < node.Children[j].Order
		})
	}

	return res, nil

}

func (g *GroupKeyImpl) ChildrenIds(ctx context.Context, id uint64) ([]uint64, errors.Error) {
	return g.childrenIds(ctx, id)
}
func (g *GroupKeyImpl) childrenIds(ctx context.Context, id uint64) ([]uint64, errors.Error) {

	if id == 0 {
		return nil, ctx.Error().Errorf(code.InstNotFoundKVCode, "group key", "id", id)
	}

	rows := make([]model.GroupKeyTab, 0)
	if err := g.db.Where("parent_id = ? and group_key_status = ?", id, model.StatusNormal).
		Find(&rows).Error; err != nil {
		ctx.Log().Errorf("get group key children ids error: %v", err)
		return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}

	res := make([]uint64, 0, len(rows))
	for _, row := range rows {
		res = append(res, row.ID)
	}
	return res, nil
}

func (g *GroupKeyImpl) AllChildrenIds(ctx context.Context, id uint64) ([]uint64, errors.Error) {

	if id == 0 {
		return nil, ctx.Error().Errorf(code.InstNotFoundKVCode, "group key", "id", id)
	}
	res := make([]uint64, 0)
	parentIds := []uint64{id}
	for len(parentIds) > 0 {
		rows := make([]model.GroupKeyTab, 0)
		if err := g.db.Where("parent_id IN (?) and group_key_status = ?", parentIds, model.StatusNormal).Find(&rows).Error; err != nil {
			ctx.Log().Errorf("get group key all children ids error: %v", err)
			return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
		}
		parentIds = []uint64{}
		for _, row := range rows {
			parentIds = append(parentIds, row.ID)
			res = append(res, row.ID)
		}
	}
	return res, nil
}

func (g *GroupKeyImpl) AllChildrenRelationNames(ctx context.Context, id uint64) ([]string, errors.Error) {
	ids, err := g.allRelationInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		// 如果没有子节点，直接返回空
		return nil, nil
	}
	rows := make([]model.DataGroupTab, len(ids))
	if err := g.db.Where("id IN (?) and data_group_status = ?", ids, model.StatusNormal).Find(&rows).Error; err != nil {
		ctx.Log().Errorf("get group key all children relation names error: %v", err)
		return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}
	names := make([]string, 0, len(ids))
	for _, row := range rows {
		names = append(names, row.DisplayName)
	}

	return names, nil
}

func (g *GroupKeyImpl) AllChildrenRelationIds(ctx context.Context, id uint64) ([]uint64, errors.Error) {
	return g.allRelationInfo(ctx, id)
}
func (g *GroupKeyImpl) allRelationInfo(ctx context.Context, id uint64) ([]uint64, errors.Error) {
	ids, err := g.AllChildrenIds(ctx, id)
	if err != nil {
		return nil, err
	}

	ids = append(ids, id)

	rows := make([]model.RelationTab, 0)
	if err := g.db.Where("group_key_id IN (?) and relation_status = ?", ids, model.StatusNormal).
		Find(&rows).Error; err != nil {
		ctx.Log().Errorf("get group key all children relation ids error: %v", err)
		return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}

	res := make([]uint64, 0, len(rows))
	idRela := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if _, ok := idRela[row.ID]; ok {
			continue
		}
		idRela[row.ID] = struct{}{}
		res = append(res, row.ID)
	}
	return ids, nil
}

func NewGroupKeyRepo(db *gorm.DB) engine.Engine {
	return &GroupKeyImpl{
		db: db,
	}
}
