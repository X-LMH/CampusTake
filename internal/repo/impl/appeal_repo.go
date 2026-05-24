package impl

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo/query"
	"CampusTake/pkg/db"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppealRepo interface {
	Create(ctx context.Context, appeal *model.Appeal) error

	GetList(
		ctx context.Context,
		req query.AppealListQuery,
	) (*response.PageResult, error)

	ExistsPendingByOrderIDAndType(
		ctx context.Context,
		orderID int64,
		appealType enums.AppealApplicantRole,
	) (bool, error)

	GetByID(ctx context.Context, id int64) (*model.Appeal, error)

	GetByIDForUpdate(ctx context.Context, id int64) (*model.Appeal, error)

	GetLatestByOrderIDAndApplicant(
		ctx context.Context,
		orderID int64,
		applicantID int64,
		applicantRole enums.AppealApplicantRole,
	) (*model.Appeal, error)

	GetByIDAndUserID(ctx context.Context, id, userID int64) (*model.Appeal, error)

	UpdateStatusByID(
		ctx context.Context,
		id int64,
		fromStatus enums.AppealStatus,
		toStatus enums.AppealStatus,
	) error

	UpdateHandleResult(
		ctx context.Context,
		id int64,
		fromStatus, toStatus enums.AppealStatus,
		params query.HandleAppealParams,
	) error

	UpdateReapplyByID(
		ctx context.Context,
		id int64,
		fromStatus enums.AppealStatus,
		params query.ReapplyAppealParams,
	) error
}

type appealRepo struct {
	RepoBase
}

func NewAppealRepo(db *gorm.DB, rdb redis.Cmdable) AppealRepo {
	return &appealRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func (a *appealRepo) Create(ctx context.Context, appeal *model.Appeal) error {
	return a.db.WithContext(ctx).Create(appeal).Error
}

func (a *appealRepo) GetList(
	ctx context.Context,
	req query.AppealListQuery,
) (*response.PageResult, error) {

	var (
		appeals []model.Appeal
		total   int64
	)

	buildQuery := func() *gorm.DB {
		dbQuery := a.db.WithContext(ctx).Model(&model.Appeal{})

		// 状态
		if req.Status != nil && *req.Status > 0 {
			dbQuery = dbQuery.Where("status = ?", *req.Status)
		}

		// 申诉人
		if req.ApplicantID != nil && *req.ApplicantID > 0 {
			dbQuery = dbQuery.Where("applicant_id = ?", *req.ApplicantID)
		}

		// 订单ID
		if req.OrderID != nil && *req.OrderID > 0 {
			dbQuery = dbQuery.Where("order_id = ?", *req.OrderID)
		}

		// 申诉用户类型
		if req.ApplicantRole != nil && *req.ApplicantRole > 0 {
			dbQuery = dbQuery.Where("applicant_role = ?", *req.ApplicantRole) // NOTE: 这里改成 applicant_role，和字段语义保持一致
		}

		return dbQuery
	}

	// NOTE: count 单独走一份 query，避免分页条件污染统计结果
	countQuery := buildQuery()
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// NOTE: 列表查询单独走一份 query
	listQuery := buildQuery()
	if err := listQuery.
		Scopes(db.Paginate(req.Page, req.Size)).
		Order("created_at DESC").
		Find(&appeals).
		Error; err != nil {
		return nil, err
	}

	return response.NewPageResult(total, appeals), nil
}

func (a *appealRepo) ExistsPendingByOrderIDAndType(
	ctx context.Context,
	orderID int64,
	appealType enums.AppealApplicantRole,
) (bool, error) {

	var count int64

	err := a.db.WithContext(ctx).
		Model(&model.Appeal{}).
		Where(
			"order_id = ? AND applicant_role = ? AND status = ?", // NOTE: 这里和上面统一，用 applicant_role
			orderID,
			appealType,
			enums.AppealStatusPending,
		).
		Count(&count).Error

	return count > 0, err
}

func (a *appealRepo) GetByID(ctx context.Context, id int64) (*model.Appeal, error) {
	appeal := new(model.Appeal)

	err := a.db.WithContext(ctx).Where("id = ?", id).First(appeal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAppealNotFound
		}
		return nil, err
	}
	return appeal, nil
}

func (a *appealRepo) GetByIDForUpdate(ctx context.Context, id int64) (*model.Appeal, error) {

	var appeal model.Appeal

	err := a.db.WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ?", id).
		First(&appeal).Error

	if err != nil {
		return nil, err
	}

	return &appeal, nil
}

func (a *appealRepo) GetLatestByOrderIDAndApplicant(
	ctx context.Context,
	orderID int64,
	applicantID int64,
	applicantRole enums.AppealApplicantRole,
) (*model.Appeal, error) {
	appeal := new(model.Appeal)

	err := a.db.WithContext(ctx).
		Where(
			"order_id = ? AND applicant_id = ? AND applicant_role = ?",
			orderID,
			applicantID,
			applicantRole,
		).
		Order("id DESC").
		First(appeal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAppealNotFound
		}
		return nil, err
	}

	return appeal, nil
}

func (a *appealRepo) GetByIDAndUserID(ctx context.Context, id, userID int64) (*model.Appeal, error) {
	appeal := new(model.Appeal)

	err := a.db.WithContext(ctx).Where("id = ? AND applicant_id = ?", id, userID).First(appeal).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAppealNotFound
		}
		return nil, err
	}
	return appeal, nil
}

func (a *appealRepo) UpdateStatusByID(
	ctx context.Context,
	id int64,
	fromStatus enums.AppealStatus,
	toStatus enums.AppealStatus,
) error {
	if !fromStatus.CanTransferTo(toStatus) {
		return errs.ErrAppealStatusInvalid
	}

	result := a.db.WithContext(ctx).
		Model(&model.Appeal{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Update("status", toStatus)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrAppealStatusChanged
	}

	return nil
}

func (a *appealRepo) UpdateHandleResult(
	ctx context.Context,
	id int64,
	fromStatus, toStatus enums.AppealStatus,
	params query.HandleAppealParams,
) error {
	if !fromStatus.CanTransferTo(toStatus) {
		return errs.ErrAppealStatusInvalid
	}

	result := a.db.WithContext(ctx).
		Model(&model.Appeal{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Updates(map[string]any{
			"status":     toStatus,
			"handled_by": params.HandledBy,
			"handled_at": params.HandledAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrAppealStatusChanged
	}

	return nil
}

func (a *appealRepo) UpdateReapplyByID(
	ctx context.Context,
	id int64,
	fromStatus enums.AppealStatus,
	params query.ReapplyAppealParams,
) error {
	result := a.db.WithContext(ctx).
		Model(&model.Appeal{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Updates(map[string]any{
			"appeal_type":   params.AppealType,
			"content":       params.Content,
			"evidence_urls": params.EvidenceUrls,
			"status":        enums.AppealStatusPending,
			"handled_by":    nil,
			"handled_at":    nil,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrAppealStatusChanged
	}

	return nil
}
