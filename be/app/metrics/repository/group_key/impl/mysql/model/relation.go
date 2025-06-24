package model

type RelationTab struct {
	ID uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`

	GroupKeyId  uint64 `gorm:"column:group_key_id" json:"group_key_id"`
	DataGroupId uint64 `gorm:"column:data_group_id" json:"data_group_id"`
	Status      uint   `gorm:"column:data_source_status;NOT NULL" json:"status"`

	Mtime uint32 `gorm:"column:mtime;NOT NULL" json:"mtime"`
	Ctime uint32 `gorm:"column:ctime;NOT NULL" json:"ctime"`
}

func (m *RelationTab) TableName() string {
	return "dashboard_relation_tab"
}

const StatusNormal = 1
const StatusDelete = 20
