package code

const (
	// JSONDecodeErrCode request body decode error. err: %s
	JSONDecodeErrCode int32 = 1010000
	// ConvertToErrCode field %s convert type(%s) error(%s).
	ConvertToErrCode int32 = 1010002
	// InstNotFoundKVCode %v not found. %v: %v
	InstNotFoundKVCode int32 = 1010003
	// TemplateFormatErr  template format error. %v
	TemplateFormatErr int32 = 1010004
	// TemplateBuildErr  template build error. %v
	TemplateBuildErr int32 = 1010005

	// MapperActionErrCode  mapper error. action: %s, err: %s
	MapperActionErrCode int32 = 1010006
	// FileNotFoundErrCode file not found. file name: %s
	FileNotFoundErrCode int32 = 1010007
)
