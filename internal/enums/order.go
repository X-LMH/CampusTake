package enums

// ======================
// 订单状态
// ======================

type OrderStatus int8

const (
	OrderDefault OrderStatus = 0 // 默认

	OrderPendingPay  OrderStatus = 1 // 待支付
	OrderPendingGrab OrderStatus = 2 // 待接单

	OrderAccepted  OrderStatus = 3 // 已接单
	OrderPickedUp  OrderStatus = 4 // 已取件
	OrderDelivered OrderStatus = 5 // 已送达

	OrderUserCancelled OrderStatus = 6 // 用户取消
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

	case OrderUserCancelled:
		return "用户取消"

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
		OrderUserCancelled:
		return true

	default:
		return false
	}
}

// ======================
// 状态对应时间字段
// ======================

var OrderStatusTimeFieldMap = map[OrderStatus]string{
	OrderPendingGrab:   "paid_at",
	OrderAccepted:      "accepted_at",
	OrderPickedUp:      "picked_up_at",
	OrderDelivered:     "delivered_at",
	OrderUserCancelled: "cancelled_at",
}

func GetOrderStatusTimeField(status OrderStatus) string {
	return OrderStatusTimeFieldMap[status]
}

// ======================
// 状态流转
// ======================

var validStatusFlow = map[OrderStatus][]OrderStatus{

	// 待支付 -> 待接单 / 用户取消
	OrderPendingPay: {
		OrderPendingGrab,
		OrderUserCancelled,
	},

	// 待接单 -> 已接单 / 用户取消
	OrderPendingGrab: {
		OrderAccepted,
		OrderUserCancelled,
	},

	// 已接单 -> 已取件 / 骑手取消
	OrderAccepted: {
		OrderPendingGrab,
		OrderPickedUp,
	},

	OrderPickedUp: {
		OrderDelivered,
	},

	// 已送达 -> 无
	OrderDelivered: {},

	// 用户取消 -> 无
	OrderUserCancelled: {},
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
	OperatorTypeUser  OperatorType = 1 // 用户
	OperatorTypeRider OperatorType = 2 // 骑手
	OperatorTypeAdmin OperatorType = 3 // 管理员
)

// String 输出操作人类型中文描述
func (t OperatorType) String() string {
	switch t {
	case OperatorTypeUser:
		return "用户"
	case OperatorTypeRider:
		return "骑手"
	case OperatorTypeAdmin:
		return "管理员"
	default:
		return "未知操作人"
	}
}
