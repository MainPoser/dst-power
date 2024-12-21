package access

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/MainPoser/dst-power/internal/dto"
	"github.com/MainPoser/dst-power/internal/model"
	"github.com/MainPoser/dst-power/pkg/apis/vo"
)

type AdminMenu interface {
	List(ctx context.Context, params *dto.GetAdminMenuListRequest) ([]*model.AdminMenu, int64, error)
	Get(ctx context.Context, id uint) (*model.AdminMenu, error)
	Create(ctx context.Context, adminMenu *model.AdminMenu) error
	Update(ctx context.Context, adminMenu *model.AdminMenu) error
	Delete(ctx context.Context, id uint) error
	GetAdminMenuListByUid(ctx context.Context, userId uint) interface{}
}

func GetAdminMenuInterface() AdminMenu {
	return adminMenuImplObj
}

var adminMenuImplObj AdminMenu

type adminMenuImpl struct {
	db *gorm.DB
}

func registryAdminMenu(db *gorm.DB) error {
	impl := &adminMenuImpl{db: db}
	adminMenuImplObj = impl
	return nil
}

func (m *adminMenuImpl) List(ctx context.Context, params *dto.GetAdminMenuListRequest) ([]*model.AdminMenu, int64, error) {
	tx := m.db.Model(&model.AdminMenu{})

	if params.Status > 0 {
		tx.Where("status = ?", params.Status)
	}

	var count int64
	tx.Count(&count)

	if params.Page > 0 && params.PageSize > 0 {
		tx = tx.Offset(GetPageOffset(int(params.Page), int(params.PageSize))).Limit(int(params.PageSize))
	}

	var list []*model.AdminMenu
	if err := tx.Order("sort asc").WithContext(ctx).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// Get 根据条件查询单条数据
func (m *adminMenuImpl) Get(ctx context.Context, id uint) (*model.AdminMenu, error) {
	tx := m.db.Model(&model.AdminMenu{})

	tx.Where("id = ? ", id)

	var adminMenu *model.AdminMenu
	if err := tx.WithContext(ctx).First(&adminMenu).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return adminMenu, err
	}
	return adminMenu, nil
}

// Create 插入数据
func (m *adminMenuImpl) Create(ctx context.Context, adminMenu *model.AdminMenu) error {
	return m.db.Model(&model.AdminMenu{}).WithContext(ctx).Create(adminMenu).Error
}

// Update 更新
func (m *adminMenuImpl) Update(ctx context.Context, adminMenu *model.AdminMenu) error {
	return m.db.Model(&model.AdminMenu{}).WithContext(ctx).Save(adminMenu).Error
}

// Delete 删除
func (m *adminMenuImpl) Delete(ctx context.Context, id uint) error {
	return m.db.Model(&model.AdminMenu{}).WithContext(ctx).Where("id = ?", id).Delete(&m).Error
}

func makeTree(menu []*model.AdminMenu, tn *vo.MenuTreeNode) {
	for _, c := range menu {
		if (tn.AdminMenu == nil && c.Pid == 0) || (tn.AdminMenu != nil && c.Pid == tn.Model.ID) {
			child := &vo.MenuTreeNode{}
			child.AdminMenu = c
			tn.Children = append(tn.Children, child)
			makeTree(menu, child)
		}
	}
}

// GetAdminMenuListByUid 通过用户id获取菜单列表
func (m *adminMenuImpl) GetAdminMenuListByUid(ctx context.Context, userId uint) interface{} {
	var list []*model.AdminMenu
	m.db.Model(&model.AdminMenu{}).WithContext(ctx).
		Joins("left join admin_role_menu as rm on admin_menu.id = rm.menu_id").
		Joins("left join admin_user_role as ur on rm.role_id = ur.role_id").
		Order("admin_menu.sort asc, admin_menu.id asc").
		Where("ur.user_id=? AND admin_menu.type=0 AND admin_menu.`status`=1 AND admin_menu.is_show=1", userId).
		Select("admin_menu.*").
		Find(&list)
	var menuNode vo.MenuTreeNode
	makeTree(list, &menuNode)

	return menuNode.Children
}
