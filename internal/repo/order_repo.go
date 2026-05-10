package repo

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OrderRepo interface {
	// Order:
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
	GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error)
	GetByIDAndRiderID(ctx context.Context, orderID int64, riderID int64) (*model.Order, error)
	UserUpdateStatusAndTime(ctx context.Context, orderID int64, userID int64, fromStatus enums.OrderStatus, toStatus enums.OrderStatus, at time.Time) error
	RiderUpdateStatusAndTime(ctx context.Context, orderID int64, riderID int64, fromStatus enums.OrderStatus, toStatus enums.OrderStatus, at time.Time) error
	GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error)
	GetAvailableForRider(ctx context.Context, page, size int, sortBy, order string, minReward, maxReward float64) (*response.PageResult, error)
	GrabOrder(ctx context.Context, orderID int64, riderID int64, acceptedAt time.Time) error
	GetListByRiderID(ctx context.Context, riderID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error)
	// LOG:
	CreateLog(ctx context.Context, log *model.OrderLog) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{db: db}
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

func (o *orderRepo) UserUpdateStatusAndTime(
	ctx context.Context,
	orderID int64,
	userID int64,
	fromStatus enums.OrderStatus,
	toStatus enums.OrderStatus,
	at time.Time,
) error {

	return o.updateStatusAndTime(
		ctx,
		map[string]interface{}{
			"id":      orderID,
			"user_id": userID,
			"status":  fromStatus,
		},
		toStatus,
		at,
	)
}

func (o *orderRepo) RiderUpdateStatusAndTime(
	ctx context.Context,
	orderID int64,
	riderID int64,
	fromStatus enums.OrderStatus,
	toStatus enums.OrderStatus,
	at time.Time,
) error {

	return o.updateStatusAndTime(
		ctx,
		map[string]interface{}{
			"id":       orderID,
			"rider_id": riderID,
			"status":   fromStatus,
		},
		toStatus,
		at,
	)
}

func (o *orderRepo) updateStatusAndTime(
	ctx context.Context,
	where map[string]interface{},
	toStatus enums.OrderStatus,
	at time.Time,
) error {

	updates := map[string]interface{}{
		"status": toStatus,
	}

	// 自动更新状态对应时间字段
	if field := enums.GetOrderStatusTimeField(toStatus); field != "" {
		updates[field] = at
	}

	res := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where(where).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	// 乐观锁失败
	if res.RowsAffected == 0 {
		return errs.ErrOrderStatusInvalid
	}

	return nil
}
func (o *orderRepo) GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	query := o.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)
	if status != enums.OrderDefault {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	err := query.Scopes(db.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&list).Error
	return response.NewPageResult(total, list), err
}

func (o *orderRepo) GetAvailableForRider(ctx context.Context, page, size int, sortBy, order string, minReward, maxReward float64) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	query := o.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", enums.OrderPendingGrab)

	if minReward > 0 {
		query = query.Where("reward_amount >= ?", minReward)
	}
	if maxReward > 0 {
		query = query.Where("reward_amount <= ?", maxReward)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 拼接排序字符串，默认按创建时间倒序
	orderClause := "created_at DESC"
	if sortBy != "" && order != "" {
		orderClause = sortBy + " " + order
	}

	err := query.Scopes(db.Paginate(page, size)).
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
			"id = ? AND status = ? AND rider_id IS NULL",
			orderID,
			enums.OrderPendingGrab,
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

func (o *orderRepo) GetListByRiderID(ctx context.Context, riderID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	query := o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("rider_id = ?", riderID)

	// 筛查
	if status != enums.OrderDefault {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	query.Scopes(db.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&list)

	return response.NewPageResult(total, list), nil
}

func (o *orderRepo) CreateLog(ctx context.Context, log *model.OrderLog) error {
	return o.db.WithContext(ctx).Create(log).Error
}
