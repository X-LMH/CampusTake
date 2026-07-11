package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/mqs"
	"CampusTake/internal/repo"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"CampusTake/pkg/snowflake"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

// -----------------------------
// CreateOrderLogic
// -----------------------------
type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (*types.CreateOrderResponse, error) {
	userID := ctxx.MustUserID(l.ctx)

	ok, err := l.svcCtx.Repo.Order.AllowCreateOrderLimit(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.ErrRateLimitError
	}

	// -----------------------------
	// 1. 参数校验
	// -----------------------------
	if req.PickupAddressID == req.DeliveryAddressID {
		return nil, errors.NewParamError("取件地址和收货地址不能相同")
	}

	// -----------------------------
	// 2. 地址校验（放在事务外，避免占用事务时间）
	// -----------------------------
	pickupAddress, err := l.svcCtx.Repo.Address.GetByIDAndUserID(l.ctx, req.PickupAddressID, userID)
	if err != nil {
		return nil, err
	}
	if pickupAddress.Type != enums.AddressTypePickup {
		return nil, errors.ErrAddressTypeInvalid
	}

	deliveryAddress, err := l.svcCtx.Repo.Address.GetByIDAndUserID(l.ctx, req.DeliveryAddressID, userID)
	if err != nil {
		return nil, err
	}
	if deliveryAddress.Type != enums.AddressTypeDelivery {
		return nil, errors.ErrAddressTypeInvalid
	}

	var order *model.Order

	// -----------------------------
	// 3. 事务：订单 + 支付
	// -----------------------------
	err = l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		orderNo := snowflake.GenerateID()

		order = &model.Order{
			OrderNo:           orderNo,
			UserID:            userID,
			OrderType:         enums.OrderType(req.OrderType),
			PickupAddressID:   req.PickupAddressID,
			DeliveryAddressID: req.DeliveryAddressID,
			RewardAmount:      req.RewardAmount,
			Status:            enums.OrderPendingPay,
			PaymentStatus:     enums.OrderPayStatusUnpaid,
			Remark:            req.Remark,
		}

		if err := tx.Order.Create(l.ctx, order); err != nil {
			return err
		}

		payNo := snowflake.GenerateID()

		payment := &model.Payment{
			OrderID: order.ID,
			PayNo:   payNo,
			Amount:  order.RewardAmount,
			Status:  enums.PaymentStatusWaitPay,
			Method:  enums.PaymentMethodMock,
		}

		if err := tx.Payment.Create(l.ctx, payment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 4. 事务成功后：异步副作用（重点）
	// -----------------------------

	// 4.1 BloomFilter（必须放这里）
	l.svcCtx.Repo.Order.AddToBloom(l.ctx, order.ID)

	// 4.2 延迟取消订单 MQ
	if err := mqs.PublishDelayCancelOrder(
		l.svcCtx,
		order.ID,
	); err != nil {
		l.Errorf("发送延迟取消订单消息失败，orderID=%d err=%v", order.ID, err)
	}

	// -----------------------------
	// 5. 返回结果
	// -----------------------------
	return &types.CreateOrderResponse{
		OrderID: order.ID,
		OrderNo: order.OrderNo,
	}, nil
}
