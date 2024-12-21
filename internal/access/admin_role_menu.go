package access

import (
	"context"
	"github.com/MainPoser/dst-power/internal/dto"
	"github.com/MainPoser/dst-power/internal/model"
	"gorm.io/gorm"
)

type AdminRoleMenu interface {
	List(ctx context.Context, param *dto.AdminRoleMenuSelector) ([]*model.AdminRoleMenu, error)
	Delete(ctx context.Context, id uint) error
	BatchCreate(ctx context.Context, record []*AdminRoleMenu) (int64, error)
}

func GetAdminRoleMenuInterface() AdminRoleMenu {
	return adminRoleMenuImplObj
}

var adminRoleMenuImplObj AdminRoleMenu

type adminRoleMenuImpl struct {
	db *gorm.DB
}

func registryAdminRoleMenu(db *gorm.DB) error {
	impl := &adminRoleMenuImpl{db: db}
	adminRoleMenuImplObj = impl
	return nil
}

// List 列表
func (roleMenu *adminRoleMenuImpl) List(ctx context.Context, param *dto.AdminRoleMenuSelector) ([]*model.AdminRoleMenu, error) {

	tx := roleMenu.db.Model(&model.AdminRoleMenu{})
	if param.RoleId != 0 {
		tx.Where("role_id = ?", param.RoleId)
	}

	var list []*model.AdminRoleMenu
	if err := tx.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Delete 删除
func (roleMenu *adminRoleMenuImpl) Delete(ctx context.Context, id uint) error {
	return roleMenu.db.Model(&model.AdminRoleMenu{}).Where("id = ?", id).WithContext(ctx).Delete(&model.AdminRoleMenu{}).Error
}

// BatchCreate 批量创建
func (roleMenu *adminRoleMenuImpl) BatchCreate(ctx context.Context, record []*AdminRoleMenu) (int64, error) {
	result := roleMenu.db.Model(&model.AdminRoleMenu{}).WithContext(ctx).Create(&record)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
