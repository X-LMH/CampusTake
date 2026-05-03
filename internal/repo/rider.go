package repo

import (
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Rider interface {
	GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error)
	UpsertProfile(ctx context.Context, profile *model.RiderProfile) error
	UpdateStatus(ctx context.Context, userID int64, status enum.RiderAuditStatus) error
}

type riderRepo struct {
	db *gorm.DB
}

func (r *riderRepo) GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error) {
	profile := new(model.RiderProfile)
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrUserNotFound
		}
		return nil, err
	}
	return profile, nil
}

func (r *riderRepo) UpsertProfile(ctx context.Context, profile *model.RiderProfile) error {
	// 使用 OnConflict 处理 user_id 冲突的情况
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		// 判定冲突的列：user_id (对应你表里的 uk_user_id)
		Columns: []clause.Column{{Name: "user_id"}},
		// 如果冲突，则执行更新以下字段
		DoUpdates: clause.AssignmentColumns([]string{
			"real_name", "student_no", "id_card_no",
			"dormitory_building", "dormitory_room",
			"campus_card_front", "campus_card_back",
			"audit_status", "updated_at",
		}),
	}).Create(profile).Error
}

func (r *riderRepo) UpdateStatus(ctx context.Context, userID int64, status enum.RiderAuditStatus) error {
	err := r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", userID).
		Update("audit_status", status).Error

	if err != nil {
		return err
	}
	return nil
}
