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
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepo interface {
	// Order:
	Create(ctx context.Context, order *model.Order) error
	AllowCreateOrderLimit(ctx context.Context, userID int64) (bool, error)
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

	CreateLog(ctx context.Context, log *model.OrderLog) error

	UpdateAppealStatusByID(
		ctx context.Context,
		id int64,
		fromStatus enums.OrderAppealStatus,
		toStatus enums.OrderAppealStatus,
	) error

	CreateBloom(ctx context.Context) error
	InitBloom(ctx context.Context) error
	AddToBloom(ctx context.Context, orderID int64)
}
type orderRepo struct {
	RepoBase
	sf singleflight.Group
}

type orderListCache struct {
	Total int64          `json:"total"`
	List  []*model.Order `json:"list"`
}

const claimGrabOrderScript = `
return redis.call(
    "SET",
    KEYS[1],
    ARGV[1],
    "NX",
    "EX",
    ARGV[2]
) and 1 or 0
`

const releaseGrabOrderClaimScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
end
return 0
`

func NewOrderRepo(db *gorm.DB, rdb *redis.Client) OrderRepo {
	return &orderRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func (o *orderRepo) getDetailCacheKey(orderID int64) string {
	return fmt.Sprintf(constants.RedisKeyPrefixOrderDetail, orderID)
}

func (o *orderRepo) getGrabClaimKey(orderID int64) string {
	return fmt.Sprintf(constants.RedisKeyPrefixOrderGrabClaim, orderID)
}

func (o *orderRepo) getOrderDetailCacheTTL() time.Duration {
	base := constants.OrderDetailCacheTTL

	jitter := base / 10

	return base + time.Duration(rand.Int63n(int64(jitter)))
}

func (o *orderRepo) Create(ctx context.Context, order *model.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

func (o *orderRepo) AllowCreateOrderLimit(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprintf("rl:create_order:%d", userID)

	cnt, err := o.rdb.Incr(ctx, key).Result()
	if err != nil {
		return true, nil // Redis异常放行
	}

	if cnt == 1 {
		o.rdb.Expire(ctx, key, time.Second)
	}

	if cnt > 3 {
		return false, nil
	}

	return true, nil
}

func (o *orderRepo) GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error) {
	order, err := o.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, errs.ErrOrderNotFound
	}
	return order, nil
}

func (o *orderRepo) GetByIDAndRiderID(ctx context.Context, orderID int64, riderID int64) (*model.Order, error) {
	order, err := o.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.RiderID != nil && *order.RiderID != riderID {
		return nil, errs.ErrOrderNotFound
	}
	return order, nil
}

func (o *orderRepo) getOrderFromCache(
	ctx context.Context,
	cacheKey string,
) (*model.Order, bool, error) {

	val, err := o.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, false, nil
	}

	// 空值缓存
	if val == constants.EmptyOrderCacheValue {
		return nil, true, errs.ErrOrderNotFound
	}

	order := new(model.Order)
	if err := json.Unmarshal([]byte(val), order); err != nil {
		_ = o.rdb.Del(ctx, cacheKey).Err()
		return nil, false, nil
	}

	return order, true, nil
}
func (o *orderRepo) queryOrderFromDB(
	ctx context.Context,
	orderID int64,
) (*model.Order, error) {

	order := new(model.Order)

	err := o.db.WithContext(ctx).
		Where("id = ?", orderID).
		First(order).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}
func (o *orderRepo) setOrderCache(
	ctx context.Context,
	cacheKey string,
	order *model.Order,
) {

	data, err := json.Marshal(order)
	if err != nil {
		return
	}

	_ = o.rdb.Set(
		ctx,
		cacheKey,
		data,
		o.getOrderDetailCacheTTL(),
	).Err()
}

func (o *orderRepo) setEmptyOrderCache(
	ctx context.Context,
	cacheKey string,
) {

	_ = o.rdb.Set(
		ctx,
		cacheKey,
		constants.EmptyOrderCacheValue,
		constants.OrderEmptyCacheTTL,
	).Err()
}

func (o *orderRepo) ExistInBloom(
	ctx context.Context,
	orderID int64,
) bool {

	res, err := o.rdb.Do(
		ctx,
		"BF.EXISTS",
		constants.OrderBloomKey,
		orderID,
	).Int()

	if err != nil {
		// Bloom异常直接放行
		return true
	}

	return res == 1
}
func (o *orderRepo) CreateBloom(ctx context.Context) error {

	_, err := o.rdb.Do(
		ctx,
		"BF.RESERVE",
		constants.OrderBloomKey,
		0.001,
		1000000,
	).Result()

	if err != nil {
		// 已存在直接忽略
		if strings.Contains(err.Error(), "item exists") {
			return nil
		}
		return err
	}

	return nil
}

func (o *orderRepo) InitBloom(
	ctx context.Context,
) error {

	var ids []int64

	if err := o.db.
		Model(&model.Order{}).
		Pluck("id", &ids).
		Error; err != nil {
		return err
	}

	pipe := o.rdb.Pipeline()

	for _, id := range ids {
		pipe.Do(
			ctx,
			"BF.ADD",
			constants.OrderBloomKey,
			id,
		)
	}

	_, err := pipe.Exec(ctx)

	return err
}
func (o *orderRepo) AddToBloom(
	ctx context.Context,
	orderID int64,
) {

	_, _ = o.rdb.Do(
		ctx,
		"BF.ADD",
		constants.OrderBloomKey,
		orderID,
	).Result()
}

func (o *orderRepo) GetByID(
	ctx context.Context,
	orderID int64,
) (*model.Order, error) {

	// 参数校验
	if orderID <= 0 {
		return nil, errs.ErrOrderNotFound
	}

	// Bloom Filter
	if !o.ExistInBloom(ctx, orderID) {
		return nil, errs.ErrOrderNotFound
	}

	cacheKey := o.getDetailCacheKey(orderID)

	// 第一次查缓存
	if order, hit, err := o.getOrderFromCache(ctx, cacheKey); hit {
		return order, err
	}

	v, err, _ := o.sf.Do(
		cacheKey,
		func() (interface{}, error) {

			// Double Check
			if order, hit, err := o.getOrderFromCache(ctx, cacheKey); hit {
				return order, err
			}

			order, err := o.queryOrderFromDB(ctx, orderID)

			if err != nil {

				if errors.Is(err, errs.ErrOrderNotFound) {

					o.setEmptyOrderCache(ctx, cacheKey)
				}

				return nil, err
			}

			o.setOrderCache(ctx, cacheKey, order)

			return order, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return v.(*model.Order), nil
}

func (o *orderRepo) GetByIDForUpdate(ctx context.Context, orderID int64) (*model.Order, error) {
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
		list  []*model.Order
		total int64
	)
	cacheKey := fmt.Sprintf(constants.RedisKeyPrefixUserOrderList, userID, status, page, size)
	if cached, ok := o.getOrderListCache(ctx, cacheKey); ok {
		return response.NewPageResult(cached.Total, cached.List), nil
	}

	v, err, _ := o.sf.Do(cacheKey, func() (interface{}, error) {

		// 双重检查
		if cached, ok := o.getOrderListCache(ctx, cacheKey); ok {
			return cached, nil
		}

		buildQuery := func() *gorm.DB {
			orderQuery := o.db.WithContext(ctx).
				Model(&model.Order{}).
				Where("user_id = ?", userID)

			if status != enums.OrderDefault {
				orderQuery = orderQuery.Where("status = ?", status)
			}

			return orderQuery
		}

		countQuery := buildQuery()
		if err := countQuery.Count(&total).Error; err != nil {
			return nil, err
		}

		if total == 0 {
			result := &orderListCache{
				Total: total,
				List:  make([]*model.Order, 0),
			}

			o.setOrderListCache(ctx, cacheKey, result)
			return result, nil
		}

		listQuery := buildQuery()

		if err := listQuery.
			Scopes(db.Paginate(page, size)).
			Order("created_at DESC").
			Find(&list).Error; err != nil {
			return nil, err
		}

		result := &orderListCache{
			Total: total,
			List:  list,
		}

		o.setOrderListCache(ctx, cacheKey, result)

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	result := v.(*orderListCache)

	return response.NewPageResult(result.Total, result.List), nil
}

func (o *orderRepo) GetAvailableForRider(
	ctx context.Context,
	page, size int,
	sortBy, order string,
	minReward, maxReward float64,
) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	cacheKey := fmt.Sprintf(
		constants.RedisKeyPrefixAvailableOrderList,
		page,
		size,
		sortBy,
		order,
		strconv.FormatFloat(minReward, 'f', 2, 64),
		strconv.FormatFloat(maxReward, 'f', 2, 64),
	)
	if cached, ok := o.getOrderListCache(ctx, cacheKey); ok {
		return response.NewPageResult(cached.Total, cached.List), nil
	}

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

	orderClause := "created_at DESC"
	order = strings.ToUpper(order)
	switch sortBy {
	case "created_at", "reward_amount", "id":
		if order == "ASC" || order == "DESC" {
			orderClause = sortBy + " " + order
		}
	}

	err := orderQuery.Scopes(db.Paginate(page, size)).
		Order(orderClause).
		Find(&list).Error
	if err == nil {
		o.setOrderListCache(ctx, cacheKey, &orderListCache{Total: total, List: list})
	}

	return response.NewPageResult(total, list), err
}

func (o *orderRepo) getOrderListCache(ctx context.Context, key string) (*orderListCache, bool) {
	val, err := o.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	cached := new(orderListCache)
	if err := json.Unmarshal([]byte(val), cached); err != nil {
		_ = o.rdb.Del(ctx, key).Err()
		return nil, false
	}

	if cached.List == nil {
		cached.List = make([]*model.Order, 0)
	}
	return cached, true
}

func (o *orderRepo) setOrderListCache(ctx context.Context, key string, value *orderListCache) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	ttl := constants.OrderListCacheTTL
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	jitter := ttl / 10
	if jitter <= 0 {
		jitter = 30 * time.Second
	}
	_ = o.rdb.Set(ctx, key, data, ttl+time.Duration(rand.Int63n(int64(jitter)))).Err()
}

func (o *orderRepo) TryClaimGrabOrder(
	ctx context.Context,
	orderID int64,
	riderID int64,
) (bool, error) {
	key := o.getGrabClaimKey(orderID)

	result, err := o.rdb.Eval(
		ctx,
		claimGrabOrderScript,
		[]string{key},
		riderID,
		int(constants.OrderGrabClaimTTL.Seconds()),
	).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (o *orderRepo) ReleaseGrabOrderClaim(
	ctx context.Context,
	orderID int64,
	riderID int64,
) error {
	return o.rdb.Eval(
		ctx,
		releaseGrabOrderClaimScript,
		[]string{o.getGrabClaimKey(orderID)},
		riderID,
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
		Where("id = ? and rider_id is null", orderID).
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

	_ = o.rdb.Del(ctx, o.getDetailCacheKey(orderID)).Err()

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

	err := orderQuery.Scopes(db.Paginate(page, size)).
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return response.NewPageResult(total, list), nil
}

func (o *orderRepo) UpdatePaymentStatusByIDAndUserID(ctx context.Context, orderID int64, userID int64, status enums.OrderPaymentStatus) error {
	res := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? and user_id = ?", orderID, userID).
		Update("payment_status", status)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected > 0 {
		_ = o.rdb.Del(ctx, o.getDetailCacheKey(orderID)).Err()
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

	_ = o.rdb.Del(ctx, o.getDetailCacheKey(q.OrderID)).Err()
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

	_ = o.rdb.Del(ctx, o.getDetailCacheKey(id)).Err()
	return nil
}
