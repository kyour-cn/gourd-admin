// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"encoding/json"
)

const TableNameRule = "rule"

// Rule 权限规则表
type Rule struct {
	ID      int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AppID   int32  `gorm:"column:app_id;not null;comment:应用ID" json:"app_id"`              // 应用ID
	Name    string `gorm:"column:name;not null;comment:名字" json:"name"`                    // 名字
	Alias_  string `gorm:"column:alias;not null;comment:英文别名" json:"alias"`                // 英文别名
	Path    string `gorm:"column:path;not null;comment:规则" json:"path"`                    // 规则
	Pid     int32  `gorm:"column:pid;not null;comment:上级Id" json:"pid"`                    // 上级Id
	Status  int32  `gorm:"column:status;not null;default:1;comment:状态 0:1" json:"status"`  // 状态 0:1
	Sort    int32  `gorm:"column:sort;not null;comment:排序" json:"sort"`                    // 排序
	AddonID int32  `gorm:"column:addon_id;not null;comment:插件ID 为0=不验证插件" json:"addon_id"` // 插件ID 为0=不验证插件
}

// MarshalBinary 支持json序列化
func (m *Rule) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *Rule) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// TableName Rule's table name
func (*Rule) TableName() string {
	return TableNameRule
}
