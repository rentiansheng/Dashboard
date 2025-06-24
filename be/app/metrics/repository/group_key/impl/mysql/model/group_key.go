package model

type GroupKeyTab struct {
	ID          uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DisplayName string `gorm:"column:display_name;NOT NULL" json:"name"`
	RootID      string `gorm:"column:root_id;NOT NULL" json:"root_id"`
	ParentID    uint64 `gorm:"column:parent_id;NOT NULL" json:"parent_id"`
	Remark      string `gorm:"column:remark" json:"remark"`
	Status      uint   `gorm:"column:group_key_status;NOT NULL" json:"status"`
	OrderIdx    uint16 `gorm:"column:order_index;NOT NULL" json:"order_index"`

	Mtime uint32 `gorm:"column:mtime;NOT NULL" json:"mtime"`
	Ctime uint32 `gorm:"column:ctime;NOT NULL" json:"ctime"`
}

func (m *GroupKeyTab) TableName() string {
	return "dashboard_group_key_tab"
}
