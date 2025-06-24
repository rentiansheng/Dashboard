package code

/***************************
    @author: tiansheng.ren
    @date: 2022/7/3
    @desc:

***************************/

// 1014 开头表示es操作错误
const (
	// ESTranslateSQLErrCode es translate sql command error.
	ESTranslateSQLErrCode int32 = 1014000
	// ESSearchErrCode es search reply error. reply(status: %v, type: %v, error: %v)
	ESSearchErrCode int32 = 1014001
	// ESSearchResponseDecodeErrCode  es search response body decode error.  err: %s
	ESSearchResponseDecodeErrCode int32 = 1014002
	ESGetMetricParamErrCode       int32 = 1014003

	ESUpdateByIdErrCode  int32 = 1014004
	EsBatchCreateErrCode int32 = 1014005
	// EsCreateErrCode es create data error
	EsCreateErrCode int32 = 1014006
	// ESSearchCountErrCode es count error: %v
	ESSearchCountErrCode    int32 = 1014007
	ESGetSearchParamErrCode int32 = 1014008

	// ESOrmSearchErrCode search error. %v
	ESOrmSearchErrCode int32 = 1014009
	// ESOrmSearchNotFoundErrCode  not found. data source: %v, value: %v"
	ESOrmSearchNotFoundErrCode int32 = 1014010
	// ESOrmChartParamTypeErrCode chart param type err. input: %v, expected: %v
	ESOrmChartParamTypeErrCode int32 = 1014011
	// ESDeleteErrCode  es delete error. err: %v
	ESDeleteErrCode int32 = 1014012

	// ESDALSearchErrCode es search error. err: %v
	ESDALSearchErrCode int32 = 1014013
	// ESDALMUpsertErrCode es multi-upsert error. err: %v
	ESDALMUpsertErrCode int32 = 1014014
	// ESDALUpdateByIdErrCode es update error. err: %v
	ESDALUpdateByIdErrCode int32 = 1014015
	// ESDALSearchLogicTipErrCode es search error. action: %v, err: %v
	ESDALSearchLogicTipErrCode int32 = 1014016
	// ESDALSaveErrCode es save error. inst type: %v, err: %v
	ESDALSaveErrCode int32 = 1014017
	// ESDALDElByQueryErrCode es delete by query error. err: %v
	ESDALDElByQueryErrCode int32 = 1014018

	// ESDALUpsertErrCode es upsert error. err: %v
	ESDALUpsertErrCode int32 = 1014019
)
