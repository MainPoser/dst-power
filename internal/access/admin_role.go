package access

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/MainPoser/dst-power/internal/dto"
	"github.com/MainPoser/dst-power/internal/model"
)

type AdminRole interface {
	List(ctx context.Context, params *dto.GetAdminRoleListRequest) ([]*model.AdminRole, int64, error)
	Create(ctx context.Context, adminRole *model.AdminRole) error
	Get(ctx context.Context, id int64) (*model.AdminRole, error)
	Update(ctx context.Context, adminRole *model.AdminRole) error
	DeleteBatch(ctx context.Context, ids []int64) error
}

func GetAdminRoleInterface() AdminRole {
	return adminRoleImplObj
}

var adminRoleImplObj AdminRole

type adminRoleImpl struct {
	db *gorm.DB
}

func registryAdminRole(db *gorm.DB) error {
	impl := &adminRoleImpl{db: db}
	adminRoleImplObj = impl
	return nil
}

// List 列表
func (m *adminRoleImpl) List(ctx context.Context, params *dto.GetAdminRoleListRequest) ([]*model.AdminRole, int64, error) {
	tx := m.db.Model(&model.AdminRole{})
	if params.Name != "" {
		tx.Where("name like ?", fmt.Sprintf("%%%s%%", params.Name))
	}
	if params.Status != "" {
		// TODO convert.Int(params.Status)
		//tx.Where("status = ?", convert.Int(params.Status))
	}
	var count int64
	tx.Count(&count)

	if params.Page > 0 && params.PageSize > 0 {
		tx.Offset(GetPageOffset(int(params.Page), int(params.PageSize))).Limit(int(params.PageSize))
	}

	var list []*model.AdminRole
	if err := tx.Order("id asc").WithContext(ctx).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// Create 创建
func (m *adminRoleImpl) Create(ctx context.Context, adminRole *model.AdminRole) error {
	return m.db.Model(&model.AdminRole{}).WithContext(ctx).Create(adminRole).Error
}

// Get 查询
func (m *adminRoleImpl) Get(ctx context.Context, id int64) (*model.AdminRole, error) {
	adminRole := &model.AdminRole{}
	if err := m.db.Model(&model.AdminRole{}).WithContext(ctx).Where("id = ? ", id).First(adminRole).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return adminRole, nil
}

// Update 更新
func (m *adminRoleImpl) Update(ctx context.Context, adminRole *model.AdminRole) error {
	return m.db.Model(&model.AdminRole{}).WithContext(ctx).Save(adminRole).Error
}

// DeleteBatch 批量删除
func (m *adminRoleImpl) DeleteBatch(ctx context.Context, ids []int64) error {
	bind := &model.AdminRole{}
	return m.db.Model(bind).WithContext(ctx).Where("id in ?", ids).Delete(bind).Error
}
