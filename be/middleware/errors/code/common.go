package code

const (
	// JSONDecodeErrCode request body decode error. err: %s
	JSONDecodeErrCode int32 = 10000
	// ConvertToErrCode field %s convert type(%s) error(%s).
	ConvertToErrCode int32 = 10006
	// InstNotFoundKVCode %v not found. %v: %v
	InstNotFoundKVCode int32 = 10008
	// TemplateFormatErr  template format error. %v
	TemplateFormatErr int32 = 10011
	// TemplateBuildErr  template build error. %v
	TemplateBuildErr int32 = 10012

	// MapperActionErrCode  mapper error. action: %s, err: %s
	MapperActionErrCode int32 = 10019
	// FileNotFoundErrCode file not found. file name: %s
	FileNotFoundErrCode int32 = 10027
)
