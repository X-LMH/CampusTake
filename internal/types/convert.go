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
		CompletedAt:       utils.FormatTimePtr(order.CompletedAt),
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
		CompletedAt:       utils.FormatTimePtr(order.CompletedAt),
		CancelledAt:       utils.FormatTimePtr(order.CancelledAt),
		RefundedAt:        utils.FormatTimePtr(order.RefundedAt),
	}
}

func ModelAppealToAppealItem(appeal model.Appeal) AppealItem {
	return AppealItem{
		ID:           appeal.ID,
		OrderID:      appeal.OrderID,
		ApplicantID:  appeal.ApplicantID,
		Type:         int8(appeal.Type),
		TypeText:     appeal.Type.String(),
		Content:      appeal.Content,
		Status:       int8(appeal.Status),
		StatusText:   appeal.Status.String(),
		RefundAmount: appeal.RefundAmount,
		PunishRider:  appeal.PunishRider,
		HandleRemark: appeal.HandleRemark,
		CreatedAt:    utils.FormatCreatedAt(appeal.CreatedAt),
		HandledAt:    utils.FormatTimePtr(appeal.HandledAt),
	}
}
