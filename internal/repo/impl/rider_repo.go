package impl

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RiderRepo interface {
	GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error)

	GetProfileByRiderID(ctx context.Context, riderID int64) (*model.RiderProfile, error)

	UpsertProfile(ctx context.Context, profile *model.RiderProfile) error

	UpdateStatusByUserID(
		ctx context.Context,
		userID int64,
		status enums.RiderAuditStatus,
	) error

	UpdateStatusAndRemarkByUserID(
		ctx context.Context,
		userID int64,
		status enums.RiderAuditStatus,
		remark string,
	) error

	GetProfileList(
		ctx context.Context,
		status enums.RiderAuditStatus,
		page,
		pageSize int,
	) (*response.PageResult, error)

	UpdateRatingByRiderID(
		ctx context.Context,
		riderID int64,
		newCount int64,
		newAvg float64,
	) error

	IncrementAcceptedOrderCount(ctx context.Context, riderID int64) error

	IncrementCompletedOrderCount(ctx context.Context, riderID int64) error

	DecrementCompletedOrderCount(ctx context.Context, riderID int64) error

	CreateLog(ctx context.Context, log *model.RiderAuditLog) error
}

type riderRepo struct {
	RepoBase
}

func NewRiderRepo(db *gorm.DB, rdb redis.Cmdable) RiderRepo {
	return &riderRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

// =====================================
// 获取骑手资料（user_id）
// =====================================

func (r *riderRepo) GetProfileByUserID(
	ctx context.Context,
	userID int64,
) (*model.RiderProfile, error) {

	key := fmt.Sprintf(
		constants.RedisKeyPrefixRiderProfileUserID,
		userID,
	)

	// =========================
	// 1. 查询 Redis
	// =========================

	val, err := r.rdb.Get(ctx, key).Result()

	if err == nil {

		profile := new(model.RiderProfile)

		if json.Unmarshal([]byte(val), profile) == nil {
			return profile, nil
		}

		// 删除脏缓存
		_ = r.rdb.Del(ctx, key).Err()
	}

	// redis 异常降级
	if err != nil && !errors.Is(err, redis.Nil) {
		// ignore
	}

	// =========================
	// 2. 查询 MySQL
	// =========================

	profile := new(model.RiderProfile)

	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(profile).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	// =========================
	// 3. 回填缓存
	// =========================

	r.setProfileCache(ctx, profile)

	return profile, nil
}

// =====================================
// 获取骑手资料（rider_id）
// =====================================

func (r *riderRepo) GetProfileByRiderID(
	ctx context.Context,
	riderID int64,
) (*model.RiderProfile, error) {

	key := fmt.Sprintf(
		constants.RedisKeyPrefixRiderProfileID,
		riderID,
	)

	// =========================
	// 1. 查询 Redis
	// =========================

	val, err := r.rdb.Get(ctx, key).Result()

	if err == nil {

		profile := new(model.RiderProfile)

		if json.Unmarshal([]byte(val), profile) == nil {
			return profile, nil
		}

		// 删除脏缓存
		_ = r.rdb.Del(ctx, key).Err()
	}

	// redis 异常降级
	if err != nil && !errors.Is(err, redis.Nil) {
		// ignore
	}

	// =========================
	// 2. 查询 MySQL
	// =========================

	profile := new(model.RiderProfile)

	err = r.db.WithContext(ctx).
		Where("id = ?", riderID).
		First(profile).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	// =========================
	// 3. 回填缓存
	// =========================

	r.setProfileCache(ctx, profile)

	return profile, nil
}

// =====================================
// 新增/更新骑手资料
// =====================================

func (r *riderRepo) UpsertProfile(
	ctx context.Context,
	profile *model.RiderProfile,
) error {

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"real_name",
				"student_no",
				"id_card_no",
				"dormitory_building",
				"dormitory_room",
				"campus_card_front",
				"campus_card_back",
				"audit_status",
				"updated_at",
			}),
		}).
		Create(profile).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 更新审核状态
// =====================================

func (r *riderRepo) UpdateStatusByUserID(
	ctx context.Context,
	userID int64,
	status enums.RiderAuditStatus,
) error {

	profile, err := r.loadProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", userID).
		Update("audit_status", status).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 更新审核状态+备注
// =====================================

func (r *riderRepo) UpdateStatusAndRemarkByUserID(
	ctx context.Context,
	userID int64,
	status enums.RiderAuditStatus,
	remark string,
) error {

	profile, err := r.loadProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"audit_status": status,
			"audit_remark": remark,
		}).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 获取骑手列表
