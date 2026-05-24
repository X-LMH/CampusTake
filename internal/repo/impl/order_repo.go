package impl

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo/query"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepo interface {
	// Order:
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
	GetByIDForUpdate(ctx context.Context, orderID int64) (*model.Order, error)
	GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error)
	GetByIDAndRiderID(ctx context.Context, orderID int64, riderID int64) (*model.Order, error)

	UpdateStatusAndTime(
		ctx context.Context,
		q query.OrderStatusUpdateQuery,
		fromStatus enums.OrderStatus,
		toStatus enums.OrderStatus,
		at time.Time,
		extraUpdates map[string]interface{},
	) error

	GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, size int) (*response.PageResult, error)

	GetAvailableForRider(
		ctx context.Context,
		page, size int,
		sortBy, order string,
		minReward, maxReward float64,
	) (*response.PageResult, error)

	TryClaimGrabOrder(ctx context.Context, orderID int64, riderID int64) (bool, error)
	ReleaseGrabOrderClaim(ctx context.Context, orderID int64, riderID int64) error
	GrabOrder(ctx context.Context, orderID int64, riderID int64, acceptedAt time.Time) error
	GetListByRiderID(ctx context.Context, riderID int64, status enums.OrderStatus, page, size int) (*response.PageResult, error)
	UpdatePaymentStatusByIDAndUserID(ctx context.Context, orderID int64, userID int64, status enums.OrderPaymentStatus) error
	//GetCountByRiderIDAndStatuses(ctx context.Context, riderID int64, statuses []enums.OrderStatus) (int64, error)

	CreateLog(ctx context.Context, log *model.OrderLog) error

	UpdateAppealStatusByID(
		ctx context.Context,
		id int64,
		fromStatus enums.OrderAppealStatus,
		toStatus enums.OrderAppealStatus,
	) error
}

type orderRepo struct {
	RepoBase
}

const claimGrabOrderScript = `
local current = redis.call("GET", KEYS[1])
if current then
	if current == ARGV[1] then
		return 2
	end
	return 0
end
redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
return 1
`

const releaseGrabOrderClaimScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

func NewOrderRepo(db *gorm.DB, rdb redis.Cmdable) OrderRepo {
	return &orderRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

// 🔑 统一管理订单详情的 Redis Key 规则
func (o *orderRepo) GetDetailCacheKey(orderID int64) string {
	return fmt.Sprintf(constants.RedisKeyPrefixOrderDetail, orderID)
}

func (o *orderRepo) getGrabClaimKey(orderID int64) string {
	return fmt.Sprintf(constants.RedisKeyPrefixOrderGrabClaim, orderID)
}

func (o *orderRepo) Create(ctx context.Context, order *model.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

func (o *orderRepo) GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error) {
	// 1. 直接调用 GetByID，优先享受 Redis 缓存红利
	order, err := o.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 2. 在内存里进行越权校验
	if order.UserID != userID {
		return nil, errs.ErrOrderNotFound
	}

	return order, nil
}

func (o *orderRepo) GetByIDAndRiderID(ctx context.Context, orderID int64, riderID int64) (*model.Order, error) {
	// 1. 同样复用 GetByID 缓存
	order, err := o.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 2. 内存校验骑手身份
	if order.RiderID != nil && *order.RiderID != riderID {
		return nil, errs.ErrOrderNotFound
	}

	return order, nil
}

func (o *orderRepo) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	cacheKey := o.GetDetailCacheKey(orderID)

	// 1. 先去 Redis 捞一网
	val, err := o.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 🚀 命中缓存！直接反序列化返回，不用惊动 MySQL
		order := new(model.Order)
		if json.Unmarshal([]byte(val), order) == nil {
			return order, nil
		}
	}

	// 2. 缓存没命中（或者 Redis 挂了），去查 MySQL
	order := new(model.Order)
	err = o.db.WithContext(ctx).Where("id = ?", orderID).First(order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}

	// 3. 查出来了，顺手把最新的完整 model 塞回 Redis
	if data, err := json.Marshal(order); err == nil {
		_ = o.rdb.Set(ctx, cacheKey, data, constants.OrderDetailCacheTTL).Err()
	}

	return order, nil
}

func (o *orderRepo) GetByIDForUpdate(ctx context.Context, orderID int64) (*model.Order, error) {
	// ⚠️ 注意：带有排他锁 (FOR UPDATE) 的查询绝不能走 Redis 缓存，必须强一致性查库
	order := new(model.Order)
	err := o.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", orderID).
		First(order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
}

func (o *orderRepo) GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, size int) (*response.PageResult, error) {
	var (
		list  []model.Order
		total int64
	)

	buildQuery := func() *gorm.DB {
		orderQuery := o.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

		if status != enums.OrderDefault {
			orderQuery = orderQuery.Where("status = ?", status)
		}
		return orderQuery
	}

	// 1. 先查总数
	countQuery := buildQuery()
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	if total == 0 {
		return response.NewPageResult(total, make([]model.Order, 0)), nil
	}

	// 2. 再查具体分页列表
	listQuery := buildQuery()
	err := listQuery.
		Scopes(db.Paginate(page, size)).
		Order("created_at DESC").
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	return response.NewPageResult(total, list), nil
}

