package model

type DataGroupTab struct {
	ID          uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DisplayName string `gorm:"column:display_name;NOT NULL" json:"name"`
	Status      uint   `gorm:"column:data_source_status;NOT NULL" json:"status"`
	Mtime       uint32 `gorm:"column:mtime;NOT NULL" json:"mtime"`
	Ctime       uint32 `gorm:"column:ctime;NOT NULL" json:"ctime"`
}

func (m *DataGroupTab) TableName() string {
	return "dashboard_data_group_tab"
}
