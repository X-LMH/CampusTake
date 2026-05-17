package enums

// ======================
// 订单状态
// ======================

type OrderStatus int8

const (

	// ======================
	// 默认状态
	// ======================

	OrderDefault OrderStatus = 0 // 默认

	// ======================
	// 主流程状态
	// ======================

	// 1 待支付
	OrderPendingPay OrderStatus = 1

	// 2 待接单
	OrderPendingGrab OrderStatus = 2

	// 3 已接单
	OrderAccepted OrderStatus = 3

	// 4 配送中
	OrderDelivering OrderStatus = 4

	// 5 已送达（骑手确认）
	OrderDelivered OrderStatus = 5

	// 6 已完成（用户确认 / 系统自动完成）
	OrderCompleted OrderStatus = 6

	// ======================
	// 关闭 / 异常状态
	// ======================

	// 7 用户取消
	OrderUserCancelled OrderStatus = 7

	// 8 骑手取消
	OrderRiderCancelled OrderStatus = 8

	// 9 超时关闭
	OrderTimeoutClosed OrderStatus = 9

	// 10 已退款
	OrderRefunded OrderStatus = 10

	// 11 异常订单
	OrderException OrderStatus = 11

	// 12 申诉中
	OrderAppealing OrderStatus = 12
)

func (s OrderStatus) String() string {

	switch s {

	case OrderPendingPay:
		return "待支付"

	case OrderPendingGrab:
		return "待接单"

	case OrderAccepted:
		return "已接单"

	case OrderDelivering:
		return "配送中"

	case OrderDelivered:
		return "已送达"

	case OrderCompleted:
		return "已完成"

	case OrderUserCancelled:
		return "用户取消"

	case OrderRiderCancelled:
		return "骑手取消"

	case OrderTimeoutClosed:
		return "超时关闭"

	case OrderRefunded:
		return "已退款"

	case OrderException:
		return "异常订单"

	case OrderAppealing:
		return "申诉中"

	default:
		return "未知状态"
	}
}

func IsOrderStatus(s OrderStatus) bool {

	switch s {

	case OrderPendingPay,
		OrderPendingGrab,
		OrderAccepted,
		OrderDelivering,
		OrderDelivered,
		OrderCompleted,
		OrderUserCancelled,
		OrderRiderCancelled,
		OrderTimeoutClosed,
		OrderRefunded,
		OrderException,
		OrderAppealing:
		return true

	default:
		return false
	}
}

// ======================
// 是否终态
// ======================

func (s OrderStatus) IsFinalStatus() bool {

	switch s {

	case OrderCompleted,
		OrderRefunded:
		return true

	default:
		return false
	}
}

// ======================
// 状态对应时间字段
// ======================

var OrderStatusTimeFieldMap = map[OrderStatus]string{

	OrderPendingGrab: "paid_at",

	OrderAccepted: "accepted_at",

	OrderDelivered: "delivered_at",

	OrderCompleted: "completed_at",

	OrderUserCancelled: "cancelled_at",

	OrderRiderCancelled: "cancelled_at",

	OrderTimeoutClosed: "cancelled_at",

	OrderRefunded: "refunded_at",
}

func GetOrderStatusTimeField(status OrderStatus) string {
	return OrderStatusTimeFieldMap[status]
}

// ======================
// 状态流转
// ======================

var validStatusFlow = map[OrderStatus][]OrderStatus{

	// 待支付
	OrderPendingPay: {
		OrderPendingGrab,
		OrderUserCancelled,
		OrderTimeoutClosed,
	},

	// 待接单
	OrderPendingGrab: {
		OrderAccepted,
		OrderUserCancelled,
		OrderTimeoutClosed,
	},

	// 已接单
	OrderAccepted: {
		OrderDelivering,
		OrderRiderCancelled,
		OrderException,
		OrderAppealing,
	},

	// 配送中
	OrderDelivering: {
		OrderDelivered,
		OrderException,
		OrderAppealing,
	},

	// 已送达
	OrderDelivered: {
		OrderCompleted,
		OrderAppealing,
	},

	// 已完成
	OrderCompleted: {},

	// 用户取消
	OrderUserCancelled: {
		OrderRefunded,
	},

	// 骑手取消
	OrderRiderCancelled: {
		OrderPendingGrab,
		OrderRefunded,
	},

	// 超时关闭
	OrderTimeoutClosed: {
		OrderRefunded,
	},

	// 已退款
	OrderRefunded: {},

	// 异常订单
	OrderException: {
		OrderRefunded,
		OrderCompleted,
		OrderAppealing,
	},

	// 申诉中
	OrderAppealing: {
		OrderCompleted,
		OrderRefunded,
		OrderException,
	},
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

// String 输出订单类型中文描述
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
	OperatorTypeUser   OperatorType = 1 // 用户
	OperatorTypeRider  OperatorType = 2 // 骑手
	OperatorTypeSystem OperatorType = 3 // 系统
	OperatorTypeAdmin  OperatorType = 4 // 管理员
)

// String 输出操作人类型中文描述
func (t OperatorType) String() string {
	switch t {
	case OperatorTypeUser:
		return "用户"
	case OperatorTypeRider:
		return "骑手"
	case OperatorTypeSystem:
		return "系统"
	case OperatorTypeAdmin:
		return "管理员"
	default:
		return "未知操作人"
	}
}
