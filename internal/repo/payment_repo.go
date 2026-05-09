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
	UpdateByOrderID(ctx context.Context, orderID int64, status enums.PayStatus) error
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

// 更新支付状态（原子操作），自动填充 PaidAt/RefundedAt
func (p *paymentRepo) UpdateByOrderID(ctx context.Context, orderID int64, status enums.PayStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == enums.PayStatusPaid {
		updates["paid_at"] = time.Now()
	} else if status == enums.PayStatusRefunded {
		updates["refunded_at"] = time.Now()
	}

	res := p.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("order_id = ? AND status = ?", orderID, enums.PayStatusPending).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errs.ErrPaymentStatusInvalid
	}
	return nil
}
