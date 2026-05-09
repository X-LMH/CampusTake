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
	Create(ctx context.Context, order *model.Order) error
	GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status enums.OrderStatus) error
	GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error)
	UpdatePaidAt(ctx context.Context, orderID int64, paidAt time.Time) error
	GetAvailableForRider(ctx context.Context, page, size int, sortBy, order string, minReward, maxReward float64) (*response.PageResult, error)
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

func (o *orderRepo) UpdateStatus(ctx context.Context, orderID int64, status enums.OrderStatus) error {
	return o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
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

func (o *orderRepo) UpdatePaidAt(ctx context.Context, orderID int64, paidAt time.Time) error {
	return o.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("paid_at", paidAt).Error
}

func (o *orderRepo) GetAvailableForRider(ctx context.Context, page, size int, sortBy, order string, minReward, maxReward float64) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	query := o.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", enums.OrderPendingGrab)

	if minReward > 0 {
		query = query.Where("reward >= ?", minReward)
	}
	if maxReward > 0 {
		query = query.Where("reward <= ?", maxReward)
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

func (o *orderRepo) CreateLog(ctx context.Context, log *model.OrderLog) error {
	return o.db.WithContext(ctx).Create(log).Error
}
