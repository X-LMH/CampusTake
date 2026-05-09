package repo

import (
	"CampusTake/internal/enums"
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
	GetByID(ctx context.Context, userID int64) (*model.User, error)
	UpdatePasswordByID(ctx context.Context, userID int64, newPassword string) error
	UpdateAvatarByID(ctx context.Context, userID int64, avatar string) error
	UpdateProfileByID(ctx context.Context, userID int64, nickname string, gender int8) error
	UpdatePhoneByID(ctx context.Context, userID int64, newPhone string) error
	UpdateRoleByID(ctx context.Context, userID int64, newRole enums.RoleType) error
}

type userRepo struct {
	db *gorm.DB
	logx.Logger
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

// Create 创建用户，唯一冲突由数据库索引保障
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

// GetByPhone 根据手机号获取用户
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

// GetByID 根据ID获取用户
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

// UpdatePasswordByID 更新密码
func (u *userRepo) UpdatePasswordByID(ctx context.Context, userID int64, newPassword string) error {
	return u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", newPassword).Error
}

// UpdateAvatarByID 更新头像
func (u *userRepo) UpdateAvatarByID(ctx context.Context, userID int64, avatar string) error {
	return u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("avatar", avatar).Error
}

// UpdateProfileByID 更新个人资料
func (u *userRepo) UpdateProfileByID(ctx context.Context, userID int64, nickname string, gender int8) error {
	updates := map[string]interface{}{
		"gender": gender,
	}
	if nickname != "" {
		updates["nickname"] = nickname
	}

	return u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

// UpdatePhoneByID 更新手机号，需额外处理冲突
func (u *userRepo) UpdatePhoneByID(ctx context.Context, userID int64, newPhone string) error {
	err := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("phone", newPhone).Error

	if err != nil && db.IsDuplicateErr(err) {
		return errs.ErrPhoneAlreadyBound
	}
	return err
}

// UpdateRoleByID 更新角色
func (u *userRepo) UpdateRoleByID(ctx context.Context, userID int64, newRole enums.RoleType) error {
	return u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("role", newRole).Error
}
