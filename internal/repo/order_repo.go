package repo

import (
	"CampusTake/internal/model"
	"context"

	"gorm.io/gorm"
)

type OrderRepo interface {
	Create(ctx context.Context, order *model.Order) error
}

type orderRepo struct {
	db *gorm.DB
}

func (o *orderRepo) Create(ctx context.Context, order *model.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}
