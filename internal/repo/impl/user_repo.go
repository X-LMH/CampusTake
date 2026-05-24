package impl

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
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
	RepoBase
}

func NewUserRepo(
	db *gorm.DB,
	rdb redis.Cmdable,
) UserRepo {
	return &userRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

// =========================
// Create 创建用户
// =========================

func (u *userRepo) Create(
	ctx context.Context,
	user *model.User,
) error {

	err := u.db.WithContext(ctx).
		Create(user).Error

	if err != nil {

		if db.IsDuplicateErr(err) {
			return errs.ErrUserExist
		}

		return err
	}

	return nil
}

// =========================
// GetByPhone 根据手机号获取用户
// =========================

func (u *userRepo) GetByPhone(
	ctx context.Context,
	phone string,
) (*model.User, error) {
	key := fmt.Sprintf(constants.RedisKeyPrefixUserPhone, phone)

	val, err := u.rdb.Get(ctx, key).Result()
	if err == nil {
		user := new(model.User)
		if json.Unmarshal([]byte(val), user) == nil {
			return user, nil
		}
		_ = u.rdb.Del(ctx, key).Err()
	}

	user := new(model.User)

	err = u.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	u.setUserCache(ctx, user)

	return user, nil
}

// =========================
// GetByID 根据ID获取用户
// =========================

func (u *userRepo) GetByID(
	ctx context.Context,
	userID int64,
) (*model.User, error) {

	key := fmt.Sprintf(
		constants.RedisKeyPrefixUser,
		userID,
	)

	// =========================
	// 1. 查询 Redis
	// =========================

	val, err := u.rdb.Get(ctx, key).Result()

	if err == nil {

		user := new(model.User)

		if json.Unmarshal([]byte(val), user) == nil {
			return user, nil
		}

		// JSON 解析失败
		// 删除脏缓存
		_ = u.rdb.Del(ctx, key).Err()
	}

	// redis 异常降级
	if err != nil && !errors.Is(err, redis.Nil) {
		// ignore
	}

	// =========================
	// 2. 查询 MySQL
	// =========================

	user := new(model.User)

	err = u.db.WithContext(ctx).
		Where("id = ?", userID).
		First(user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	// =========================
	// 3. 回填 Redis
	// =========================

	u.setUserCache(ctx, user)

	return user, nil
}

func (u *userRepo) setUserCache(ctx context.Context, user *model.User) {
	bytes, err := json.Marshal(user)
	if err != nil {
		return
	}

	idKey := fmt.Sprintf(constants.RedisKeyPrefixUser, user.ID)
	phoneKey := fmt.Sprintf(constants.RedisKeyPrefixUserPhone, user.Phone)

	pipe := u.rdb.Pipeline()
	pipe.Set(ctx, idKey, bytes, constants.UserCacheTTL)
	pipe.Set(ctx, phoneKey, bytes, constants.UserCacheTTL)
	_, _ = pipe.Exec(ctx)
}

// =========================
// UpdatePasswordByID 更新密码
// =========================

func (u *userRepo) UpdatePasswordByID(
	ctx context.Context,
	userID int64,
	newPassword string,
) error {

	tx := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", newPassword)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	u.deleteUserCache(ctx, userID)

	return nil
}

// =========================
// UpdateAvatarByID 更新头像
// =========================

func (u *userRepo) UpdateAvatarByID(
	ctx context.Context,
	userID int64,
	avatar string,
) error {

	tx := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("avatar", avatar)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	u.deleteUserCache(ctx, userID)

	return nil
}

// =========================
// UpdateProfileByID 更新个人资料
// =========================

func (u *userRepo) UpdateProfileByID(
	ctx context.Context,
	userID int64,
	nickname string,
	gender int8,
) error {

	updates := map[string]interface{}{
		"gender": gender,
	}

	if nickname != "" {
		updates["nickname"] = nickname
	}

	tx := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	u.deleteUserCache(ctx, userID)

	return nil
}

// =========================
// UpdatePhoneByID 更新手机号
// =========================

func (u *userRepo) UpdatePhoneByID(
	ctx context.Context,
	userID int64,
	newPhone string,
) error {
	var oldPhone string
	_ = u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Select("phone").
		Scan(&oldPhone).Error

	tx := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("phone", newPhone)

	if tx.Error != nil {

		if db.IsDuplicateErr(tx.Error) {
			return errs.ErrPhoneAlreadyBound
		}

		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	u.deleteUserCache(ctx, userID)
	if oldPhone != "" && oldPhone != newPhone {
		_ = u.rdb.Del(ctx, fmt.Sprintf(constants.RedisKeyPrefixUserPhone, oldPhone)).Err()
	}

	return nil
}

// =========================
// UpdateRoleByID 更新角色
// =========================

func (u *userRepo) UpdateRoleByID(
	ctx context.Context,
	userID int64,
	newRole enums.RoleType,
) error {

	tx := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("role", newRole)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	u.deleteUserCache(ctx, userID)

	return nil
}

// =========================
// deleteUserCache 删除用户缓存
// =========================

func (u *userRepo) deleteUserCache(
	ctx context.Context,
	userID int64,
) {
	keys := []string{
		fmt.Sprintf(
			constants.RedisKeyPrefixUser,
			userID,
		),
	}

	var phone string
	if err := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Select("phone").
		Scan(&phone).Error; err == nil && phone != "" {
		keys = append(keys, fmt.Sprintf(constants.RedisKeyPrefixUserPhone, phone))
	}

	_ = u.rdb.Del(ctx, keys...).Err()
}
