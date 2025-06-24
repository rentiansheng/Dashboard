package code

/***************************
    @author: tiansheng.ren
    @date: 2025/6/10
    @desc:

***************************/

const (
	// MetricDataSourceUnsupportedErrCode "data source(%s) not support"
	MetricDataSourceUnsupportedErrCode int32 = 1015001
	// MetricDetailCopyRowFieldErrCode %s copy field value from %s to %s error, to field must be not exists"
	MetricDetailCopyRowFieldErrCode = 1015002
	// MetricDataSourceQueryFilterRoleNotErrCode not found query filter role, operator: %s
	MetricDataSourceQueryFilterRoleNotErrCode = 1015003
	// MetricDataSourceIDNotFoundErrCode "data source id(%v) not found"
	MetricDataSourceIDNotFoundErrCode = 1015004
	// MetricDataSourceStorageUnsupportedErrCode "data source(%s) storage(%s) not support"
	MetricDataSourceStorageUnsupportedErrCode = 1015005

	// MetricDataSourceQueryFilterNeedValueErrCode "data source filter operator(%s) field(%s) need value"
	MetricDataSourceQueryFilterNeedValueErrCode = 1015006
)
