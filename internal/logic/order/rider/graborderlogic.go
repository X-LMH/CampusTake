// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/mqs"
	"CampusTake/internal/repo"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"
	"time"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrabOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGrabOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrabOrderLogic {
	return &GrabOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GrabOrderLogic) GrabOrder(req *types.GrabOrderRequest) error {
	userID := ctxx.MustUserID(l.ctx)
	acceptedAt := time.Now()

	// Redis Claim 抢占
	claimed, err := l.svcCtx.Repo.Order.TryClaimGrabOrder(
		l.ctx,
		req.OrderID,
		userID,
	)
	if err != nil {
		l.Errorf(
			"claim order failed, orderID=%d riderID=%d err=%v",
			req.OrderID,
			userID,
			err,
		)
		return err
	}

	if !claimed {
		return errs.ErrOrderHaveGrabbed
	}

	defer func() {
		_ = l.svcCtx.Repo.Order.ReleaseGrabOrderClaim(
			l.ctx,
			req.OrderID,
			userID,
		)
	}()

	err = l.svcCtx.Repo.WithTx(
		l.ctx,
		func(tx *repo.RepoTx) error {

			err := tx.Order.GrabOrder(
				l.ctx,
				req.OrderID,
				userID,
				acceptedAt,
			)
			if err != nil {
				return err
			}

			return tx.Order.CreateLog(
				l.ctx,
				&model.OrderLog{
					OrderID:      req.OrderID,
					FromStatus:   enums.OrderPendingGrab,
					ToStatus:     enums.OrderAccepted,
					OperatorType: enums.OperatorTypeRider,
					OperatorID:   userID,
					Remark:       "rider grab order",
					CreatedAt:    acceptedAt,
				},
			)
		},
	)

	if err != nil {
		l.Errorf(
			"grab order failed, orderID=%d riderID=%d err=%v",
			req.OrderID,
			userID,
			err,
		)
		return err
	}

	_ = mqs.PublishDelayRiderCheck(
		l.svcCtx,
		req.OrderID,
		userID,
	)

	return nil
}
