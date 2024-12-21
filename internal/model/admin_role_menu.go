package model

type AdminRoleMenu struct {
	Model
	RoleId int `gorm:"column:role_id;default:0;NOT NULL" json:"role_id"` // 角色id
	MenuId int `gorm:"column:menu_id;default:0;NOT NULL" json:"menu_id"` // 菜单id
}

func (roleMenu *AdminRoleMenu) TableName() string {
	return "admin_role_menu"
}
