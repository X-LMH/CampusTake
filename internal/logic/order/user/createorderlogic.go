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
	order := new(model.Order)

	// 1. 校验取件地址和收货地址
	if req.PickupAddressID == req.DeliveryAddressID {
		return nil, errors.NewParamError("取件地址和收货地址不能相同")
	}
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

	err = l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		// -----------------------------
		// 1. 创建订单
		// -----------------------------
		orderNo := snowflake.GenerateID()
		order = &model.Order{
			OrderNo:           orderNo,
			UserID:            userID,
			OrderType:         enums.OrderType(req.OrderType),
			PickupAddressID:   req.PickupAddressID,
			DeliveryAddressID: req.DeliveryAddressID,
			RewardAmount:      req.RewardAmount,
			Status:            enums.OrderPendingPay,      // 待支付
			PaymentStatus:     enums.OrderPayStatusUnpaid, // 未支付
			Remark:            req.Remark,
		}

		if err := tx.Order.Create(l.ctx, order); err != nil {
			return err
		}

		// -----------------------------
		// 2. 创建支付记录
		// -----------------------------
		payNo := snowflake.GenerateID()
		payment := &model.Payment{
			OrderID: order.ID,
			PayNo:   payNo,
			Amount:  order.RewardAmount,
			Status:  enums.PaymentStatusUnpaid,
			Method:  enums.PaymentMethodVirtual,
		}

		if err := tx.Payment.Create(l.ctx, payment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = mqs.PublishDelayCancelOrder(
		l.svcCtx,
		order.ID,
	)
	if err != nil {
		l.Errorf("发送延迟取消订单消息失败，orderID=%d err=%v", order.ID, err)
	}

	return &types.CreateOrderResponse{
		OrderID:           order.ID,
		OrderNo:           order.OrderNo,
		Status:            int8(order.Status),
		PaymentStatus:     int8(order.PaymentStatus),
		PaymentStatusText: order.PaymentStatus.String(),
	}, nil
}
