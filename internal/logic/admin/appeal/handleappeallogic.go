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
	"CampusTake/pkg/snowflake"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleAppealLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleAppealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleAppealLogic {
	return &HandleAppealLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleAppealLogic) HandleAppeal(req *types.HandleAppealRequest) error {
	adminID := ctxx.MustUserID(l.ctx)
	at := time.Now()

	toStatus, handleResult, err := l.resolveHandleResult(req.Result)
	if err != nil {
		return err
	}

	params := query.HandleAppealParams{
		HandleRemark:   req.Remark,
		RefundAmount:   req.RefundAmount,
		PunishRider:    req.PunishRider,
		TerminateOrder: req.TerminateOrder,
		HandledBy:      adminID,
		HandledAt:      at,
	}
	if err := params.Validate(req.Result); err != nil {
		return err
	}

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		appeal, err := tx.Appeal.GetByIDForUpdate(l.ctx, req.AppealID)
		if err != nil {
			return err
		}
		if appeal.Status != enums.AppealStatusPending {
			return errs.ErrAppealHaveProcessed
		}

		order, err := tx.Order.GetByIDForUpdate(l.ctx, appeal.OrderID)
		if err != nil {
			return err
		}

		if err := tx.Appeal.UpdateHandleResult(
			l.ctx,
			req.AppealID,
			enums.AppealStatusPending,
			toStatus,
			params,
		); err != nil {
			return err
		}

		if err := tx.AppealHandle.Create(l.ctx, &model.AppealHandle{
			AppealID:       req.AppealID,
			HandlerID:      adminID,
			Result:         handleResult,
			Remark:         req.Remark,
			RefundAmount:   req.RefundAmount,
			PunishRider:    req.PunishRider,
			PunishUser:     req.PunishUser,
			TerminateOrder: req.TerminateOrder,
			CreatedAt:      at,
		}); err != nil {
			return err
		}

		if err := l.updateOrderAppealStatus(tx, order, toStatus); err != nil {
			return err
		}

		if toStatus == enums.AppealStatusRejected {
			return nil
		}

		if err := l.adjustRiderCompletedCount(tx, appeal, order); err != nil {
			return err
		}

		if req.TerminateOrder == 1 {
			if err := l.handleTerminateOrder(tx, adminID, order, at); err != nil {
				return err
			}
		}

		// Rider punishment is intentionally reserved for a later rules module.

		if req.RefundAmount > 0 {
			if err := l.handleRefund(tx, adminID, order, req.RefundAmount, at); err != nil {
				return err
			}
		}

		return nil
	})
}

func (l *HandleAppealLogic) resolveHandleResult(
	result int8,
) (enums.AppealStatus, enums.AppealHandleResult, error) {
	switch result {
	case int8(enums.AppealHandleResultSuccess):
		return enums.AppealStatusApproved, enums.AppealHandleResultSuccess, nil
	case int8(enums.AppealHandleResultFailed):
		return enums.AppealStatusRejected, enums.AppealHandleResultFailed, nil
	default:
		return enums.AppealStatusDefault, enums.AppealHandleResultDefault, errs.ErrInvalidParam
	}
}

func (l *HandleAppealLogic) updateOrderAppealStatus(
	tx *repo.RepoTx,
	order *model.Order,
	appealStatus enums.AppealStatus,
) error {
	var toStatus enums.OrderAppealStatus
	switch appealStatus {
	case enums.AppealStatusApproved:
		toStatus = enums.OrderAppealStatusApproved
	case enums.AppealStatusRejected:
		toStatus = enums.OrderAppealStatusRejected
	default:
		return errs.ErrAppealStatusInvalid
	}

	if order.AppealStatus == toStatus {
		return nil
	}
	if order.AppealStatus != enums.OrderAppealStatusOngoing {
		return errs.ErrOrderAppealStatusChanged
	}

	return tx.Order.UpdateAppealStatusByID(
		l.ctx,
		order.ID,
		enums.OrderAppealStatusOngoing,
		toStatus,
	)
}

func (l *HandleAppealLogic) adjustRiderCompletedCount(
	tx *repo.RepoTx,
	appeal *model.Appeal,
	order *model.Order,
) error {
	if appeal.ApplicantRole != enums.AppealApplicantRoleUser {
		return nil
	}
	if order.RiderID == nil {
		return nil
	}
	if order.Status != enums.OrderCompleted {
		return nil
	}

	return tx.Rider.DecrementCompletedOrderCount(l.ctx, *order.RiderID)
}

