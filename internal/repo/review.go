package repo

import (
	"CampusTake/internal/model"
	"context"

	"gorm.io/gorm"
)

type ReviewRepo interface {
	Create(ctx context.Context, review *model.Review) error
	GetReviewCountAndScoreByRiderID(ctx context.Context, riderID int64) (int64, float64, error)
}

type reviewRepo struct {
	db *gorm.DB
}

func NewReviewRepo(db *gorm.DB) ReviewRepo {
	return &reviewRepo{db: db}
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