func (o *orderRepo) GetAvailableForRider(
	ctx context.Context,
	page, size int,
	sortBy, order string,
	minReward, maxReward float64,
) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	orderQuery := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("rider_id IS NULL").
		Where("(status = ? or (status = ? and can_reassign = ?))",
			enums.OrderPendingGrab,
			enums.OrderRiderCancelled,
			enums.OrderCanReassignYes,
		)

	if minReward > 0 {
		orderQuery = orderQuery.Where("reward_amount >= ?", minReward)
	}
	if maxReward > 0 {
		orderQuery = orderQuery.Where("reward_amount <= ?", maxReward)
	}

	if err := orderQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// 白名单排序
	orderClause := "created_at DESC"
	switch sortBy {
	case "created_at", "reward_amount", "id":
		if order == "ASC" || order == "DESC" {
			orderClause = sortBy + " " + order
		}
	}

	err := orderQuery.Scopes(db.Paginate(page, size)).
		Order(orderClause).
		Find(&list).Error

	return response.NewPageResult(total, list), err
}

func (o *orderRepo) TryClaimGrabOrder(ctx context.Context, orderID int64, riderID int64) (bool, error) {
	result, err := o.rdb.Eval(
		ctx,
		claimGrabOrderScript,
		[]string{o.getGrabClaimKey(orderID)},
		strconv.FormatInt(riderID, 10),
		int(constants.OrderGrabClaimTTL.Seconds()),
	).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (o *orderRepo) ReleaseGrabOrderClaim(ctx context.Context, orderID int64, riderID int64) error {
	return o.rdb.Eval(
		ctx,
		releaseGrabOrderClaimScript,
		[]string{o.getGrabClaimKey(orderID)},
		strconv.FormatInt(riderID, 10),
	).Err()
}

func (o *orderRepo) GrabOrder(
	ctx context.Context,
	orderID int64,
	riderID int64,
	acceptedAt time.Time,
) error {
	res := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where(
			"id = ? and rider_id is null",
			orderID,
		).
		Where("(status = ? or (status = ? and can_reassign = ?))",
			enums.OrderPendingGrab,
			enums.OrderRiderCancelled,
			enums.OrderCanReassignYes,
		).
		Updates(map[string]interface{}{
			"rider_id":    riderID,
			"accepted_at": acceptedAt,
			"status":      enums.OrderAccepted,
		})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrOrderHaveGrabbed
	}

	// 🚀 骑手抢单成功后，利用统一方法生成 Key 并清除缓存
	cacheKey := o.GetDetailCacheKey(orderID)
	_ = o.rdb.Del(ctx, cacheKey).Err()

	return nil
}

func (o *orderRepo) GetListByRiderID(ctx context.Context, riderID int64, status enums.OrderStatus, page, size int) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	orderQuery := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("rider_id = ?", riderID)

	if status != enums.OrderDefault {
		orderQuery = orderQuery.Where("status = ?", status)
	}

	if err := orderQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	orderQuery.Scopes(db.Paginate(page, size)).
		Order("created_at DESC").
		Find(&list)

	return response.NewPageResult(total, list), nil
}

func (o *orderRepo) UpdatePaymentStatusByIDAndUserID(ctx context.Context, orderID int64, userID int64, status enums.OrderPaymentStatus) error {
	res := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? and user_id = ?",
			orderID,
			userID,
		).
		Update("payment_status", status)

	if res.Error != nil {
		return res.Error
	}

	// 🚀 支付状态发生改变，在更新成功时同步清除缓存
	if res.RowsAffected > 0 {
		cacheKey := o.GetDetailCacheKey(orderID)
		_ = o.rdb.Del(ctx, cacheKey).Err()
	}

	return nil
}

func (o *orderRepo) CreateLog(ctx context.Context, log *model.OrderLog) error {
	return o.db.WithContext(ctx).Create(log).Error
}

func (o *orderRepo) UpdateStatusAndTime(
	ctx context.Context,
	q query.OrderStatusUpdateQuery,
	fromStatus enums.OrderStatus,
	toStatus enums.OrderStatus,
	at time.Time,
	extraUpdates map[string]interface{},
) error {

	if !enums.CheckOrderStatusFlow(fromStatus, toStatus) {
		return errs.ErrOrderStatusInvalid
	}

	where := map[string]interface{}{
		"id":     q.OrderID,
		"status": fromStatus,
	}

	if q.UserID != nil {
		where["user_id"] = *q.UserID
	}

	if q.RiderID != nil {
		where["rider_id"] = *q.RiderID
	}

	updates := map[string]interface{}{
		"status": toStatus,
	}

	if field := enums.GetOrderStatusTimeField(toStatus); field != "" {
		updates[field] = at
	}

	for k, v := range extraUpdates {
		updates[k] = v
	}

	res := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where(where).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrOrderStatusInvalid
	}

	// ⭐ 状态修改成功，同步利用统一方法清除详情缓存
	cacheKey := o.GetDetailCacheKey(q.OrderID)
	_ = o.rdb.Del(ctx, cacheKey).Err()

	return nil
}

func (o *orderRepo) UpdateAppealStatusByID(
	ctx context.Context,
	id int64,
	fromStatus enums.OrderAppealStatus,
	toStatus enums.OrderAppealStatus,
) error {
	if !fromStatus.CanTransitionTo(toStatus) {
		return errs.ErrOrderStatusInvalid
	}

	result := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? AND appeal_status = ?", id, fromStatus).
		Update("appeal_status", toStatus)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrOrderAppealStatusChanged
	}

	// ⭐ 申诉状态流转成功后，同步清除详情缓存
	cacheKey := o.GetDetailCacheKey(id)
	_ = o.rdb.Del(ctx, cacheKey).Err()

	return nil
}
