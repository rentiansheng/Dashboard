package _default

import (
	"github.com/rentiansheng/dashboard/middleware/errors/code"
	"github.com/rentiansheng/dashboard/middleware/errors/register"
)

/*
**************************

	@author: tiansheng.ren
	@date: 2022/7/3
	@desc:

**************************
*/
func init() {
	register.Register(langName, asstCodes)
}

var asstCodes = map[int32]string{
	/* ========= 10 开头 表示公共 ========= */

	code.JSONDecodeErrCode:   "request body decode error. err: %s",
	code.ConvertToErrCode:    "field %s convert type(%s) error(%s)",
	code.InstNotFoundKVCode:  "%v not found. %v: %v",
	code.TemplateFormatErr:   "template format error. %v",
	code.TemplateBuildErr:    "template build error. %v",
	code.MapperActionErrCode: "mapper error. action: %s, err: %s",
	code.FileNotFoundErrCode: "file not found. file name: %s",

	/* ========= 12 开头 表示db操作错误 ========= */

	code.DBExecuteErrCode:               "执行数据库语句失败",
	code.DBRecordDuplicateErrCode:       "数据重复，DB唯一校验失败",
	code.CacheSetErrCode:                "cache set error. err: %v",
	code.CacheGetErrCode:                "cache get error. err: %v",
	code.DBCommitErrCode:                "db transaction commit error. err: %v",
	code.DBRollbackErrCode:              "db transaction rollback error. err: %v",
	code.DBNotFoundErrCode:              "db record not found. unique parameters: %s",
	code.DBSaveErrCode:                  "insert data error. err: %s",
	code.DBNotFoundTipFieldErrCode:      "db record not found. %v: %v",
	code.DBBeginTxnErrCode:              "db transaction begin error. err: %v",
	code.DBUpdateAffectedRowsNumErrCode: "update affected rows num wrong, err: %v",

	/* ========= 14 开头 表示es操作错误 ========= */

	code.ESTranslateSQLErrCode:         "es translate sql command error",
	code.ESSearchErrCode:               "es search reply error. reply(status: %v, type: %v, error: %v)",
	code.ESSearchResponseDecodeErrCode: "es search response body decode error. err: %s",
	code.ESGetMetricParamErrCode:       "es get metric param error",
	code.ESUpdateByIdErrCode:           "es update by id error",
	code.EsBatchCreateErrCode:          "es batch create error",
	code.EsCreateErrCode:               "es create data error",
	code.ESSearchCountErrCode:          "es count error: %v",
	code.ESGetSearchParamErrCode:       "es get search param error",
	code.ESOrmSearchErrCode:            "search error. %v",
	code.ESOrmSearchNotFoundErrCode:    "not found. data source: %v, value: %v",
	code.ESOrmChartParamTypeErrCode:    "chart param type err. input: %v, expected: %v",
	code.ESDeleteErrCode:               "es delete error. err: %v",
	code.ESDALSearchErrCode:            "es search error. err: %v",
	code.ESDALMUpsertErrCode:           "es multi-upsert error. err: %v",
	code.ESDALUpdateByIdErrCode:        "es update error. err: %v",
	code.ESDALSearchLogicTipErrCode:    "es search error. action: %v, err: %v",
	code.ESDALSaveErrCode:              "es save error. inst type: %v, err: %v",
	code.ESDALDElByQueryErrCode:        "es delete by query error. err: %v",
	code.ESDALUpsertErrCode:            "es upsert error. err: %v",

	/* ========= 15 开头 表示数据源相关错误 ========= */

	code.MetricDataSourceUnsupportedErrCode:          "data source(%s) not support",
	code.MetricDetailCopyRowFieldErrCode:             "%s copy field value from %s to %s error, to field must be not exists",
	code.MetricDataSourceQueryFilterRoleNotErrCode:   "not found query filter role, operator: %s",
	code.MetricDataSourceIDNotFoundErrCode:           "data source id(%v) not found",
	code.MetricDataSourceStorageUnsupportedErrCode:   "data source(%s) storage(%s) not support",
	code.MetricDataSourceQueryFilterNeedValueErrCode: "data source filter operator(%s) field(%s) need value",

	/* ========= 99 开头 表示原始错误包装 ========= */

	code.RawErrWrap: "raw error wrap: %v",
}
