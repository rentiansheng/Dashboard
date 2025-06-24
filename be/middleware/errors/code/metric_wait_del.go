package code

/***************************
    @author: tiansheng.ren
    @date: 2022/7/14
    @desc:

***************************/

/* ========= 15 开头 表示app metric领域内逻辑处理操作错误 ========= */

const (
	// UnimplementedChartHandlerErrCode handler %s unimplemented
	UnimplementedChartHandlerErrCode = 15001
	// SearchMetricListErrCode search metric value err
	SearchMetricListErrCode    = 15002
	HandleEsSearchParamErrCode = 15003
	TransactionQlangErrCode    = 15004
	// MetricListDetailValueCountErrCode metric list detail value count expect %d, actual %d
	MetricListDetailValueCountErrCode = 15005

	// MetricListDetailAPIKeyFnUnimplementedCode chart key(%s) fetch data api unimplemented, func name: %s
	MetricListDetailAPIKeyFnUnimplementedCode = 15006

	// MetricChartListCycleGroupNameTaskNotFoundErrCode find chart cycle group name %s cycle %d not found
	MetricChartListCycleGroupNameTaskNotFoundErrCode = 15008

	// MetricDataSourceQueryFilterBetweenParamErrCode "field(%s) between parameter error. value must be 2"
	MetricDataSourceQueryFilterBetweenParamErrCode = 15010

	// MetricChartBuildTipsErrCode "metric chart build tips err."
	MetricChartBuildTipsErrCode = 15015

	// MetricFieldFilterOperatorNotSupportErrCode "metric search filter field(%s)  operator(%s) not support"
	MetricFieldFilterOperatorNotSupportErrCode = 15016
	// MetricFieldFilterSuffixErrCode "metric search filter field(%s) suffix, need values must be string"
	MetricFieldFilterSuffixErrCode = 15017

	// MetricESOrmHandlerExtraParamNeedSetErrCode "metric esorm handler, extra parameter need set"
	MetricESOrmHandlerExtraParamNeedSetErrCode = 15018

	// MetricMergeChartErrCode "merge same charts err: %v"
	MetricMergeChartErrCode = 15019
	// MetricBuildChartErrCode "build chart err: %v"
	MetricBuildChartErrCode = 15020

	// MetricESOrmHandlerValueFnNameNotFoundErrCode "metric esorm handler, value function(%s) not found
	MetricESOrmHandlerValueFnNameNotFoundErrCode = 15021

	// MetricHandlerCfgIllegalErrCode "metric handler config illegal, err: %s"
	MetricHandlerCfgIllegalErrCode = 15022

	// MetricESOrmHandlerValueNeedSetErrCode "metric esorm handler, agg filter terms value  need set"
	MetricESOrmHandlerValueNeedSetErrCode = 15023

	// MetricESOrmSubAggOperateTypeUnImplementCode "metric esorm sub agg operate type un implement, type: %s"
	MetricESOrmSubAggOperateTypeUnImplementCode = 15024

	// MetricChartDetailNotSupportSearchErrCode  "metric chart detail not support search, chart name: %s"
	MetricChartDetailNotSupportSearchErrCode = 15025
	// MetricChartSubAggTypeUnImplementCode "metric chart sub agg type un implement,  chart name: %s, type: %v"
	MetricChartSubAggTypeUnImplementCode = 15026

	// DataSourceQueryConditionExistErrCode "data source query condition exist, name: %s, rule: %s"
	DataSourceQueryConditionExistErrCode = 15027

	// MetricESORMCalculateFormalAggNotFound "metric name %s aggregate name %s not found. plz calculate formula and legend"
	MetricESORMCalculateFormalAggNotFound = 15028

	// MetricFieldFilterNumIllegalErrCode  "metric search filter field(%s) num illegal, must be %v"
	MetricFieldFilterNumIllegalErrCode = 15029
)