// =====================================

func (r *riderRepo) GetProfileList(
	ctx context.Context,
	status enums.RiderAuditStatus,
	page,
	pageSize int,
) (*response.PageResult, error) {

	var list []*model.RiderProfile
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.RiderProfile{})

	if status > 0 {
		query = query.Where("audit_status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	err := query.
		Scopes(db.Paginate(page, pageSize)).
		Order("updated_at DESC").
		Find(&list).Error

	return response.NewPageResult(total, list), err
}

// =====================================
// 更新评分
// =====================================

func (r *riderRepo) UpdateRatingByRiderID(
	ctx context.Context,
	riderID int64,
	newCount int64,
	newAvg float64,
) error {

	profile, err := r.loadProfileByRiderID(ctx, riderID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("id = ?", riderID).
		Updates(map[string]interface{}{
			"rating_count": newCount,
			"rating_avg":   newAvg,
		}).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 已接单数 +1
// =====================================

func (r *riderRepo) IncrementAcceptedOrderCount(
	ctx context.Context,
	riderID int64,
) error {

	profile, err := r.loadProfileByRiderID(ctx, riderID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("id = ?", riderID).
		Update(
			"accepted_order_count",
			gorm.Expr("accepted_order_count + ?", 1),
		).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 已完成订单数 +1
// =====================================

func (r *riderRepo) IncrementCompletedOrderCount(
	ctx context.Context,
	riderID int64,
) error {

	profile, err := r.loadProfileByRiderID(ctx, riderID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("id = ?", riderID).
		Update(
			"completed_order_count",
			gorm.Expr("completed_order_count + ?", 1),
		).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 已完成订单数 -1
// =====================================

func (r *riderRepo) DecrementCompletedOrderCount(
	ctx context.Context,
	riderID int64,
) error {

	profile, err := r.loadProfileByRiderID(ctx, riderID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("id = ? AND completed_order_count > 0", riderID).
		Update(
			"completed_order_count",
			gorm.Expr("completed_order_count - ?", 1),
		).Error

	if err != nil {
		return err
	}

	r.deleteProfileCache(ctx, profile)

	return nil
}

// =====================================
// 创建审核日志
// =====================================

func (r *riderRepo) CreateLog(
	ctx context.Context,
	log *model.RiderAuditLog,
) error {

	return r.db.WithContext(ctx).
		Create(log).Error
}

// =====================================
// 设置缓存
// =====================================

func (r *riderRepo) setProfileCache(
	ctx context.Context,
	profile *model.RiderProfile,
) {

	bytes, err := json.Marshal(profile)

	if err != nil {
		return
	}

	userKey := fmt.Sprintf(
		constants.RedisKeyPrefixRiderProfileUserID,
		profile.UserID,
	)

	idKey := fmt.Sprintf(
		constants.RedisKeyPrefixRiderProfileID,
		profile.ID,
	)

	pipe := r.rdb.Pipeline()

	pipe.Set(
		ctx,
		userKey,
		bytes,
		constants.RiderProfileCacheTTL,
	)

	pipe.Set(
		ctx,
		idKey,
		bytes,
		constants.RiderProfileCacheTTL,
	)

	_, _ = pipe.Exec(ctx)
}

// =====================================
// 删除缓存
// =====================================

func (r *riderRepo) deleteProfileCache(
	ctx context.Context,
	profile *model.RiderProfile,
) {

	keys := []string{
		fmt.Sprintf(
			constants.RedisKeyPrefixRiderProfileUserID,
			profile.UserID,
		),

		fmt.Sprintf(
			constants.RedisKeyPrefixRiderProfileID,
			profile.ID,
		),
	}

	_ = r.rdb.Del(ctx, keys...).Err()
}

// =====================================
// 内部查询
// =====================================

func (r *riderRepo) loadProfileByUserID(
	ctx context.Context,
	userID int64,
) (*model.RiderProfile, error) {

	profile := new(model.RiderProfile)

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(profile).Error

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *riderRepo) loadProfileByRiderID(
	ctx context.Context,
	riderID int64,
) (*model.RiderProfile, error) {

	profile := new(model.RiderProfile)

	err := r.db.WithContext(ctx).
		Where("id = ?", riderID).
		First(profile).Error

	if err != nil {
		return nil, err
	}

	return profile, nil
}
