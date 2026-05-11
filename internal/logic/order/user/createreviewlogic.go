// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReviewLogic {
	return &CreateReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateReviewLogic) CreateReview(req *types.CreateReviewRequest) error {
	userID := ctxx.MustUserID(l.ctx)
	// 验证订单是否存在且属于该用户
	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(l.ctx, req.OrderID, userID)
	if err != nil {
		return err
	}
	// 状态校验，只有已完成的订单才能评价
	if order.Status != enums.OrderDelivered {
		return errs.ErrOrderNotDelivered
	}

	if order.RiderID == nil {
		return errs.ErrOrderNoRider
	}

	err = l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		// 创建评价
		review := &model.Review{
			OrderID: req.OrderID,
			UserID:  userID,
			RiderID: *order.RiderID,
			Score:   req.Score,
			Content: req.Content,
		}
		err = tx.Review.Create(l.ctx, review)
		if err != nil {
			return err
		}

		// 计算评价总数
		reviewCount, avg, err := tx.Review.GetReviewCountAndScoreByRiderID(
			l.ctx,
			*order.RiderID,
		)
		if err != nil {
			return err
		}

		err = tx.Rider.UpdateRatingByRiderID(
			l.ctx,
			*order.RiderID,
			reviewCount,
			avg,
		)
		if err != nil {
			return err
		}

		// =========================
		// 完成单数
		// =========================
		completeCount, err := tx.Order.GetCountByRiderIDAndStatuses(
			l.ctx,
			*order.RiderID,
			[]enums.OrderStatus{
				enums.OrderDelivered,
			},
		)
		if err != nil {
			return err
		}

		// =========================
		// 接单总数
		// =========================
		totalCount, err := tx.Order.GetCountByRiderIDAndStatuses(
			l.ctx,
			*order.RiderID,
			[]enums.OrderStatus{
				enums.OrderAccepted,
				enums.OrderPickedUp,
				enums.OrderDelivered,
			},
		)
		if err != nil {
			return err
		}

		// =========================
		// 计算完成率
		// =========================
		completionRate := 0.0

		if totalCount > 0 {
			completionRate = float64(completeCount) / float64(totalCount) * 100
		}

		// =========================
		// 更新骑手统计
		// =========================
		return tx.Rider.UpdateCompleteStatsByRiderID(
			l.ctx,
			*order.RiderID,
			int(completeCount),
			completionRate,
		)
	})

	return nil
}
