package code

/***************************
    @author: tiansheng.ren
    @date: 2022/7/3
    @desc:

***************************/

// 14 开头表示es操作错误
const (
	// ESTranslateSQLErrCode es translate sql command error.
	ESTranslateSQLErrCode int32 = 14000
	// ESSearchErrCode es search reply error. reply(status: %v, type: %v, error: %v)
	ESSearchErrCode int32 = 14001
	// ESSearchResponseDecodeErrCode  es search response body decode error.  err: %s
	ESSearchResponseDecodeErrCode int32 = 14002
	ESGetMetricParamErrCode       int32 = 14003

	ESUpdateByIdErrCode  int32 = 14004
	EsBatchCreateErrCode int32 = 14005
	// EsCreateErrCode es create data error
	EsCreateErrCode int32 = 14006
	// ESSearchCountErrCode es count error: %v
	ESSearchCountErrCode    int32 = 14007
	ESGetSearchParamErrCode int32 = 14008

	// ESOrmSearchErrCode search error. %v
	ESOrmSearchErrCode int32 = 14009
	// ESOrmSearchNotFoundErrCode  not found. data source: %v, value: %v"
	ESOrmSearchNotFoundErrCode int32 = 14010
	// ESOrmChartParamTypeErrCode chart param type err. input: %v, expected: %v
	ESOrmChartParamTypeErrCode int32 = 14011
	// ESDeleteErrCode  es delete error. err: %v
	ESDeleteErrCode int32 = 14012

	// ESDALSearchErrCode es search error. err: %v
	ESDALSearchErrCode int32 = 14013
	// ESDALMUpsertErrCode es multi-upsert error. err: %v
	ESDALMUpsertErrCode int32 = 14014
	// ESDALUpdateByIdErrCode es update error. err: %v
	ESDALUpdateByIdErrCode int32 = 14015
	// ESDALSearchLogicTipErrCode es search error. action: %v, err: %v
	ESDALSearchLogicTipErrCode int32 = 14016
	// ESDALSaveErrCode es save error. inst type: %v, err: %v
	ESDALSaveErrCode int32 = 14017
	// ESDALDElByQueryErrCode es delete by query error. err: %v
	ESDALDElByQueryErrCode int32 = 14018

	// ESDALUpsertErrCode es upsert error. err: %v
	ESDALUpsertErrCode int32 = 14019
)
