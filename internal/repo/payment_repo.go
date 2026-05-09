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
	UpdateByOrderID(ctx context.Context, orderID int64, status enums.PayStatus, At time.Time) error
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

func (p *paymentRepo) UpdateByOrderID(ctx context.Context, orderID int64, status enums.PayStatus, at time.Time) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 根据状态选择更新时间字段
	switch status {
	case enums.PayStatusPaid:
		updates["paid_at"] = at
	case enums.PayStatusRefunded:
		updates["refunded_at"] = at
	default:
		// 如果只是更新状态，不更新时间，可省略
	}

	res := p.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("order_id = ? AND status = ?", orderID, enums.PayStatusPending).
		Updates(updates)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errs.ErrPaymentNotFound
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrPaymentStatusInvalid
	}

	return nil
}
