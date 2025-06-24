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
	code.MapperActionErrCode: "mapper error. action: %s, err: %s",
	// FileNotFoundErrCode file not found. file name: %s
	code.FileNotFoundErrCode: "file not found. file name: %s",

	/* ========= 12 开头 表示db操作错误 ========= */

	code.DBExecuteErrCode:         "执行数据库语句失败",
	code.DBRecordDuplicateErrCode: "数据重复，DB唯一校验失败",
	code.CacheSetErrCode:          "cache set error. err: %v",
	code.CacheGetErrCode:          "cache get error. err: %v",
	code.DBCommitErrCode:          "db transaction commit error. err: %v",
	code.DBRollbackErrCode:        "db transaction rollback error. err: %v",
	code.DBNotFoundErrCode:        "db record not found. unique parameters: %s",
	code.DBSaveErrCode:            "insert data error. err: %s",
	// DBNotFoundTipFieldErrCode  "db record not found. %v: %v"
	code.DBNotFoundTipFieldErrCode: "db record not found. %v: %v",
	// DBBeginTxnErrCode db transaction begin error. err: %v
	code.DBBeginTxnErrCode: "db transaction begin error. err: %v",
}
