package impl

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	errs "CampusTake/pkg/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

func NewAddressRepo(db *gorm.DB, rdb *redis.Client) AddressRepo {
	return &addressRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

// Create 创建新地址（新地址不需要写缓存，等查询时自然会触发写入）
func (a *addressRepo) Create(ctx context.Context, address *model.Address) error {
	return a.db.WithContext(ctx).Create(address).Error
}

// DeleteByID 删除地址，包含所有权校验及缓存清除
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

	// 数据库删除成功后，清除 Redis 缓存
	cacheKey := fmt.Sprintf(constants.RedisKeyPrefixAddress, id)
	_ = a.rdb.Del(ctx, cacheKey).Err() // 这里的 error 可以选择记录日志或忽略

	return nil
}

// SetDefaultByID 设为默认地址，并清除该地址缓存
func (a *addressRepo) SetDefaultByID(ctx context.Context, id, userID int64) error {
	err := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", enums.AddressIsDefault).Error

	if err != nil {
		return err
	}

	// 清除该地址缓存
	cacheKey := fmt.Sprintf(constants.RedisKeyPrefixAddress, id)
	_ = a.rdb.Del(ctx, cacheKey).Err()

	return nil
}

// GetByIDAndUserID 获取单条地址（带高并发缓存保护）
func (a *addressRepo) GetByIDAndUserID(ctx context.Context, addressID int64, userID int64) (*model.Address, error) {
	cacheKey := fmt.Sprintf(constants.RedisKeyPrefixAddress, addressID)

	// 1. 尝试从 Redis 获取缓存
	val, err := a.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		address := new(model.Address)
		if err := json.Unmarshal([]byte(val), address); err == nil {
			// 严格越权校验：即使缓存命中，也要确保请求的 userID 匹配
			if address.UserID == userID {
				return address, nil
			}
			return nil, errs.ErrAddressNotFound
		}
	}

	// 2. 缓存未命中，查询数据库
	address := new(model.Address)
	err = a.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", addressID, userID).
		First(address).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 可选防穿透逻辑：如果数据库没有，可在 Redis 写入短时间的空值，如 a.rdb.Set(ctx, cacheKey, "", 1*time.Minute)
			return nil, errs.ErrAddressNotFound
		}
		return nil, err
	}

	// 3. 异步/同步回写缓存（包含序列化）
	if data, err := json.Marshal(address); err == nil {
		_ = a.rdb.Set(ctx, cacheKey, data, constants.AddressCacheTTL).Err()
	}

	return address, nil
}

// UpdateByID 更新地址内容，并使缓存失效
func (a *addressRepo) UpdateByID(ctx context.Context, addressID int64, userID int64, address *model.Address) error {
	err := a.db.WithContext(ctx).
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

	if err != nil {
		return err
	}

	// 强制删除旧缓存，保证后续读取到最新 DB 数据
	cacheKey := fmt.Sprintf(constants.RedisKeyPrefixAddress, addressID)
	_ = a.rdb.Del(ctx, cacheKey).Err()

	return nil
}

// GetListByUserIDAndType 获取用户的地址列表（直接查库，因为列表排序复杂且长度短，单条数据已走缓存）
func (a *addressRepo) GetListByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) ([]*model.Address, error) {
	var addresses []*model.Address
	err := a.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, addrType).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// ClearDefaultByType 重置该类型下的所有默认地址（精准清理受影响的地址缓存）
func (a *addressRepo) ClearDefaultByType(ctx context.Context, userID int64, addrType enums.AddressType) error {
	// 1. 先查出当前已经处于默认状态的地址 ID（通常只有1条），用于后续精准清除缓存
	var affectedIDs []int64
	err := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ? AND is_default = ?", userID, addrType, enums.AddressIsDefault).
		Pluck("id", &affectedIDs).Error
	if err != nil {
		return err
	}

	// 2. 执行数据库更新
	err = a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ? AND is_default = ?", userID, addrType, enums.AddressIsDefault).
		Update("is_default", enums.AddressNotDefault).Error
	if err != nil {
		return err
	}

	// 3. 循环删除受影响地址的 Redis 缓存，防止 `is_default` 状态在缓存中脏读
	for _, id := range affectedIDs {
		cacheKey := fmt.Sprintf(constants.RedisKeyPrefixAddress, id)
		_ = a.rdb.Del(ctx, cacheKey).Err()
	}

	return nil
}

// CountByUserIDAndType 统计数量（直接查库即可）
func (a *addressRepo) CountByUserIDAndType(ctx context.Context, userID int64, addrType enums.AddressType) (int64, error) {
	var count int64
	err := a.db.WithContext(ctx).
		Model(&model.Address{}).
		Where("user_id = ? AND type = ?", userID, addrType).
		Count(&count).Error
	return count, err
}
