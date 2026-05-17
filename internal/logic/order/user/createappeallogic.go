// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"

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

func (l *CreateAppealLogic) CreateAppeal(
	req *types.CreateAppealRequest,
) error {

	userID := ctxx.MustUserID(l.ctx)
	role := ctxx.GetRole(l.ctx)

	var (
		order *model.Order
		err   error
	)

	switch role {

	case enums.RoleUser:
		order, err = l.svcCtx.Repo.Order.GetByIDAndUserID(
			l.ctx,
			req.OrderID,
			userID,
		)

	case enums.RoleRider:
		order, err = l.svcCtx.Repo.Order.GetByIDAndRiderID(
			l.ctx,
			req.OrderID,
			userID,
		)

	default:
		return errs.ErrUserPermissionDenied
	}

	if err != nil {
		l.Errorf(
			"查询订单失败 orderID=%d userID=%d err=%v",
			req.OrderID,
			userID,
			err,
		)
		return err
	}

	if !order.Status.CanAppealOrderStatus() {
		return errs.ErrOrderStatusInvalid
	}

	appealType := role.ToAppealType()

	exist, err := l.svcCtx.Repo.Appeal.
		ExistsPendingByOrderIDAndType(
			l.ctx,
			req.OrderID,
			appealType,
		)

	if err != nil {
		l.Errorf(
			"查询申诉记录失败 orderID=%d userID=%d err=%v",
			req.OrderID,
			userID,
			err,
		)
		return err
	}

	if exist {
		return errs.ErrOrderHasPendingAppeal
	}

	appeal := &model.Appeal{
		OrderID:     order.ID,
		ApplicantID: userID,
		Type:        appealType,
		Content:     req.Content,
		Status:      enums.AppealStatusPending,
	}

	if err = l.svcCtx.Repo.Appeal.Create(
		l.ctx,
		appeal,
	); err != nil {

		l.Errorf(
			"创建申诉失败 orderID=%d userID=%d err=%v",
			req.OrderID,
			userID,
			err,
		)

		return err
	}

	return nil
}
