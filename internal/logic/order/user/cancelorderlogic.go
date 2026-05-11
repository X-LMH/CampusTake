// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"
	"fmt"
	"time"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderRequest) error {

	userID := ctxx.MustUserID(l.ctx)

	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(
		l.ctx,
		req.OrderID,
		userID,
	)
	if err != nil {
		return err
	}

	at := time.Now()

	switch order.Status {

	// =========================
	// 未支付 -> 直接取消
	// =========================

	case enums.OrderPendingPay:

		err = l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

			fromStatus := order.Status
			toStatus := enums.OrderUserCancelled

			// 更新订单
			err := tx.Order.UserUpdateStatusAndTime(
				l.ctx,
				order.ID,
				userID,
				fromStatus,
				toStatus,
				at,
				map[string]interface{}{
					"cancel_reason": req.Reason,
				},
			)
			if err != nil {
				return err
			}

			// 创建日志
			log := &model.OrderLog{
				OrderID:      order.ID,
				FromStatus:   fromStatus,
				ToStatus:     toStatus,
				OperatorType: enums.OperatorTypeUser,
				OperatorID:   userID,
				Remark:       fmt.Sprintf("用户取消订单，原因：%s", req.Reason),
				CreatedAt:    at,
			}

			return tx.Order.CreateLog(l.ctx, log)
		})

	// =========================
	// 已支付未接单 -> 取消 + 退款
	// =========================

	case enums.OrderPendingGrab:
		if order.PaymentStatus != enums.OrderPayStatusPaid {
			return errs.ErrPaymentStatusInvalid
		}
		err = l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

			fromStatus := order.Status
			toStatus := enums.OrderUserCancelled

			// 更新订单状态
			err := tx.Order.UserUpdateStatusAndTime(
				l.ctx,
				order.ID,
				userID,
				fromStatus,
				toStatus,
				at,
				map[string]interface{}{
					"cancel_reason":  req.Reason,
					"payment_status": enums.OrderPayStatusRefunded,
					"refunded_at":    at,
				},
			)
			if err != nil {
				return err
			}

			// 更新支付状态
			err = tx.Payment.UpdateStatusByOrderID(
				l.ctx,
				order.ID,
				enums.OrderPayStatusPaid,
				enums.OrderPayStatusRefunded,
				at,
			)
			if err != nil {
				return err
			}

			// 创建日志
			log := &model.OrderLog{
				OrderID:      order.ID,
				FromStatus:   fromStatus,
				ToStatus:     toStatus,
				OperatorType: enums.OperatorTypeUser,
				OperatorID:   userID,
				Remark:       fmt.Sprintf("用户取消订单并退款，原因：%s", req.Reason),
				CreatedAt:    at,
			}

			return tx.Order.CreateLog(l.ctx, log)
		})

	default:
		return errs.ErrOrderCannotCancel
	}

	if err != nil {
		return err
	}

	return nil
}
