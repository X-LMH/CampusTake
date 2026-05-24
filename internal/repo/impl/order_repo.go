package impl

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo/query"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"context"
	"errors"
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

	GrabOrder(ctx context.Context, orderID int64, riderID int64, acceptedAt time.Time) error
	GetListByRiderID(ctx context.Context, riderID int64, status enums.OrderStatus, page, size int) (*response.PageResult, error)
	UpdatePaymentStatusByIDAndUserID(ctx context.Context, orderID int64, userID int64, status enums.OrderPaymentStatus) error
	GetCountByRiderIDAndStatuses(ctx context.Context, riderID int64, statuses []enums.OrderStatus) (int64, error)

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

func NewOrderRepo(db *gorm.DB, rdb redis.Cmdable) OrderRepo {
	return &orderRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func (o *orderRepo) Create(ctx context.Context, order *model.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

func (o *orderRepo) GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error) {
	order := new(model.Order)
	err := o.db.WithContext(ctx).Where("id = ? AND user_id = ?", orderID, userID).First(order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
}
func (o *orderRepo) GetByIDAndRiderID(ctx context.Context, orderID int64, riderID int64) (*model.Order, error) {
	order := new(model.Order)
	err := o.db.WithContext(ctx).Where("id = ? AND rider_id = ?", orderID, riderID).First(order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
}

func (o *orderRepo) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	order := new(model.Order)
	err := o.db.WithContext(ctx).Where("id = ?", orderID).First(order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
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
		list  []model.Order // 核心修改 1：改成结构体值切片，提升内存连续性，减轻 GC 压力
		total int64
	)

	// 核心修改 2：采用你之前在 Appeal 里的闭包优秀实践，彻底隔离 Count 和 Find 的状态污染
	buildQuery := func() *gorm.DB {
		orderQuery := o.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

		// 状态过滤（只要不是默认状态，就作为有效条件）
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

	// 性能小优化：如果总数本来就是 0，直接打道回府，没必要去空查一遍数据库
	if total == 0 {
		return response.NewPageResult(total, make([]model.Order, 0)), nil
	}

	// 2. 再查具体分页列表
	listQuery := buildQuery()
	err := listQuery.
		Scopes(db.Paginate(page, size)).
		Order("created_at DESC"). // 校园取送高频场景，最新订单必须在最上面
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

	return nil
}

func (o *orderRepo) GetListByRiderID(ctx context.Context, riderID int64, status enums.OrderStatus, page, size int) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	orderQuery := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("rider_id = ?", riderID)

	// 筛查
	if status != enums.OrderDefault {
		orderQuery = orderQuery.Where("status = ?", status)
	}

	// 统计总数
	if err := orderQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	orderQuery.Scopes(db.Paginate(page, size)).
		Order("created_at DESC").
		Find(&list)

	return response.NewPageResult(total, list), nil
}

func (o *orderRepo) UpdatePaymentStatusByIDAndUserID(ctx context.Context, orderID int64, userID int64, status enums.OrderPaymentStatus) error {
	return o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ? and user_id = ?",
			orderID,
			userID,
		).
		Update("payment_status", status).Error
}

func (o *orderRepo) GetCountByRiderIDAndStatuses(ctx context.Context, riderID int64, statuses []enums.OrderStatus) (int64, error) {
	var count int64
	err := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("rider_id = ? AND status IN ?", riderID, statuses).
		Count(&count).Error
	return count, err
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

	// =========================
	// 校验状态流转
	// =========================

	if !enums.CheckOrderStatusFlow(fromStatus, toStatus) {
		return errs.ErrOrderStatusInvalid
	}

	// =========================
	// WHERE 条件
	// =========================

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

	// =========================
	// 更新字段
	// =========================

	updates := map[string]interface{}{
		"status": toStatus,
	}

	// 自动状态时间
	if field := enums.GetOrderStatusTimeField(toStatus); field != "" {
		updates[field] = at
	}

	// 额外字段
	for k, v := range extraUpdates {
		updates[k] = v
	}

	// =========================
	// 更新数据库
	// =========================

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

	return nil
}
