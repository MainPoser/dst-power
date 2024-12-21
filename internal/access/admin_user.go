package access

import (
	"context"
	"errors"
	"fmt"
	"github.com/MainPoser/dst-power/internal/dto"
	"github.com/MainPoser/dst-power/internal/model"
	"gorm.io/gorm"
)

type AdminUser interface {
	Get(ctx context.Context, id uint) (*model.AdminUser, error)
	GetWithName(ctx context.Context, name string) (*model.AdminUser, error)
	Update(ctx context.Context, adminUser *model.AdminUser) error
	Create(ctx context.Context, adminUser *model.AdminUser) (*model.AdminUser, error)
	List(ctx context.Context, params *dto.GetAdminUserListRequest) ([]*model.AdminUser, int64, error)
	Delete(ctx context.Context, id uint) error
}

func GetAdminUserInterface() AdminUser {
	return adminUserImplObj
}

var adminUserImplObj AdminUser

type adminUserImpl struct {
	db *gorm.DB
}

func registryAdminUser(db *gorm.DB) error {
	impl := &adminUserImpl{db: db}
	adminUserImplObj = impl
	return nil
}

// GetWithName 根据名称获取一条记录
func (a *adminUserImpl) GetWithName(ctx context.Context, name string) (*model.AdminUser, error) {

	var admin *model.AdminUser
	if err := a.db.Model(&model.AdminUser{}).Where("username = ?", name).
		WithContext(ctx).First(&admin).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return admin, err
	}
	return admin, nil
}

// Get 获取一条记录
func (a *adminUserImpl) Get(ctx context.Context, id uint) (*model.AdminUser, error) {

	tx := a.db.Model(&model.AdminUser{}).Where("id = ?", id)

	var admin *model.AdminUser
	if err := tx.WithContext(ctx).First(&admin).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return admin, err
	}
	return admin, nil
}

// Update 更新
func (a *adminUserImpl) Update(ctx context.Context, adminUser *model.AdminUser) error {
	if err := a.db.Model(&model.AdminUser{}).Where("id = ?", adminUser.ID).WithContext(ctx).Updates(adminUser).Error; err != nil {
		return err
	}
	return nil
}

// Create 创建
func (a *adminUserImpl) Create(ctx context.Context, adminUser *model.AdminUser) (*model.AdminUser, error) {
	if err := a.db.Model(&model.AdminUser{}).WithContext(ctx).Create(adminUser).Error; err != nil {
		return nil, err
	}
	return adminUser, nil
}

// List 列表
func (a *adminUserImpl) List(ctx context.Context, params *dto.GetAdminUserListRequest) ([]*model.AdminUser, int64, error) {
	tx := a.db.Model(&model.AdminUser{})
	if params.Name != "" {
		tx.Where("name like ?", fmt.Sprintf("%%%s%%", params.Name))
	}
	var count int64
	tx.Count(&count)

	if params.Page > 0 && params.PageSize > 0 {
		tx.Offset(GetPageOffset(int(params.Page), int(params.PageSize))).Limit(int(params.PageSize))
	}

	var list []*model.AdminUser
	if err := tx.WithContext(ctx).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// Delete 删除
func (a *adminUserImpl) Delete(ctx context.Context, id uint) error {
	bind := &model.AdminUser{}
	return a.db.Model(bind).WithContext(ctx).Where("id = ?", id).Delete(bind).Error
}
