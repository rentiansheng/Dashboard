package code

const (
	// DBExecuteErrCode 执行数据库语句失败
	DBExecuteErrCode int32 = 1012000
	// DBRecordDuplicateErrCode 数据重复，DB唯一校验失败
	DBRecordDuplicateErrCode int32 = 1012001
	// CacheSetErrCode  cache set error. err: %v
	CacheSetErrCode int32 = 1012002
	// CacheGetErrCode cache get error. err: %v
	CacheGetErrCode int32 = 1012003
	// DBCommitErrCode db transaction commit error. err: %v
	DBCommitErrCode int32 = 1012004
	// DBRollbackErrCode db transaction rollback error. err: %
	DBRollbackErrCode int32 = 1012005

	// DBNotFoundErrCode  "db record not found. unique parameters: %s"
	DBNotFoundErrCode int32 = 1012006

	// DBSaveErrCode insert data error. err: %s
	DBSaveErrCode int32 = 1012007
	// DBNotFoundTipFieldErrCode  "db record not found. %v: %v"
	DBNotFoundTipFieldErrCode int32 = 1012008

	// DBBeginTxnErrCode db transaction begin error. err: %v
	DBBeginTxnErrCode int32 = 1012009

	// DBUpdateAffectedRowsNumErrCode update affected rows num wrong, err: %v
	DBUpdateAffectedRowsNumErrCode = 1012010
)
