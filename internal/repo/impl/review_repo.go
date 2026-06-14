package impl

import (
	"CampusTake/internal/model"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReviewRepo interface {
	Create(ctx context.Context, review *model.Review) error
	GetReviewCountAndScoreByRiderID(ctx context.Context, riderID int64) (int64, float64, error)
	ExistByOrderID(ctx context.Context, orderID int64) (bool, error)
}

type reviewRepo struct {
	RepoBase
}

func NewReviewRepo(db *gorm.DB, rdb *redis.Client) ReviewRepo {
	return &reviewRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func (r *reviewRepo) Create(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepo) GetReviewCountAndScoreByRiderID(ctx context.Context, riderID int64) (int64, float64, error) {
	var count int64
	var avg float64

	query := r.db.WithContext(ctx).Model(&model.Review{}).Where("rider_id = ?", riderID)
	err := query.Select("COUNT(*)").Scan(&count).Error
	if err != nil {
		return 0, 0, err
	}
	err = query.Select("AVG(score)").Scan(&avg).Error
	if err != nil {
		return 0, 0, err
	}
	return count, avg, err
}

func (r *reviewRepo) ExistByOrderID(ctx context.Context, orderID int64) (bool, error) {
	var review model.Review

	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&review).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
