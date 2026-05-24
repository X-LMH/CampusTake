package impl

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	errs "CampusTake/pkg/errors"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AddressRepo interface {
	Create(ctx context.Context, address *model.Address) error
	DeleteByID(ctx context.Context, id int64, userID int64) error
	SetDefaultByID(ctx context.Context, id int64, userID int64) error
	GetByIDAndUserID(ctx context.Context, addressID int64, userID int64) (*model.Address, error)
	UpdateByID(ctx context.Context, addressID int64, userID int64, address *model.Address) error
	GetListByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) ([]*model.Address, error)
	ClearDefaultByType(ctx context.Context, userID int64, addrType enums.AddressType) error
	CountByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) (int64, error)
}

type addressRepo struct {
	RepoBase
}

func NewAddressRepo(db *gorm.DB, rdb redis.Cmdable) AddressRepo {
	return &addressRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

// Create 创建新地址
func (a *addressRepo) Create(ctx context.Context, address *model.Address) error {
	return a.db.WithContext(ctx).Create(address).Error
}

// DeleteByID 删除地址，包含所有权校验
func (a *addressRepo) DeleteByID(ctx context.Context, id int64, userID int64) error {
	result := a.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Address{})

	if result.Error != nil {
		return result.Error
	}
	// 删除操作 RowsAffected 为 0 说明地址不存在或不属于该用户
	if result.RowsAffected == 0 {
		return errs.ErrAddressNotFound
	}
	return nil
}

// SetDefaultByID 设为默认地址
func (a *addressRepo) SetDefaultByID(ctx context.Context, id, userID int64) error {
	return a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", enums.AddressIsDefault).Error
}

// GetByIDAndUserID 获取单条地址
func (a *addressRepo) GetByIDAndUserID(ctx context.Context, addressID int64, userID int64) (*model.Address, error) {
	address := new(model.Address)
	err := a.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", addressID, userID).
		First(address).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

// UpdateByID 更新地址内容
func (a *addressRepo) UpdateByID(ctx context.Context, addressID int64, userID int64, address *model.Address) error {
	// 使用 Updates(map) 确保只更新指定字段，且不触发 RowsAffected 误报逻辑
	return a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Updates(map[string]interface{}{
			"contact_name":  address.ContactName,
			"contact_phone": address.ContactPhone,
			"building":      address.Building,
			"room":          address.Room,
			"detail":        address.Detail,
			"type":          address.Type,
		}).Error
}

// GetListByUserIDAndType 获取用户的地址列表（默认地址排在最前）
func (a *addressRepo) GetListByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) ([]*model.Address, error) {
	var addresses []*model.Address
	err := a.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, addrType).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// ClearDefaultByType 重置该类型下的所有默认地址（用于设新默认前的清理）
func (a *addressRepo) ClearDefaultByType(ctx context.Context, userID int64, addrType enums.AddressType) error {
	return a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ? AND is_default = ?", userID, addrType, enums.AddressIsDefault).
		Update("is_default", enums.AddressNotDefault).Error
}

// CountByUserIDAndType 统计数量（用于限制用户最大地址数，如最多10个）
func (a *addressRepo) CountByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) (int64, error) {
	var count int64
	err := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ?", userID, addrType).
		Count(&count).Error
	return count, err
}
