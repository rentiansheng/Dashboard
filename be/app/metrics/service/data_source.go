package service

import (
	"github.com/rentiansheng/dashboard/app/metrics/datasource/format"
	"github.com/rentiansheng/dashboard/app/metrics/datasource/query"
	"github.com/rentiansheng/dashboard/app/metrics/define"
	"github.com/rentiansheng/dashboard/app/metrics/repository"
	"github.com/rentiansheng/dashboard/middleware/context"
	. "github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
)

type DataSourceService interface {
	DataSource(ctx context.Context) ([]define.MetricDataSource, int64, Error)
	DataSourceMeta(ctx context.Context, dataSourceID uint64) ([]define.MetricDataSourceMeta, Error)
	QueryTable(ctx context.Context, input *define.QueryReq) ([]map[string]interface{}, uint64, Error)
	QueryChart(ctx context.Context, input *define.QueryReq) (*define.Chart, Error)
	DataSourceFieldEnum(ctx context.Context, input *define.DataSourceFieldEnumReq) (*define.DataSourceFieldEnum, Error)
}

func NewDataSourceService() DataSourceService {
	return &DataSourceServiceImpl{}
}

type DataSourceServiceImpl struct {
}

func (d *DataSourceServiceImpl) DataSource(ctx context.Context) ([]define.MetricDataSource, int64, Error) {
	rows, total, err := repository.DataSource().GetDataSourceList(ctx)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (d *DataSourceServiceImpl) DataSourceMeta(ctx context.Context, dataSourceID uint64) ([]define.MetricDataSourceMeta, Error) {
	rows, err := repository.DataSource().GetDataSourceMeta(ctx, dataSourceID)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (d *DataSourceServiceImpl) QueryChart(ctx context.Context, input *define.QueryReq) (*define.Chart, Error) {

	dataSourceInfo, err := d.dataSourceInfo(ctx, input.DataSourceID)
	if err != nil {
		return nil, err
	}
	metas, err := d.adjustDataSourceQuery(ctx, dataSourceInfo.DataName, input, false)
	if err != nil {
		return nil, err
	}

	var chartObj *define.Chart
	switch dataSourceInfo.DataSourceType {
	case define.DataSourceTypeES:
		chartObj, err = query.SearchESMetric(ctx, dataSourceInfo, metas, input)
		if err != nil {
			return nil, err
		}

	default:
		return nil, ctx.Error().Errorf(code.MetricDataSourceStorageUnsupportedErrCode,
			dataSourceInfo.DataName, dataSourceInfo.DataSourceType)
	}

	format.SeriesValueFormat(chartObj)
	return chartObj, nil
}

func (d *DataSourceServiceImpl) QueryTable(ctx context.Context, input *define.QueryReq) ([]map[string]interface{}, uint64, Error) {
	dataSourceInfo, err := d.dataSourceInfo(ctx, input.DataSourceID)
	if err != nil {
		return nil, 0, err
	}

	fields, err := d.adjustDataSourceQuery(ctx, dataSourceInfo.DataName, input, true)
	if err != nil {
		return nil, 0, err
	}
	var results []map[string]interface{}
	var total uint64
	switch dataSourceInfo.DataSourceType {
	case define.DataSourceTypeES:
		results, total, err = query.SearchES(ctx, dataSourceInfo, input)
		if err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, ctx.Error().Errorf(code.MetricDataSourceStorageUnsupportedErrCode,
			dataSourceInfo.DataName, dataSourceInfo.DataSourceType)
	}

	detailMeta, err := query.DataSourceEnum(ctx, fields, input, results)
	if err != nil {
		return nil, 0, err
	}
	for idx, row := range results {
		results[idx], err = format.AdjustChartDetailPageInfoRowByFormatterTemplate(ctx, row, detailMeta, "user analysis detail")
		if err != nil {
			return nil, 0, err
		}
	}
	return results, total, nil

}

func (d *DataSourceServiceImpl) DataSourceFieldEnum(ctx context.Context, input *define.DataSourceFieldEnumReq) (*define.DataSourceFieldEnum, Error) {
	dataSourceInfo, err := d.dataSourceInfo(ctx, input.DataSourceId)
	if err != nil {
		return nil, err
	}

	switch dataSourceInfo.DataSourceType {
	case define.DataSourceTypeES:
		return query.EnumES(ctx, dataSourceInfo.DataName, input)
	}
	return nil, nil
}

func (d *DataSourceServiceImpl) dataSourceInfo(ctx context.Context, id uint64) (define.MetricDataSource, Error) {
	result := define.MetricDataSource{}
	dataSourceInfo, err := repository.DataSource().GetDataSourceByID(ctx, id)
	if err != nil {
		return result, err
	}
	if dataSourceInfo == nil {
		return define.MetricDataSource{}, ctx.Error().Errorf(code.MetricDataSourceIDNotFoundErrCode, id)
	}

	if err := ctx.Mapper("query metric. mapper data source", dataSourceInfo, &result); err != nil {
		return result, err
	}
	return result, nil
}

// adjustDataSourceQuery 根据索引新加对应的查询条件
func (d *DataSourceServiceImpl) adjustDataSourceQuery(ctx context.Context, dataSourceName string, input *define.QueryReq, isDetail bool) ([]define.MetricDataSourceMeta, Error) {
	metas, err := repository.DataSource().GetDataSourceMeta(ctx, input.DataSourceID)
	if err != nil {
		return nil, err
	}
	fields := make([]define.MetricDataSourceMeta, 0, len(metas))
	if err := ctx.Mapper("query metric. mapper data source", metas, &fields); err != nil {
		return nil, err
	}
	if err := format.AdjustQueryInput(ctx, dataSourceName, fields, input, isDetail); err != nil {
		return nil, err
	}
	return fields, nil
}
