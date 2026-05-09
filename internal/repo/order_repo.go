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
	GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status enums.OrderStatus) error
	GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error)
	UpdatePaidAt(ctx context.Context, orderID int64, paidAt time.Time) error
	// Log:
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
func (o *orderRepo) GetByIDAndUserID(ctx context.Context, orderID int64, userID int64) (*model.Order, error) {
	order := new(model.Order)
	err := o.db.WithContext(ctx).First(order, "id = ? and user_id = ?", orderID, userID).Error
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

func (o *orderRepo) GetListByUserID(ctx context.Context, userID int64, status enums.OrderStatus, page, pageSize int) (*response.PageResult, error) {
	var list []*model.Order
	var total int64

	query := o.db.WithContext(ctx).Model(&model.Order{})
	query = query.Where("user_id = ?", userID)

	// 添加过滤条件
	if status != enums.OrderDefault {
		query = query.Where("status = ?", status)
	}

	// 先查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	err := query.Scopes(db.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&list).Error
	return response.NewPageResult(total, list), err
}

func (o *orderRepo) UpdatePaidAt(ctx context.Context, orderID int64, paidAt time.Time) error {
	res := o.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", orderID).Update("paid_at", paidAt)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errs.ErrOrderNotFound
		}
	}
	if res.RowsAffected == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}

func (o *orderRepo) CreateLog(ctx context.Context, log *model.OrderLog) error {
	return o.db.WithContext(ctx).Create(log).Error
}
