package mysql

import (
	"github.com/rentiansheng/dashboard/app/metrics/define"
	dataSourceEngine "github.com/rentiansheng/dashboard/app/metrics/repository/data_source/engine"
	"github.com/rentiansheng/dashboard/app/metrics/repository/data_source/impl/mysql/model"
	"github.com/rentiansheng/dashboard/middleware/context"
	"github.com/rentiansheng/dashboard/middleware/errors"
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"gorm.io/gorm"
)

type DataSourceImpl struct {
	db *gorm.DB
}

func NewDataSourceRepo(db *gorm.DB) dataSourceEngine.Engine {
	return &DataSourceImpl{
		db: db,
	}
}

func (r *DataSourceImpl) GetDataSourceList(ctx context.Context) ([]define.MetricDataSource, int64, errors.Error) {
	rows := make([]model.DataSourceTab, 0)
	if err := r.db.Where(model.DataSourceStatusField+" = ?", model.DataSourceStatusNormal).
		Find(&rows).Error; err != nil {
		ctx.Log().Errorf("get data source list error: %v", err)
		return nil, 0, ctx.Error().Error(code.DBExecuteErrCode, err)
	}

	res := make([]define.MetricDataSource, 0, len(rows))
	if err := ctx.AllMapper("mapper data source", &rows, &res); err != nil {
		return nil, 0, err
	}
	return res, int64(len(res)), nil
}

func (r *DataSourceImpl) GetDataSourceByID(ctx context.Context, id uint64) (*define.MetricDataSource, errors.Error) {
	res := &define.MetricDataSource{}

	row := &model.DataSourceTab{}
	if err := r.db.Where(model.DataSourceStatusField+" = ?", model.DataSourceStatusNormal).
		Where("id = ?", id).
		Find(row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ctx.Error().Errorf(code.DBNotFoundTipFieldErrCode, "data source id", id).SetError(err)
		}
		return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}

	if err := ctx.Mapper("mapper data source", row, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *DataSourceImpl) GetDataSourceMeta(ctx context.Context, dataSourceID uint64) ([]define.MetricDataSourceMeta, errors.Error) {
	rows := make([]model.DataSourceMetaTab, 0)
	if err := r.db.Where("data_source_id = ? and data_source_meta_status = ?", dataSourceID, model.DataSourceStatusNormal).
		Find(&rows).Error; err != nil {
		return nil, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}
	res := make([]define.MetricDataSourceMeta, 0, len(rows))
	if err := ctx.AllMapper("mapper data source meta", rows, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *DataSourceImpl) GetDataSourceByName(ctx context.Context, name string) (define.MetricDataSource, errors.Error) {
	res := define.MetricDataSource{}

	var row model.DataSourceTab
	if err := r.db.Where("data_name = ?", name).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, ctx.Error().Errorf(code.DBNotFoundTipFieldErrCode, "data source name", name).SetError(err)
		}
		return res, ctx.Error().Errorf(code.DBExecuteErrCode, err)
	}

	if err := ctx.Mapper("mapper data source", row, &res); err != nil {
		return res, err
	}

	return res, nil
}
