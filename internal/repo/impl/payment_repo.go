package impl

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
	UpdateStatusByOrderID(ctx context.Context, orderID int64, fromStatus enums.OrderPaymentStatus, toStatus enums.OrderPaymentStatus, at time.Time) error
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

func (p *paymentRepo) UpdateStatusByOrderID(
	ctx context.Context,
	orderID int64,
	fromStatus enums.OrderPaymentStatus,
	toStatus enums.OrderPaymentStatus,
	at time.Time,
) error {

	updates := map[string]interface{}{
		"status": toStatus,
	}

	switch toStatus {
	case enums.OrderPayStatusPaid:
		updates["paid_at"] = at

	case enums.OrderPayStatusRefunded:
		updates["refunded_at"] = at
	}

	res := p.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where(
			"order_id = ? AND status = ?",
			orderID,
			fromStatus,
		).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrPaymentStatusInvalid
	}

	return nil
}
