package repo

import (
	"CampusTake/common/db"
	"CampusTake/common/errx"
	"context"
	"errors"

	"CampusTake/internal/model"

	"gorm.io/gorm"
)

// ✅ 接口
type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	UpdatePasswordByID(ctx context.Context, userID uint64, newPassword string) error
}

// ✅ 实现
type userRepo struct {
	db *gorm.DB
}

func (u *userRepo) Create(ctx context.Context, user *model.User) error {

	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if db.IsDuplicateErr(err) {
			return errx.ErrUserExist
		}
		return err
	}

	return nil
}

func (u *userRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	user := new(model.User)
	err := u.db.WithContext(ctx).Where("phone = ?", phone).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userRepo) UpdatePasswordByID(ctx context.Context, userID uint64, newPassword string) error {
	res := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND password <> ?", userID, newPassword).
		Updates(map[string]interface{}{
			"password": newPassword,
		})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		var count int64
		u.db.WithContext(ctx).Model(&model.User{}).
			Where("userID = ?", userID).
			Count(&count)

		if count == 0 {
			return errx.ErrUserNotFound
		}
		return errx.ErrPasswordNoChange
	}

	return nil
}
