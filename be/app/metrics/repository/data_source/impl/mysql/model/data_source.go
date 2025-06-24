package model

type DataSourceTab struct {
	ID          uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DisplayName string `gorm:"column:display_name;NOT NULL" json:"name"`
	Tips        string `gorm:"column:tips"  json:"tips"`
	Remark      string `gorm:"column:remark" json:"remark"`
	DataName    string `gorm:"column:data_name"  json:"data_name"`
	Status      uint   `gorm:"column:data_source_status;NOT NULL" json:"status"`
	// group key unused field, used to group data source.
	GroupKey           string `gorm:"column:group_key;NOT NULL" json:"group_key"`
	EnableGroupKeyName bool   `gorm:"column:enable_group_key_name;NOT NULL" json:"enable_group_key_name"`

	//`data_type` VARCHAR(45) NOT NULL COMMENT 'data source storage type. eg: mysql,es'  version 1 only support es
	DataSourceType string `gorm:"column:data_type;NOT NULL" json:"data_type"`
	SortFields     string `gorm:"column:sort_fields;NOT NULL"`
	Mtime          uint32 `gorm:"column:mtime;NOT NULL" json:"mtime"`
	Ctime          uint32 `gorm:"column:ctime;NOT NULL" json:"ctime"`
}

func (m *DataSourceTab) TableName() string {
	return "dashboard_data_source_tab"
}

const DataSourceStatusField = "data_source_status"

const DataSourceStatusNormal = 1
const DataSourceStatusDelete = 20
