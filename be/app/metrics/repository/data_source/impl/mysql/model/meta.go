package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rentiansheng/dashboard/app/metrics/define"
)

type DataSourceMetaTab struct {
	ID             uint                     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DataSourceID   uint                     `gorm:"column:data_source_id" json:"data_source_id"`
	FieldName      string                   `gorm:"column:field_name" json:"name"`
	DisplayName    string                   `gorm:"column:display_name;NOT NULL" json:"display_name"`
	DataType       string                   `gorm:"column:data_type;NOT NULL" json:"data_type"`
	OutputDataType string                   `gorm:"column:output_data_type;NOT NULL" json:"output_data_type"`
	FieldTips      string                   `gorm:"column:field_tips" json:"field_tips"`
	Action         DataSourceMetaEnumAction `gorm:"column:action;NOT NULL" json:"action"`
	Enum           DataSourceMetaEnum       `gorm:"column:enum;NOT NULL" json:"enum"`
	Formatter      string                   `gorm:"column:formatter;NOT NUll" json:"formatter"`
	Template       string                   `gorm:"column:template;NOT NUll" json:"template"`
	Nested         string                   `gorm:"column:nested_path" json:"nested_path"`
	Status         uint                     `gorm:"column:data_source_meta_status;NOT NULL" json:"status"`
	Mtime          uint32                   `gorm:"column:mtime;NOT NULL" json:"mtime"`
	Ctime          uint32                   `gorm:"column:ctime;NOT NULL" json:"ctime"`
}

type DataSourceMetaEnumAction define.MetricDataSourceMetaEnumAction

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *DataSourceMetaEnumAction) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, m)
	return err
}

// Value return json value, implement driver.Valuer interface
func (m DataSourceMetaEnumAction) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type DataSourceMetaEnum define.MetricDataSourceMetaEnum

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *DataSourceMetaEnum) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, m)
	return err
}

// Value return json value, implement driver.Valuer interface
func (m DataSourceMetaEnum) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m DataSourceMetaTab) TableName() string {
	return "dashboard_data_source_meta_tab"
}
