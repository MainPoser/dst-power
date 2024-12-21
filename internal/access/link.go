package access

import (
	"context"
	"errors"
	"fmt"
	"github.com/MainPoser/dst-power/internal/dto"
	"github.com/MainPoser/dst-power/internal/model"
	"gorm.io/gorm"
)

type Link interface {
	List(ctx context.Context, params *dto.GetLinkListRequest) ([]*model.Link, int64, error)
	Get(ctx context.Context, id uint) (*model.Link, error)
	Create(ctx context.Context, link *model.Link) error
	Update(ctx context.Context, link *model.Link) error
	Delete(ctx context.Context, id uint) error
}

func GetLinkInterface() Link {
	return linkImplObj
}

var linkImplObj Link

type linkImpl struct {
	db *gorm.DB
}

func registryLink(db *gorm.DB) error {
	impl := &linkImpl{db: db}
	linkImplObj = impl
	return nil
}

func (l *linkImpl) List(ctx context.Context, params *dto.GetLinkListRequest) ([]*model.Link, int64, error) {
	tx := l.db.Model(&model.Link{})
	if params.Status > 0 {
		tx.Where("status = ?", params.Status)
	}
	if params.Name != "" {
		tx.Where("name like ?", fmt.Sprintf("%%%s%%", params.Name))
	}
	var count int64
	tx.Count(&count)

	if params.Page > 0 && params.PageSize > 0 {
		tx.Offset(GetPageOffset(int(params.Page), int(params.PageSize))).Limit(int(params.PageSize))
	}

	var list []*model.Link
	if err := tx.Order("sort asc").WithContext(ctx).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// Get 根据条件查询单条数据
func (l *linkImpl) Get(ctx context.Context, id uint) (*model.Link, error) {
	var link *model.Link
	if err := l.db.Model(&model.Link{}).Where("id = ? ", id).WithContext(ctx).First(&link).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return link, err
	}
	return link, nil
}

// Create 插入数据
func (l *linkImpl) Create(ctx context.Context, link *model.Link) error {
	return l.db.Model(&model.Link{}).WithContext(ctx).Create(link).Error
}

// Update 更新
func (l *linkImpl) Update(ctx context.Context, link *model.Link) error {
	return l.db.Model(&model.Link{}).WithContext(ctx).Save(link).Error
}

// Delete 删除
func (l *linkImpl) Delete(ctx context.Context, id uint) error {
	bind := &model.Link{}
	return l.db.Model(bind).WithContext(ctx).Where("id = ?", id).Delete(bind).Error
}
