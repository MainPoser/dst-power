package vo

import "github.com/MainPoser/dst-power/internal/model"

// MenuTreeNode 菜单Vo
type MenuTreeNode struct {
	*model.AdminMenu
	Children []*MenuTreeNode `json:"children"` // 子菜单
}
