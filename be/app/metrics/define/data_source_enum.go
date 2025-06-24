package define

type DataRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type DataSourceMetaEnumReq struct {
	DepartmentId int64         `json:"department_id"`
	DataSourceId uint64        `json:"data_source_id"`
	Date         DataRange     `json:"date"`
	FieldName    string        `json:"field_name"`
	FieldValue   []interface{} `json:"field_value"`
	Relations    []struct {
		FieldName  string `json:"field_name"`
		FieldValue string `json:"field_value"`
	} `json:"relations"`
}

const (
	DataSourceMetaFieldEnumDataTypeKv    DataSourceMetaFieldEnumDataType = "kv"
	DataSourceMetaFieldEnumDataTypeArray DataSourceMetaFieldEnumDataType = "array"
)

type DataSourceMetaFieldEnumDataType string
type DataSourceMetaFieldEnumData struct {
	Typ    DataSourceMetaFieldEnumDataType `json:"type"`
	Values interface{}                     `json:"values"`
}
