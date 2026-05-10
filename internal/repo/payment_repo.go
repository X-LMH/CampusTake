package repo

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	errs "CampusTake/pkg/errors"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByOrderID(ctx context.Context, orderID int64) (*model.Payment, error)
	UpdateByOrderID(ctx context.Context, orderID int64, status enums.OrderPaymentStatus, at time.Time) error
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) PaymentRepo {
	return &paymentRepo{db: db}
}

func (p *paymentRepo) Create(ctx context.Context, payment *model.Payment) error {
	return p.db.WithContext(ctx).Create(payment).Error
}

func (p *paymentRepo) GetByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	payment := new(model.Payment)
	err := p.db.WithContext(ctx).Where("order_id = ?", orderID).First(payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPaymentNotFound
		}
		return nil, err
	}
	return payment, nil
}

// UpdateByOrderID 带有状态保护的更新
func (p *paymentRepo) UpdateByOrderID(ctx context.Context, orderID int64, status enums.OrderPaymentStatus, at time.Time) error {
	updates := map[string]interface{}{
		"status": status,
	}

	switch status {
	case enums.OrderPayStatusPaid:
		updates["paid_at"] = at
	case enums.OrderPayStatusRefunded:
		updates["refunded_at"] = at
	}

	// 这里保留 Where("status = ?", enums.PayStatusPending) 作为乐观锁/状态保护
	res := p.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("order_id = ? AND status = ?", orderID, enums.OrderPayStatusUnpaid).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	// 支付场景建议保留此判断，因为“状态已改变”在支付逻辑中是关键异常
	if res.RowsAffected == 0 {
		return errs.ErrPaymentStatusInvalid
	}

	return nil
}
