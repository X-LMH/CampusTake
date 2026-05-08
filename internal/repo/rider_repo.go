package repo

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

type Rider interface {
	GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error)
	UpsertProfile(ctx context.Context, profile *model.RiderProfile) error
	UpdateStatusByUserID(ctx context.Context, userID int64, status enums.RiderAuditStatus) error
	UpdateStatusAndRemarkByUserID(ctx context.Context, id int64, status enums.RiderAuditStatus, remark string) error
	GetProfileList(ctx context.Context, status enums.RiderAuditStatus, page, pageSize int) (*response.PageResult, error)
	CreateLog(ctx context.Context, log *model.RiderAuditLog) error
}

type riderRepo struct {
	db *gorm.DB
}

func (r *riderRepo) GetProfileByUserID(ctx context.Context, userID int64) (*model.RiderProfile, error) {
	profile := new(model.RiderProfile)
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
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

func (r *riderRepo) UpdateStatusByUserID(ctx context.Context, userID int64, status enums.RiderAuditStatus) error {
	err := r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("user_id = ?", userID).
		Update("audit_status", status).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *riderRepo) GetProfileList(ctx context.Context, status enums.RiderAuditStatus, page, pageSize int) (*response.PageResult, error) {
	var list []*model.RiderProfile
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RiderProfile{})

	if status > 0 {
		query = query.Where("audit_status = ?", status)
	}

	// 先统计 Total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 再执行自动分页查询
	err := query.Scopes(db.Paginate(page, pageSize)).
		Order("updated_at DESC").
		Find(&list).Error

	return response.NewPageResult(total, list), err
}

func (r *riderRepo) CreateLog(ctx context.Context, log *model.RiderAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *riderRepo) UpdateStatusAndRemarkByUserID(ctx context.Context, id int64, status enums.RiderAuditStatus, remark string) error {
	return r.db.WithContext(ctx).
		Model(&model.RiderProfile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"audit_status": status,
			"audit_remark": remark,
		}).Error
}