func (l *HandleAppealLogic) handleTerminateOrder(
	tx *repo.RepoTx,
	adminID int64,
	order *model.Order,
	at time.Time,
) error {
	switch order.Status {
	case enums.OrderAccepted:
		return l.moveOrderStatus(
			tx,
			adminID,
			order,
			enums.OrderException,
			at,
			"Admin approved appeal and terminated the order",
		)
	case enums.OrderDelivering, enums.OrderDelivered:
		return tx.Order.CreateLog(l.ctx, &model.OrderLog{
			OrderID:      order.ID,
			FromStatus:   order.Status,
			ToStatus:     order.Status,
			OperatorType: enums.OperatorTypeAdmin,
			OperatorID:   adminID,
			Remark:       "Admin approved appeal; order will finish current delivery without rider completion credit",
			CreatedAt:    at,
		})
	case enums.OrderCompleted:
		return tx.Order.CreateLog(l.ctx, &model.OrderLog{
			OrderID:      order.ID,
			FromStatus:   order.Status,
			ToStatus:     order.Status,
			OperatorType: enums.OperatorTypeAdmin,
			OperatorID:   adminID,
			Remark:       "Admin approved appeal after completion; rider completion credit has been revoked",
			CreatedAt:    at,
		})
	default:
		return errs.ErrOrderStatusInvalid
	}
}

func (l *HandleAppealLogic) moveOrderStatus(
	tx *repo.RepoTx,
	adminID int64,
	order *model.Order,
	toStatus enums.OrderStatus,
	at time.Time,
	remark string,
) error {
	if err := tx.Order.UpdateStatusAndTime(
		l.ctx,
		query.OrderStatusUpdateQuery{OrderID: order.ID},
		order.Status,
		toStatus,
		at,
		nil,
	); err != nil {
		return err
	}

	return tx.Order.CreateLog(l.ctx, &model.OrderLog{
		OrderID:      order.ID,
		FromStatus:   order.Status,
		ToStatus:     toStatus,
		OperatorType: enums.OperatorTypeAdmin,
		OperatorID:   adminID,
		Remark:       remark,
		CreatedAt:    at,
	})
}

func (l *HandleAppealLogic) handleRefund(
	tx *repo.RepoTx,
	adminID int64,
	order *model.Order,
	refundAmount float64,
	at time.Time,
) error {
	payment, err := tx.Payment.GetByOrderIDForUpdate(l.ctx, order.ID)
	if err != nil {
		return err
	}

	switch payment.Status {
	case enums.PaymentStatusWaitPay:
		return errs.ErrOrderNotPaid
	case enums.PaymentStatusFullRefund:
		return errs.ErrPaymentAlreadyRefunded
	}

	remainRefundAmount := payment.Amount - payment.RefundAmount
	if refundAmount <= 0 || refundAmount > remainRefundAmount {
		return errs.ErrRefundAmountInvalid
	}

	refund := &model.PaymentRefund{
		PaymentID:  payment.ID,
		RefundNo:   snowflake.GenerateID(),
		Amount:     refundAmount,
		Status:     enums.RefundStatusSuccess,
		Reason:     "Admin appeal refund",
		OperatorID: adminID,
		RefundedAt: &at,
		CreatedAt:  at,
	}

	if err := tx.Refund.Create(l.ctx, refund); err != nil {
		return err
	}

	newRefundAmount := payment.RefundAmount + refundAmount
	newPaymentStatus := enums.PaymentStatusPartRefund
	if newRefundAmount >= payment.Amount {
		newPaymentStatus = enums.PaymentStatusFullRefund
	}

	if err := tx.Payment.UpdateRefundInfo(
		l.ctx,
		payment.ID,
		newRefundAmount,
		newPaymentStatus,
		at,
	); err != nil {
		return err
	}

	if newPaymentStatus == enums.PaymentStatusFullRefund {
		return tx.Order.UpdatePaymentStatusByIDAndUserID(
			l.ctx,
			order.ID,
			order.UserID,
			enums.OrderPayStatusRefunded,
		)
	}

	return nil
}
