package impl

import (
	"CampusTake/internal/model"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RefundRepo interface {
	Create(ctx context.Context, refund *model.PaymentRefund) error
}

type refundRepo struct {
	RepoBase
}

func NewRefundRepo(db *gorm.DB, rdb *redis.Client) RefundRepo {
	return &refundRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func (r *refundRepo) Create(ctx context.Context, refund *model.PaymentRefund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}
