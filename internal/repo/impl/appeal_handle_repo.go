package impl

import (
	"CampusTake/internal/model"
	"context"
	"errors" // 引入 standard library 的 errors

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppealHandleRepo interface {
	Create(
		ctx context.Context,
		handle *model.AppealHandle,
	) error

	GetLastHandleByAppealID(
		ctx context.Context,
		appealID int64,
	) (*model.AppealHandle, error)
}

type appealHandleRepo struct {
	RepoBase
}

func NewAppealHandleRepo(db *gorm.DB, rdb redis.Cmdable) AppealHandleRepo {
	return &appealHandleRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func (a *appealHandleRepo) Create(
	ctx context.Context,
	handle *model.AppealHandle,
) error {
	return a.db.WithContext(ctx).Create(handle).Error
}

func (a *appealHandleRepo) GetLastHandleByAppealID(
	ctx context.Context,
	appealID int64,
) (*model.AppealHandle, error) {
	var handle model.AppealHandle

	err := a.db.WithContext(ctx).
		Where("appeal_id = ?", appealID).
		Order("created_at DESC, id DESC").
		First(&handle).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &handle, nil
}
