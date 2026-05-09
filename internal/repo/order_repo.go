package repo

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	errs "CampusTake/pkg/errors"
	"context"
	"errors"

	"gorm.io/gorm"
)

type OrderRepo interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status enums.OrderStatus) error
	CreateLog(ctx context.Context, log *model.OrderLog) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{db: db}
}

// Create 创建订单
func (o *orderRepo) Create(ctx context.Context, order *model.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}

// GetByID 根据 ID 查询订单
func (o *orderRepo) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	order := new(model.Order)
	err := o.db.WithContext(ctx).First(order, "id = ?", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
}

// UpdateStatus 更新订单状态（Logic 层直接调用，不带旧状态参数，状态流转由状态机校验）
func (o *orderRepo) UpdateStatus(ctx context.Context, orderID int64, status enums.OrderStatus) error {
	res := o.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", orderID).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errs.ErrOrderStatusInvalid
	}
	return nil
}

// 创建订单日志
func (o *orderRepo) CreateLog(ctx context.Context, log *model.OrderLog) error {
	return o.db.WithContext(ctx).Create(log).Error
}
