// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/repo/query"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"CampusTake/pkg/utils"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAppealLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAppealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppealLogic {
	return &CreateAppealLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAppealLogic) CreateAppeal(req *types.CreateAppealRequest) error {
	userID := ctxx.MustUserID(l.ctx)
	role := ctxx.GetRole(l.ctx)

	var (
		order *model.Order
		err   error
	)

	// 1. 先查订单并校验归属
	switch role {
	case enums.RoleUser:
		order, err = l.svcCtx.Repo.Order.GetByIDAndUserID(l.ctx, req.OrderID, userID)
	case enums.RoleRider:
		order, err = l.svcCtx.Repo.Order.GetByIDAndRiderID(l.ctx, req.OrderID, userID)
	default:
		l.Errorf("创建申诉权限不足，orderID=%d，userID=%d，role=%v", req.OrderID, userID, role)
		return errs.ErrUserPermissionDenied
	}

	if err != nil {
		l.Errorf("查询订单失败 orderID=%d userID=%d err=%v", req.OrderID, userID, err)
		return err
	}

	// 2. 校验订单状态是否允许申诉
	if !order.Status.CanAppealOrderStatus() {
		l.Errorf("订单状态不允许申诉，orderID=%d，userID=%d，status=%v", req.OrderID, userID, order.Status)
		return errs.ErrOrderStatusInvalid
	}

	appealRole := role.ToAppealType()

	// 3. 事务内创建/更新申诉 + 更新订单状态
	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		latestAppeal, err := tx.Appeal.GetLatestByOrderIDAndApplicant(
			l.ctx,
			order.ID,
			userID,
			appealRole,
		)
		if err != nil && !errors.Is(err, errs.ErrAppealNotFound) {
			l.Errorf("查询最新申诉失败 orderID=%d userID=%d err=%v", req.OrderID, userID, err)
			return err
		}

		if latestAppeal != nil && latestAppeal.Status == enums.AppealStatusCancelled {
			err = tx.Appeal.UpdateReapplyByID(
				l.ctx,
				latestAppeal.ID,
				enums.AppealStatusCancelled,
				query.ReapplyAppealParams{
					AppealType:   req.AppealType,
					Content:      req.Content,
					EvidenceUrls: utils.ParseStringSliceToJSON(req.EvidenceUrls),
				},
			)
			if err != nil {
				l.Errorf("更新申诉失败 orderID=%d userID=%d err=%v", req.OrderID, userID, err)
				return err
			}
		} else {
			appeal := &model.Appeal{
				OrderID:       order.ID,
				ApplicantID:   userID,
				ApplicantRole: appealRole,
				AppealType:    req.AppealType,
				Content:       req.Content,
				EvidenceUrls:  utils.ParseStringSliceToJSON(req.EvidenceUrls),
				Status:        enums.AppealStatusPending,
			}

			err = tx.Appeal.Create(l.ctx, appeal)
			if err != nil {
				l.Errorf("创建申诉失败 orderID=%d userID=%d err=%v", req.OrderID, userID, err)
				return err
			}
		}

		if order.AppealStatus != enums.OrderAppealStatusOngoing {
			err = tx.Order.UpdateAppealStatusByID(
				l.ctx,
				order.ID,
				order.AppealStatus,
				enums.OrderAppealStatusOngoing,
			)
			if err != nil {
				l.Errorf("更新订单申诉状态失败 orderID=%d err=%v", order.ID, err)
				return err
			}
		}

		return nil
	})
}
