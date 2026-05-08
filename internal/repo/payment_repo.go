package repo

import (
	"CampusTake/internal/model"
	"context"

	"gorm.io/gorm"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment *model.Payment) error
}

type paymentRepo struct {
	db *gorm.DB
}

func (p *paymentRepo) Create(ctx context.Context, payment *model.Payment) error {
	return p.db.WithContext(ctx).Create(payment).Error
}
