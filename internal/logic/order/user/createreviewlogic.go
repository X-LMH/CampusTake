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

	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(
		l.ctx,
		req.OrderID,
		userID,
	)
	if err != nil {
		l.Errorf("查询评价订单失败，orderID=%d，userID=%d，err=%v", req.OrderID, userID, err)
		return err
	}

	// =========================
	// 修改：已送达/已完成都允许评价
	// =========================
	if order.Status != enums.OrderDelivered &&
		order.Status != enums.OrderCompleted {
		l.Errorf("订单状态不允许评价，orderID=%d，userID=%d，status=%v", req.OrderID, userID, order.Status)
		return errs.ErrOrderNotDelivered
	}

	if order.RiderID == nil {
		l.Errorf("订单缺少骑手信息，无法评价，orderID=%d，userID=%d", req.OrderID, userID)
		return errs.ErrOrderNoRider
	}

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		// =========================
		// 新增：检查是否重复评价
		// =========================
		exist, err := tx.Review.ExistByOrderID(
			l.ctx,
			req.OrderID,
		)
		if err != nil {
			l.Errorf("查询评价是否已存在失败，orderID=%d，userID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		if exist {
			l.Errorf("订单已评价，无法重复评价，orderID=%d，userID=%d", req.OrderID, userID)
			return errs.ErrReviewAlreadyExists
		}

		// =========================
		// 创建评价
		// =========================
		review := &model.Review{
			OrderID: req.OrderID,
			UserID:  userID,
			RiderID: *order.RiderID,
			Score:   req.Score,
			Content: req.Content,
		}

		err = tx.Review.Create(l.ctx, review)
		if err != nil {
			l.Errorf("创建评价失败，orderID=%d，userID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// =========================
		// 修改：增量更新评分
		// =========================
		rider, err := tx.Rider.GetProfileByRiderID(
			l.ctx,
			*order.RiderID,
		)
		if err != nil {
			l.Errorf("查询骑手评分信息失败，orderID=%d，riderID=%d，err=%v", req.OrderID, *order.RiderID, err)
			return err
		}

		newCount := rider.RatingCount + 1

		newAvg :=
			(rider.RatingAvg*float64(rider.RatingCount) +
				float64(req.Score)) / float64(newCount)

		return tx.Rider.UpdateRatingByRiderID(
			l.ctx,
			*order.RiderID,
			int64(newCount),
			newAvg,
		)
	})
}
