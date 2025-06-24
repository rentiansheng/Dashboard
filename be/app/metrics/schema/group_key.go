package schema

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	. "github.com/rentiansheng/dashboard/middleware/errors"
)

func (g *GroupKeySchema) Roots(ctx context.Contexts) Error {

	res, err := g.groupKeySvc.Roots(ctx)
	if err != nil {
		return err
	}
	ctx.SetData(res)
	return nil
}

func (g *GroupKeySchema) Tree(ctx context.Contexts) Error {
	req := &define.GroupKeyTreeReq{}
	if err := ctx.Decode(req); err != nil {
		return err
	}

	res, err := g.groupKeySvc.Tree(ctx, req.RootID)
	if err != nil {
		return err
	}
	ctx.SetData(res)
	return nil
}
