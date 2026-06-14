package impl

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	errs "CampusTake/pkg/errors"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByOrderID(ctx context.Context, orderID int64) (*model.Payment, error)
	GetByOrderIDForUpdate(ctx context.Context, orderID int64) (*model.Payment, error)

	UpdateStatusByOrderID(ctx context.Context, orderID int64, fromStatus enums.PaymentStatus, toStatus enums.PaymentStatus, at time.Time) error

	UpdateRefundInfo(
		ctx context.Context,
		paymentID int64,
		refundAmount float64,
		status enums.PaymentStatus,
		at time.Time,
	) error
}

type paymentRepo struct {
	RepoBase
}

func NewPaymentRepo(db *gorm.DB, rdb *redis.Client) PaymentRepo {
	return &paymentRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
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

func (p *paymentRepo) GetByOrderIDForUpdate(ctx context.Context, orderID int64) (*model.Payment, error) {
	payment := new(model.Payment)
	err := p.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPaymentNotFound
		}
		return nil, err
	}
	return payment, nil
}

func (p *paymentRepo) UpdateStatusByOrderID(ctx context.Context, orderID int64, fromStatus enums.PaymentStatus, toStatus enums.PaymentStatus, at time.Time) error {

	updates := map[string]interface{}{
		"status": toStatus,
	}

	switch toStatus {
	case enums.PaymentStatusWaitPay:
		updates["paid_at"] = at

	case enums.PaymentStatusPartRefund,
		enums.PaymentStatusFullRefund:
		updates["refunded_at"] = at
	}

	res := p.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where(
			"order_id = ? and status = ?",
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

func (p *paymentRepo) UpdateRefundInfo(
	ctx context.Context,
	paymentID int64,
	refundAmount float64,
	status enums.PaymentStatus,
	at time.Time,
) error {

	result := p.db.
		WithContext(ctx).
		Model(&model.Payment{}).
		Where(
			"id = ?",
			paymentID,
		).
		Updates(map[string]interface{}{
			"refund_amount": refundAmount,
			"status":        status,
			"updated_at":    at,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrPaymentNotFound
	}

	return nil
}
