package repo

import (
	"CampusTake/internal/model"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	UpdatePasswordByID(ctx context.Context, userID int64, newPassword string) error
	GetByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateAvatarByID(ctx context.Context, userID int64, avatar string) error
	UpdateProfileByID(ctx context.Context, userID int64, nickname string, gender int8) error
	UpdatePhoneByID(ctx context.Context, userID int64, newPhone string) error
}

type userRepo struct {
	db *gorm.DB
	logx.Logger
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (u *userRepo) Create(ctx context.Context, user *model.User) error {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if db.IsDuplicateErr(err) {
			return errs.ErrUserExist
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
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userRepo) UpdatePasswordByID(ctx context.Context, userID int64, newPassword string) error {
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
			Where("id = ?", userID).
			Count(&count)

		if count == 0 {
			return errs.ErrUserNotFound
		}
		return errs.ErrPasswordNoChange
	}

	return nil
}

func (u *userRepo) GetByID(ctx context.Context, userID int64) (*model.User, error) {
	user := new(model.User)
	err := u.db.WithContext(ctx).Where("id = ?", userID).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userRepo) UpdateAvatarByID(ctx context.Context, userID int64, avatar string) error {
	res := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND avatar <> ?", userID, avatar).
		Updates(map[string]interface{}{"avatar": avatar})

	if res.Error != nil {
		return res.Error
	}

	// 如果没有行受影响，检查用户是否存在
	if res.RowsAffected == 0 {
		var count int64
		u.db.WithContext(ctx).Model(&model.User{}).
			Where("id = ?", userID).
			Count(&count)

		if count == 0 {
			return errs.ErrUserNotFound
		}
		// avatar 没有变更，视为成功（也可以返回特定错误，但项目中无 ErrAvatarNoChange）
	}

	return nil
}

func (u *userRepo) UpdateProfileByID(ctx context.Context, userID int64, nickname string, gender int8) error {
	updates := map[string]interface{}{}

	if nickname != "" {
		updates["nickname"] = nickname
	}
	updates["gender"] = gender

	res := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		var count int64
		u.db.WithContext(ctx).Model(&model.User{}).
			Where("id = ?", userID).
			Count(&count)

		if count == 0 {
			return errs.ErrUserNotFound
		}
	}

	return nil
}

func (u *userRepo) UpdatePhoneByID(ctx context.Context, userID int64, newPhone string) error {
	res := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND phone <> ?", userID, newPhone).
		Update("phone", newPhone)

	if res.Error != nil {
		if db.IsDuplicateErr(res.Error) {
			return errs.ErrPhoneAlreadyBound
		}
		return res.Error
	}

	return nil
}
