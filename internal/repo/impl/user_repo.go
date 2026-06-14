package impl

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/cache"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
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
	sf singleflight.Group
}

func NewUserRepo(
	db *gorm.DB,
	rdb *redis.Client,
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

func (u *userRepo) getUserFromCache(
	ctx context.Context,
	key string,
) (*model.User, bool, error) {

	val, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, false, nil
	}

	// 空值缓存
	if val == constants.EmptyUserCacheValue {
		return nil, true, errs.ErrUserNotFound
	}

	user := new(model.User)
	if err := json.Unmarshal([]byte(val), user); err != nil {
		// 缓存损坏，删除重建
		_ = u.rdb.Del(ctx, key).Err()
		return nil, false, nil
	}

	return user, true, nil
}

func (u *userRepo) GetByPhone(
	ctx context.Context,
	phone string,
) (*model.User, error) {

	key := fmt.Sprintf(
		constants.RedisKeyPrefixUserPhone,
		phone,
	)

	// 第一次查缓存
	if user, hit, err := u.getUserFromCache(ctx, key); hit {
		return user, err
	}

	v, err, _ := u.sf.Do(
		key,
		func() (interface{}, error) {

			// Double Check
			if user, hit, err := u.getUserFromCache(ctx, key); hit {
				return user, err
			}

			user := new(model.User)

			err := u.db.WithContext(ctx).
				Where("phone = ?", phone).
				First(user).Error

			if err != nil {

				// 不存在 -> 空值缓存
				if errors.Is(err, gorm.ErrRecordNotFound) {

					_ = u.rdb.Set(
						ctx,
						key,
						constants.EmptyUserCacheValue,
						constants.UserEmptyCacheTTL,
					).Err()

					return nil, errs.ErrUserNotFound
				}

				return nil, err
			}

			// 写入缓存（同时写 user:id 和 user:phone）
			u.setUserCache(ctx, user)

			return user, nil
		},
	)

	if err != nil {
		return nil, err
	}

	user, ok := v.(*model.User)
	if !ok {
		return nil, errs.ErrServiceError
	}

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

	// 第一次查缓存
	if user, hit, err := u.getUserFromCache(ctx, key); hit {
		return user, err
	}

	v, err, _ := u.sf.Do(
		key,
		func() (interface{}, error) {

			// Double Check
			if user, hit, err := u.getUserFromCache(ctx, key); hit {
				return user, err
			}

			user := new(model.User)

			err := u.db.WithContext(ctx).
				Where("id = ?", userID).
				First(user).Error

			if err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {

					_ = u.rdb.Set(
						ctx,
						key,
						constants.EmptyUserCacheValue,
						constants.UserEmptyCacheTTL,
					).Err()

					return nil, errs.ErrUserNotFound
				}

				return nil, err
			}

			u.setUserCache(ctx, user)

			return user, nil
		},
	)

	if err != nil {
		return nil, err
	}

	user, ok := v.(*model.User)
	if !ok {
		return nil, errs.ErrServiceError
	}

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
	pipe.Set(ctx, idKey, bytes, cache.JitterTTL(constants.UserCacheTTL))
	pipe.Set(ctx, phoneKey, bytes, cache.JitterTTL(constants.UserCacheTTL))
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
