// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/repo/query"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderLogic) PayOrder(req *types.PayOrderRequest) error {
	userID := ctxx.MustUserID(l.ctx)

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		// -----------------------------
		// 1. 查询订单和支付记录
		// -----------------------------
		order, err := tx.Order.GetByIDAndUserID(l.ctx, req.OrderID, userID)
		if err != nil {
			l.Errorf("查询支付订单失败，orderID=%d，userID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		payment, err := tx.Payment.GetByOrderID(l.ctx, req.OrderID)
		if err != nil {
			l.Errorf("查询支付记录失败，orderID=%d，userID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// -----------------------------
		// 2. 幂等性检查
		// -----------------------------
		if payment.Status == enums.PaymentStatusPaid && order.PaymentStatus == enums.OrderPayStatusPaid {
			return nil // 已支付，直接返回成功
		}

		if order.Status != enums.OrderPendingPay || payment.Status != enums.PaymentStatusWaitPay {
			l.Errorf("订单状态不允许支付，orderID=%d，userID=%d，orderStatus=%v，paymentStatus=%v", req.OrderID, userID, order.Status, payment.Status)
			return errs.ErrOrderStatusInvalid
		}
		toStatus := enums.OrderPendingGrab
		if !enums.CheckOrderStatusFlow(order.Status, toStatus) {
			l.Errorf("订单流转校验失败，orderID=%d，userID=%d，from=%v，to=%v", req.OrderID, userID, order.Status, toStatus)
			return errs.ErrOrderStatusInvalid
		}

		paidAt := time.Now()
		// -----------------------------
		// 3. 更新支付状态（先更新支付状态保证原子性）
		// -----------------------------
		if err := tx.Payment.UpdateStatusByOrderID(
			l.ctx,
			req.OrderID,
			enums.PaymentStatusWaitPay,
			enums.PaymentStatusPaid,
			paidAt,
		); err != nil {
			l.Errorf("更新支付状态失败, orderID: %d, err: %v", req.OrderID, err)
			return err
		}

		// -----------------------------
		// 4. 更新订单状态
		// -----------------------------
		if err := tx.Order.UpdateStatusAndTime(
			l.ctx,
			query.OrderStatusUpdateQuery{
				OrderID: order.ID,
				UserID:  &userID,
			},
			order.Status,
			toStatus,
			paidAt,
			map[string]interface{}{
				"payment_status": enums.OrderPayStatusPaid,
			},
		); err != nil {
			l.Errorf("更新订单支付状态失败，orderID=%d，userID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// -----------------------------
		// 5. 写订单日志
		// -----------------------------
		log := &model.OrderLog{
			OrderID:      req.OrderID,
			FromStatus:   order.Status,
			ToStatus:     toStatus,
			OperatorType: enums.OperatorTypeUser,
			OperatorID:   userID,
			Remark:       "用户支付订单",
		}

		return tx.Order.CreateLog(l.ctx, log)
	})

	if err != nil {
		return err
	}

	return nil
}
