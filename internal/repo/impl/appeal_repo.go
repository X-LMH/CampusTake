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

	"gorm.io/gorm"
)

type AppealRepo interface {
	Create(ctx context.Context, appeal *model.Appeal) error

	GetList(
		ctx context.Context,
		req query.AppealQuery,
	) (*response.PageResult, error)

	ExistsPendingByOrderIDAndType(
		ctx context.Context,
		orderID int64,
		appealType enums.AppealType,
	) (bool, error)
	GetByIDAndUserID(ctx context.Context, id, userID int64) (*model.Appeal, error)

	UpdateStatus(
		ctx context.Context,
		id int64,
		fromStatus enums.AppealStatus,
		toStatus enums.AppealStatus,
	) error
}

type appealRepo struct {
	db *gorm.DB
}

func NewAppealRepo(db *gorm.DB) AppealRepo {
	return &appealRepo{db: db}
}

func (a *appealRepo) Create(ctx context.Context, appeal *model.Appeal) error {
	return a.db.WithContext(ctx).Create(appeal).Error
}

func (a *appealRepo) GetList(
	ctx context.Context,
	req query.AppealQuery,
) (*response.PageResult, error) {

	var (
		appeals []model.Appeal
		total   int64
	)

	dbQuery := a.db.WithContext(ctx).
		Model(&model.Appeal{})

	// =========================
	// 动态条件
	// =========================

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

	// 申诉类型
	if req.Type != nil && *req.Type > 0 {
		dbQuery = dbQuery.Where("type = ?", *req.Type)
	}

	// =========================
	// 查询总数
	// =========================

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// =========================
	// 分页查询
	// =========================

	if err := dbQuery.
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
	appealType enums.AppealType,
) (bool, error) {

	var count int64

	err := a.db.WithContext(ctx).
		Model(&model.Appeal{}).
		Where(
			"order_id = ? AND type = ? AND status = ?",
			orderID,
			appealType,
			enums.AppealStatusPending,
		).
		Count(&count).Error

	return count > 0, err
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

func (a *appealRepo) UpdateStatus(
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
