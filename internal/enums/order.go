package enums

// ======================
// 订单状态
// ======================

type OrderStatus int8

const (
	OrderDefault OrderStatus = 0 // 默认

	OrderPendingPay  OrderStatus = 1 // 待支付
	OrderPendingGrab OrderStatus = 2 // 待接单
	OrderAccepted    OrderStatus = 3 // 已接单
	OrderPickedUp    OrderStatus = 4 // 已取件
	OrderDelivered   OrderStatus = 5 // 已送达

	OrderCancelled OrderStatus = 6 // 已取消
	OrderRefunded  OrderStatus = 7 // 已退款
)

func (s OrderStatus) String() string {
	switch s {

	case OrderPendingPay:
		return "待支付"

	case OrderPendingGrab:
		return "待接单"

	case OrderAccepted:
		return "已接单"

	case OrderPickedUp:
		return "已取件"

	case OrderDelivered:
		return "已送达"

	case OrderCancelled:
		return "已取消"

	case OrderRefunded:
		return "已退款"

	default:
		return "未知状态"
	}
}

func IsOrderStatus(s OrderStatus) bool {
	switch s {

	case OrderPendingPay,
		OrderPendingGrab,
		OrderAccepted,
		OrderPickedUp,
		OrderDelivered,
		OrderCancelled,
		OrderRefunded:

		return true

	default:
		return false
	}
}

var OrderStatusTimeFieldMap = map[OrderStatus]string{
	OrderPendingGrab: "paid_at",
	OrderAccepted:    "accepted_at",
	OrderPickedUp:    "picked_up_at",
	OrderDelivered:   "delivered_at",
	OrderCancelled:   "cancelled_at",
	OrderRefunded:    "refunded_at",
}

func GetOrderStatusTimeField(status OrderStatus) string {
	return OrderStatusTimeFieldMap[status]
}

// ======================
// 状态流转
// ======================

var validStatusFlow = map[OrderStatus][]OrderStatus{

	// 待支付 -> 待接单 / 已取消
	OrderPendingPay: {
		OrderPendingGrab,
		OrderCancelled,
	},

	// 待接单 -> 已接单 / 已取消
	OrderPendingGrab: {
		OrderAccepted,
		OrderCancelled,
	},

	// 已接单 -> 已取件 / 已取消
	OrderAccepted: {
		OrderPickedUp,
		OrderCancelled,
	},

	// 已取件 -> 已送达
	OrderPickedUp: {
		OrderDelivered,
	},

	// 已送达 -> 无
	OrderDelivered: {},

	// 已取消 -> 已退款
	OrderCancelled: {
		OrderRefunded,
	},

	// 已退款 -> 无
	OrderRefunded: {},
}

func CheckOrderStatusFlow(current, target OrderStatus) bool {

	nextList, ok := validStatusFlow[current]
	if !ok {
		return false
	}

	for _, s := range nextList {
		if s == target {
			return true
		}
	}

	return false
}

// ======================
// 支付状态
// ======================

type OrderPaymentStatus int8

const (
	OrderPayStatusUnpaid   OrderPaymentStatus = 0 // 未支付
	OrderPayStatusPaid     OrderPaymentStatus = 1 // 已支付
	OrderPayStatusRefunded OrderPaymentStatus = 2 // 已退款
)

func (s OrderPaymentStatus) String() string {
	switch s {
	case OrderPayStatusUnpaid:
		return "未支付"
	case OrderPayStatusPaid:
		return "已支付"
	case OrderPayStatusRefunded:
		return "已退款"
	default:
		return "未知状态"
	}
}

// ======================
// 订单类型
// ======================

type OrderType int8

const (
	OrderTypeExpress OrderType = 1 // 快递代取
	OrderTypeFood    OrderType = 2 // 外卖代取
)

func (t OrderType) String() string {
	switch t {

	case OrderTypeExpress:
		return "快递代取"

	case OrderTypeFood:
		return "外卖代取"

	default:
		return "未知类型"
	}
}

// ======================
// 操作人类型
// ======================

type OperatorType int8

const (
	OperatorUser  OperatorType = 1 // 用户
	OperatorRider OperatorType = 2 // 骑手
	OperatorAdmin OperatorType = 3 // 管理员
)

func (t OperatorType) String() string {
	switch t {

	case OperatorUser:
		return "用户"

	case OperatorRider:
		return "骑手"

	case OperatorAdmin:
		return "管理员"

	default:
		return "未知操作人"
	}
}
