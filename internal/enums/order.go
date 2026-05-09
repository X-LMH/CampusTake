package enums

// ====================== 1. 订单类型 ======================
type OrderType int8

const (
	OrderTypeExpress  OrderType = 1 // 快递
	OrderTypeTakeAway OrderType = 2 // 外卖
)

func (t OrderType) String() string {
	switch t {
	case OrderTypeExpress:
		return "快递"
	case OrderTypeTakeAway:
		return "外卖"
	default:
		return "未知类型"
	}
}

// ====================== 2. 订单主状态 ======================
type OrderStatus int8

const (
	OrderDefault     OrderStatus = 0  // 默认状态
	OrderPendingPay  OrderStatus = 1  // 待支付
	OrderPendingGrab OrderStatus = 2  // 待接单
	OrderAccepted    OrderStatus = 3  // 已接单
	OrderPicking     OrderStatus = 4  // 取件中
	OrderDelivering  OrderStatus = 5  // 配送中
	OrderFinished    OrderStatus = 6  // 已完成
	OrderCancelled   OrderStatus = 7  // 已取消
	OrderRefunded    OrderStatus = 8  // 已退款
	OrderTimeout     OrderStatus = 9  // 超时
	OrderException   OrderStatus = 10 // 异常
	OrderAppealing   OrderStatus = 11 // 申诉中
	OrderWithdrawn   OrderStatus = 12 // 已撤单
)

// String 输出中文描述
func (s OrderStatus) String() string {
	switch s {
	case OrderPendingPay:
		return "待支付"
	case OrderPendingGrab:
		return "待接单"
	case OrderAccepted:
		return "已接单"
	case OrderPicking:
		return "取件中"
	case OrderDelivering:
		return "配送中"
	case OrderFinished:
		return "已完成"
	case OrderCancelled:
		return "已取消"
	case OrderRefunded:
		return "已退款"
	case OrderTimeout:
		return "超时"
	case OrderException:
		return "异常"
	case OrderAppealing:
		return "申诉中"
	case OrderWithdrawn:
		return "已撤单"
	default:
		return "未知状态"
	}
}

func IsOrderStatus(s OrderStatus) bool {
	switch s {
	case OrderPendingPay, OrderPendingGrab, OrderAccepted, OrderPicking, OrderDelivering,
		OrderFinished, OrderCancelled, OrderRefunded, OrderTimeout, OrderException,
		OrderAppealing, OrderWithdrawn:
		return true
	default:
		return false
	}
}

var validStatusFlow = map[OrderStatus][]OrderStatus{
	// 待支付 → 待接单 / 已取消
	OrderPendingPay: {
		OrderPendingGrab,
		OrderCancelled,
	},

	// 待接单 → 已接单 / 已取消 / 超时
	OrderPendingGrab: {
		OrderAccepted,
		OrderCancelled,
		OrderTimeout,
	},

	// 已接单 → 取件中 / 异常
	OrderAccepted: {
		OrderPicking,
		OrderException,
	},

	// 取件中 → 配送中 / 异常
	OrderPicking: {
		OrderDelivering,
		OrderException,
	},

	// 配送中 → 已完成 / 异常
	OrderDelivering: {
		OrderFinished,
		OrderException,
	},

	// 已完成 → 申诉中 / 结束
	OrderFinished: {
		OrderAppealing,
	},

	// 已取消 → 无流转
	OrderCancelled: {},

	// 已退款 → 无流转
	OrderRefunded: {},

	// 超时 → 无流转
	OrderTimeout: {},

	// 异常 → 无流转
	OrderException: {},

	// 申诉中 → 已完成 / 已撤单
	OrderAppealing: {
		OrderFinished,
		OrderWithdrawn,
	},

	// 已撤单 → 无流转
	OrderWithdrawn: {},
}

func CheckOrderStatusFlow(current, target OrderStatus) bool {
	// 获取当前状态允许的下一个状态列表
	nextList, ok := validStatusFlow[current]
	if !ok {
		return false
	}

	// 检查 target 是否在允许列表中
	for _, s := range nextList {
		if s == target {
			return true
		}
	}

	return false
}

// ====================== 3. 支付状态 ======================
type PaymentStatus int8

const (
	PaymentStatusUnpaid PaymentStatus = 0 // 未支付
	PaymentStatusPaid   PaymentStatus = 1 // 已支付
	PaymentStatusRefund PaymentStatus = 2 // 已退款
)

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusUnpaid:
		return "未支付"
	case PaymentStatusPaid:
		return "已支付"
	case PaymentStatusRefund:
		return "已退款"
	default:
		return "未知支付状态"
	}
}

type OperatorType int8

const (
	OperatorUser   OperatorType = 1 // 用户
	OperatorRider  OperatorType = 2 // 骑手
	OperatorSystem OperatorType = 3 // 系统
	OperatorAdmin  OperatorType = 4 // 管理员
)

// String 返回中文描述
func (o OperatorType) String() string {
	switch o {
	case OperatorUser:
		return "用户"
	case OperatorRider:
		return "骑手"
	case OperatorSystem:
		return "系统"
	case OperatorAdmin:
		return "管理员"
	default:
		return "未知"
	}
}
