package repo

import (
	"CampusTake/common/errx"
	"CampusTake/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type AddressRepo interface {
	Create(ctx context.Context, address *model.Address) error
	GetListByUserID(ctx context.Context, userID uint64) ([]*model.Address, error)
	DeleteByID(ctx context.Context, id uint64, userID uint64) error
	ClearDefaultByUserID(ctx context.Context, userID uint64) error
	SetDefaultByID(ctx context.Context, id uint64, userID uint64) error
	GetByIDAndUserID(ctx context.Context, addressID uint64, userID uint64) (*model.Address, error)
	UpdateByID(ctx context.Context, addressID uint64, userID uint64, address *model.Address) error
	GetFirstByUserID(ctx context.Context, userID uint64) (*model.Address, error)
}
type addressRepo struct {
	db *gorm.DB
}

func (a *addressRepo) Create(ctx context.Context, address *model.Address) error {
	return a.db.WithContext(ctx).Create(address).Error
}

func (a *addressRepo) GetListByUserID(ctx context.Context, userID uint64) ([]*model.Address, error) {
	var addresses []*model.Address
	err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC").
		Order("created_at ASC").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (a *addressRepo) DeleteByID(ctx context.Context, id uint64, userID uint64) error {
	result := a.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Address{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errx.ErrAddressNotFound
	}
	return nil
}

func (a *addressRepo) ClearDefaultByUserID(ctx context.Context, userID uint64) error {
	return a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND is_default = ?", userID, 1).
		Update("is_default", 0).Error
}

func (a *addressRepo) SetDefaultByID(ctx context.Context, id, userID uint64) error {
	result := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", 1)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errx.ErrAddressNotFound
	}
	return nil
}

func (a *addressRepo) GetByIDAndUserID(ctx context.Context, addressID uint64, userID uint64) (*model.Address, error) {
	address := new(model.Address)
	err := a.db.WithContext(ctx).Where("id = ? AND user_id = ?", addressID, userID).First(address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

func (a *addressRepo) UpdateByID(ctx context.Context, addressID uint64, userID uint64, address *model.Address) error {
	result := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Updates(map[string]interface{}{
			"contact_name":  address.ContactName,
			"contact_phone": address.ContactPhone,
			"building":      address.Building,
			"room":          address.Room,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errx.ErrAddressNotFound
	}
	return nil
}

func (a *addressRepo) GetFirstByUserID(ctx context.Context, userID uint64) (*model.Address, error) {
	address := new(model.Address)
	err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		First(address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}
