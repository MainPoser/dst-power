package model

type AdminRole struct {
	Model
	Name   string `gorm:"column:name;NOT NULL" json:"name"`      // 角色名称
	Code   string `gorm:"column:code;NOT NULL" json:"code"`      // 角色编码
	Sort   int    `gorm:"column:sort;default:0" json:"sort"`     // 排序
	Status int    `gorm:"column:status;default:1" json:"status"` // 状态 1-启用2-禁用
}

func (role *AdminRole) TableName() string {
	return "admin_role"
}
