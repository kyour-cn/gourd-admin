// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"encoding/json"
)

const TableNameMenu = "menu"

// Menu 菜单
type Menu struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AppID     int32     `gorm:"column:app_id;not null;comment:应用ID" json:"app_id"`           // 应用ID
	Pid       int32     `gorm:"column:pid;not null;comment:上级ID" json:"pid"`                 // 上级ID
	Name      string    `gorm:"column:name;not null;comment:别名" json:"name"`                 // 别名
	Title     string    `gorm:"column:title;not null;comment:显示名称" json:"title"`             // 显示名称
	Type      string    `gorm:"column:type;not null;comment:类型" json:"type"`                 // 类型
	Path      string    `gorm:"column:path;not null;comment:路由地址" json:"path"`               // 路由地址
	Component string    `gorm:"column:component;not null;comment:组件地址" json:"component"`     // 组件地址
	Status    int32     `gorm:"column:status;not null;default:1;comment:是否启用" json:"status"` // 是否启用
	Sort      int32     `gorm:"column:sort;not null;comment:排序" json:"sort"`                 // 排序
	Meta      string    `gorm:"column:meta;not null;comment:meta路由参数" json:"meta"`           // meta路由参数
	App       App       `gorm:"foreignKey:app_id;references:id" json:"app"`
	MenuApi   []MenuAPI `gorm:"foreignKey:menu_id;references:id" json:"menu_api"`
}

// MarshalBinary 支持json序列化
func (m *Menu) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *Menu) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// TableName Menu's table name
func (*Menu) TableName() string {
	return TableNameMenu
}
