package types

import (
	"CampusTake/internal/model"
	"CampusTake/pkg/utils"
)

// ==========================================
// Order 转换模块
// ==========================================

func ModelOrderToOrderItem(order *model.Order) *OrderItem {
	if order == nil {
		return nil
	}
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

// ==========================================
// Appeal 转换模块
// ==========================================

func ModelAppealToAppealItem(appeal *model.Appeal) *AppealItem {
	if appeal == nil {
		return nil
	}
	return &AppealItem{
		ID:                appeal.ID,
		OrderID:           appeal.OrderID,
		ApplicantID:       appeal.ApplicantID,
		ApplicantRole:     int8(appeal.ApplicantRole),
		ApplicantRoleText: appeal.ApplicantRole.String(),
		AppealType:        appeal.AppealType,
		AppealTypeText:    getAppealTypeText(appeal.AppealType),
		Content:           appeal.Content,
		EvidenceUrls:      utils.ParseJSONToStringSlice(appeal.EvidenceUrls),
		Status:            int8(appeal.Status),
		StatusText:        appeal.Status.String(),
		CreatedAt:         utils.FormatCreatedAt(appeal.CreatedAt),
		HandledAt:         utils.FormatTimePtr(appeal.HandledAt),
	}
}

func getAppealTypeText(appealType int8) string {
	types := map[int8]string{
		1: "未收到货",
		2: "骑手态度恶劣",
		3: "骑手提前点送达",
		4: "用户恶意退款",
		5: "用户恶意投诉",
		6: "订单丢失",
		7: "物品损坏",
	}
	if text, ok := types[appealType]; ok {
		return text
	}
	return "其他" // 将 default 优雅归类为其他或未知
}

// ==========================================
// AppealHandle 转换模块
// ==========================================

// ModelAppealHandleToItem 命名对齐，增加防空
func ModelAppealHandleToItem(handle *model.AppealHandle) *AppealHandleItem {
	if handle == nil {
		return nil
	}
	return &AppealHandleItem{
		ID:             handle.ID,
		Result:         int8(handle.Result),
		ResultText:     handle.Result.String(),
		RefundAmount:   handle.RefundAmount,
		PunishRider:    handle.PunishRider,
		PunishUser:     handle.PunishUser,
		TerminateOrder: handle.TerminateOrder,
		Remark:         handle.Remark,
		HandlerID:      handle.HandlerID,
		CreatedAt:      utils.FormatCreatedAt(handle.CreatedAt),
	}
}
