// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/pkg/ctxx"
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

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		order, err := tx.Order.GetByID(l.ctx, req.OrderID)
		if err != nil {
			return err
		}

		err = tx.Order.GrabOrder(
			l.ctx,
			req.OrderID,
			userID,
			acceptedAt,
		)
		if err != nil {
			return err
		}

		log := &model.OrderLog{
			OrderID:      req.OrderID,
			FromStatus:   order.Status,
			ToStatus:     enums.OrderAccepted,
			OperatorType: enums.OperatorTypeRider,
			OperatorID:   userID,
			Remark:       "代取员抢单",
			CreatedAt:    acceptedAt,
		}

		return tx.Order.CreateLog(l.ctx, log)
	})
}
