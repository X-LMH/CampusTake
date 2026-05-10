package types

import (
	"CampusTake/internal/model"
	"CampusTake/pkg/utils"
)

func ModelOrderToOrderItem(order *model.Order) OrderItem {
	return OrderItem{
		ID:                order.ID,
		OrderNo:           order.OrderNo,
		UserID:            order.UserID,
		RiderID:           order.RiderID,
		OrderType:         int8(order.OrderType),
		PickupAddressID:   order.PickupAddressID,
		DeliveryAddressID: order.DeliveryAddressID,
		RewardAmount:      order.RewardAmount,
		Status:            int8(order.Status),
		StatusText:        order.Status.String(),
		PaymentStatus:     int8(order.PaymentStatus),
		PaymentStatusText: order.PaymentStatus.String(),
		Remark:            order.Remark,
		CancelReason:      order.CancelReason,
		CreatedAt:         utils.FormatCreatedAt(order.CreatedAt),
		PaidAt:            utils.FormatTimePtr(order.PaidAt),
		AcceptedAt:        utils.FormatTimePtr(order.AcceptedAt),
		PickedUpAt:        utils.FormatTimePtr(order.PickedUpAt),
		DeliveredAt:       utils.FormatTimePtr(order.DeliveredAt),
		CancelledAt:       utils.FormatTimePtr(order.CancelledAt),
		RefundedAt:        utils.FormatTimePtr(order.RefundedAt),
	}
}

func ModelOrderToOrderDetail(order *model.Order) *OrderItem {
	return &OrderItem{
		ID:                order.ID,
		OrderNo:           order.OrderNo,
		UserID:            order.UserID,
		RiderID:           order.RiderID,
		OrderType:         int8(order.OrderType),
		PickupAddressID:   order.PickupAddressID,
		DeliveryAddressID: order.DeliveryAddressID,
		RewardAmount:      order.RewardAmount,
		Status:            int8(order.Status),
		StatusText:        order.Status.String(),
		PaymentStatus:     int8(order.PaymentStatus),
		PaymentStatusText: order.PaymentStatus.String(),
		Remark:            order.Remark,
		CancelReason:      order.CancelReason,
		CreatedAt:         utils.FormatCreatedAt(order.CreatedAt),
		PaidAt:            utils.FormatTimePtr(order.PaidAt),
		AcceptedAt:        utils.FormatTimePtr(order.AcceptedAt),
		PickedUpAt:        utils.FormatTimePtr(order.PickedUpAt),
		DeliveredAt:       utils.FormatTimePtr(order.DeliveredAt),
		CancelledAt:       utils.FormatTimePtr(order.CancelledAt),
		RefundedAt:        utils.FormatTimePtr(order.RefundedAt),
	}
}
