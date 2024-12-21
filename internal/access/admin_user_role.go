package access

import (
	"context"

	"gorm.io/gorm"

	"github.com/MainPoser/dst-power/internal/model"
)

type AdminUserRole interface {
	BatchCreate(ctx context.Context, db *gorm.DB, record []*model.AdminUserRole) (int64, error)
	Delete(ctx context.Context, id uint) error
}

func GetAdminUserRoleInterface() AdminUserRole {
	return adminUserRoleImplObj
}

var adminUserRoleImplObj AdminUserRole

type adminUserRoleImpl struct {
	db *gorm.DB
}

func registryAdminUserRole(db *gorm.DB) error {
	impl := &adminUserRoleImpl{db: db}
	adminUserRoleImplObj = impl
	return nil
}

// BatchCreate 批量创建
func (userRole *adminUserRoleImpl) BatchCreate(ctx context.Context, db *gorm.DB, record []*model.AdminUserRole) (int64, error) {
	result := userRole.db.Model(&model.AdminUserRole{}).WithContext(ctx).Create(&record)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// Delete 删除
func (userRole *adminUserRoleImpl) Delete(ctx context.Context, id uint) error {
	return userRole.db.Model(&model.AdminUserRole{}).WithContext(ctx).Delete(&model.AdminUserRole{
		Model: model.Model{ID: id},
	}).Error
}
