// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RiderCancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRiderCancelOrderLogic(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
) *RiderCancelOrderLogic {

	return &RiderCancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RiderCancelOrderLogic) RiderCancelOrder(
	req *types.RiderCancelOrderRequest,
) error {

	userID := ctxx.MustUserID(l.ctx)

	// =========================
	// 查询订单
	// =========================

	order, err := l.svcCtx.Repo.Order.GetByID(l.ctx, req.OrderID)
	if err != nil {
		return err
	}

	// =========================
	// 校验权限
	// =========================

	if order.RiderID == nil || *order.RiderID != userID {
		return errs.ErrOrderNoPermission
	}

	// =========================
	// 只允许：
	// 已接单 -> 待接单
	// =========================

	if order.Status != enums.OrderAccepted {
		return errs.ErrOrderStatusInvalid
	}

	at := time.Now()

	err = l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		fromStatus := order.Status
		toStatus := enums.OrderPendingGrab

		// =========================
		// 回退订单状态
		// =========================

		err := tx.Order.RiderUpdateStatusAndTime(
			l.ctx,
			order.ID,
			userID,
			fromStatus,
			toStatus,
			at,
			map[string]interface{}{
				"rider_id":      nil,
				"accepted_at":   nil,
				"cancel_reason": req.Reason,
			},
		)
		if err != nil {
			l.Errorf("骑手取消订单失败，orderID：%d, riderID: %d, err: %v", order.ID, userID, err)
			return err
		}

		// =========================
		// 创建订单日志
		// =========================

		log := &model.OrderLog{
			OrderID:      order.ID,
			FromStatus:   fromStatus,
			ToStatus:     toStatus,
			OperatorType: enums.OperatorTypeRider,
			OperatorID:   userID,
			Remark:       fmt.Sprintf("骑手取消接单，原因：%s", req.Reason),
			CreatedAt:    at,
		}

		return tx.Order.CreateLog(l.ctx, log)
	})

	if err != nil {
		return err
	}

	return nil
}
