package repo

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	errs "CampusTake/pkg/errors"
	"context"
	"errors"

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
	db *gorm.DB
}

func (a *addressRepo) Create(ctx context.Context, address *model.Address) error {
	return a.db.WithContext(ctx).Create(address).Error
}

func (a *addressRepo) DeleteByID(ctx context.Context, id int64, userID int64) error {
	result := a.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Address{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrAddressNotFound
	}
	return nil
}

func (a *addressRepo) SetDefaultByID(ctx context.Context, id, userID int64) error {
	result := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", enums.AddressIsDefault)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrAddressNotFound
	}
	return nil
}

func (a *addressRepo) GetByIDAndUserID(ctx context.Context, addressID int64, userID int64) (*model.Address, error) {
	address := new(model.Address)
	err := a.db.WithContext(ctx).Where("id = ? AND user_id = ?", addressID, userID).First(address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

func (a *addressRepo) UpdateByID(ctx context.Context, addressID int64, userID int64, address *model.Address) error {
	result := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Updates(map[string]interface{}{
			"contact_name":  address.ContactName,
			"contact_phone": address.ContactPhone,
			"building":      address.Building,
			"room":          address.Room,
			"detail":        address.Detail,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrAddressNotFound
	}
	return nil
}

func (a *addressRepo) GetListByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) ([]*model.Address, error) {
	var addresses []*model.Address
	err := a.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, addrType).
		Order("is_default DESC").
		Order("created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

func (a *addressRepo) ClearDefaultByType(ctx context.Context, userID int64, addrType enums.AddressType) error {
	return a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ? AND is_default = ?", userID, addrType, enums.AddressIsDefault).
		Update("is_default", enums.AddressNotDefault).Error
}

func (a *addressRepo) CountByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) (int64, error) {
	var count int64
	err := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ?", userID, addrType).
		Count(&count).Error
	return count, err
}
