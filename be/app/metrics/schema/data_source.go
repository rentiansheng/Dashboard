package schema

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
)

func (m *DataSourceSchema) List(ctx context.Contexts) errors.Error {
	dataSources, total, err := m.dataSourceSvc.DataSource(ctx)
	if err != nil {
		return err
	}

	resp := define.DataSourceDataResp{}
	if err := ctx.Mapper("data source response", dataSources, &resp.List); err != nil {
		return err
	}
	resp.Total = (total)

	ctx.SetData(resp)
	return nil
}

func (m *DataSourceSchema) DataSourceMeta(ctx context.Contexts) errors.Error {

	req := &define.DataSourceMetaReq{}
	if err := ctx.Decode(req); err != nil {
		return err
	}

	dataSources, err := m.dataSourceSvc.DataSourceMeta(ctx, req.ID)
	if err != nil {
		return err
	}

	resp := define.DataSourceFieldMetaResp{}
	if err := ctx.Mapper("data source response", dataSources, &resp.Fields); err != nil {
		return err
	}
	ctx.SetData(resp)
	return nil
}

func (m *DataSourceSchema) DataSourceQuery(ctx context.Contexts) errors.Error {
	req := &define.QueryReq{}
	if err := ctx.Decode(req); err != nil {
		return err
	}

	list, total, err := m.dataSourceSvc.QueryTable(ctx, req)
	if err != nil {
		return err
	}
	ctx.SetData(map[string]interface{}{"list": list, "total": total})
	return nil

}

func (m *DataSourceSchema) DataSourceChart(ctx context.Contexts) errors.Error {
	req := &define.QueryReq{}
	if err := ctx.Decode(req); err != nil {
		return err
	}

	result, err := m.dataSourceSvc.QueryChart(ctx, req)
	if err != nil {
		return err
	}
	ctx.SetData(result)
	return nil

}

func (m *DataSourceSchema) DataSourceEnum(ctx context.Contexts) errors.Error {
	req := &define.DataSourceFieldEnumReq{}
	if err := ctx.Decode(req); err != nil {
		return err
	}

	result, err := m.dataSourceSvc.DataSourceFieldEnum(ctx, req)
	if err != nil {
		return err
	}
	ctx.SetData(result)
	return nil

}

