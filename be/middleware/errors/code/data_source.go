package code

/***************************
    @author: tiansheng.ren
    @date: 2025/6/10
    @desc:

***************************/

const (
	// MetricDataSourceUnsupportedErrCode "data source(%s) not support"
	MetricDataSourceUnsupportedErrCode int32 = 15001
	// MetricDetailCopyRowFieldErrCode %s copy field value from %s to %s error, to field must be not exists"
	MetricDetailCopyRowFieldErrCode = 15007
	// MetricDataSourceQueryFilterRoleNotErrCode not found query filter role, operator: %s
	MetricDataSourceQueryFilterRoleNotErrCode = 15009
	// MetricDataSourceIDNotFoundErrCode "data source id(%v) not found"
	MetricDataSourceIDNotFoundErrCode = 15012
	// MetricDataSourceStorageUnsupportedErrCode "data source(%s) storage(%s) not support"
	MetricDataSourceStorageUnsupportedErrCode = 15013

	// MetricDataSourceQueryFilterNeedValueErrCode "data source filter operator(%s) field(%s) need value"
	MetricDataSourceQueryFilterNeedValueErrCode = 15014
)
