package impl

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RiderRepo interface {
	GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error)
	GetProfileByRiderID(ctx context.Context, riderID int64) (*model.RiderProfile, error)
	UpsertProfile(ctx context.Context, profile *model.RiderProfile) error
	UpdateStatusByUserID(ctx context.Context, userID int64, status enums.RiderAuditStatus) error
	UpdateStatusAndRemarkByUserID(ctx context.Context, userID int64, status enums.RiderAuditStatus, remark string) error
	GetProfileList(ctx context.Context, status enums.RiderAuditStatus, page, pageSize int) (*response.PageResult, error)
	UpdateCompleteStatsByRiderID(ctx context.Context, riderID int64, completedCount int, completionRate float64) error
	UpdateRatingByRiderID(ctx context.Context, riderID int64, count int64, avg float64) error
	CreateLog(ctx context.Context, log *model.RiderAuditLog) error
}

type riderRepo struct {
	db *gorm.DB
}

func NewRiderRepo(db *gorm.DB) RiderRepo {
	return &riderRepo{db: db}
}

func (r *riderRepo) GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error) {
	profile := new(model.RiderProfile)
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return profile, nil
}

func (r *riderRepo) GetProfileByRiderID(ctx context.Context, riderID int64) (*model.RiderProfile, error) {
	profile := new(model.RiderProfile)
	err := r.db.WithContext(ctx).Where("id = ?", riderID).First(profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return profile, nil
}

func (r *riderRepo) UpsertProfile(ctx context.Context, profile *model.RiderProfile) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"real_name", "student_no", "id_card_no",
			"dormitory_building", "dormitory_room",
			"campus_card_front", "campus_card_back",
			"audit_status", "updated_at",
		}),
	}).Create(profile).Error
}

func (r *riderRepo) UpdateStatusByUserID(ctx context.Context, userID int64, status enums.RiderAuditStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", userID).
		Update("audit_status", status).Error
}

func (r *riderRepo) UpdateStatusAndRemarkByUserID(ctx context.Context, userID int64, status enums.RiderAuditStatus, remark string) error {
	return r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"audit_status": status,
			"audit_remark": remark,
		}).Error
}

func (r *riderRepo) GetProfileList(ctx context.Context, status enums.RiderAuditStatus, page, pageSize int) (*response.PageResult, error) {
	var list []*model.RiderProfile
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RiderProfile{})
	if status > 0 {
		query = query.Where("audit_status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	err := query.Scopes(db.Paginate(page, pageSize)).
		Order("updated_at DESC").
		Find(&list).Error

	return response.NewPageResult(total, list), err
}

func (r *riderRepo) UpdateRatingByRiderID(ctx context.Context, riderID int64, count int64, avg float64) error {
	return r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("id = ?", riderID).
		Updates(map[string]interface{}{
			"rating_count": count,
			"rating_avg":   avg,
		}).Error
}

func (r *riderRepo) UpdateCompleteStatsByRiderID(
	ctx context.Context,
	riderID int64,
	completedCount int,
	completionRate float64,
) error {

	return r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", riderID).
		Updates(map[string]interface{}{
			"completed_order_count": completedCount,
			"completion_rate":       completionRate,
		}).Error
}

func (r *riderRepo) CreateLog(ctx context.Context, log *model.RiderAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
