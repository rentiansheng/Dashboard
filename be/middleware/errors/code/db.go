package code

const (
	// DBExecuteErrCode 执行数据库语句失败
	DBExecuteErrCode int32 = 12000
	// DBRecordDuplicateErrCode 数据重复，DB唯一校验失败
	DBRecordDuplicateErrCode int32 = 12001
	// CacheSetErrCode  cache set error. err: %v
	CacheSetErrCode int32 = 12002
	// CacheGetErrCode cache get error. err: %v
	CacheGetErrCode int32 = 12003
	// DBCommitErrCode db transaction commit error. err: %v
	DBCommitErrCode int32 = 12004
	// DBRollbackErrCode db transaction rollback error. err: %
	DBRollbackErrCode int32 = 12005

	// DBNotFoundErrCode  "db record not found. unique parameters: %s"
	DBNotFoundErrCode int32 = 12006

	// DBSaveErrCode insert data error. err: %s
	DBSaveErrCode int32 = 12009
	// DBNotFoundTipFieldErrCode  "db record not found. %v: %v"
	DBNotFoundTipFieldErrCode int32 = 12010

	// DBBeginTxnErrCode db transaction begin error. err: %v
	DBBeginTxnErrCode int32 = 12011

	// DBUpdateAffectedRowsNumErrCode update affected rows num wrong, err: %v
	DBUpdateAffectedRowsNumErrCode = 12012
)
