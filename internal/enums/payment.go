package enums

// PaymentStatus 支付状态
type PaymentStatus int8

const (
	PaymentStatusWaitPay    PaymentStatus = 1 // 待支付
	PaymentStatusPaid       PaymentStatus = 2 // 已支付
	PaymentStatusPartRefund PaymentStatus = 3 // 部分退款
	PaymentStatusFullRefund PaymentStatus = 4 // 已全额退款
)

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusWaitPay:
		return "待支付"
	case PaymentStatusPaid:
		return "已支付"
	case PaymentStatusPartRefund:
		return "部分退款"
	case PaymentStatusFullRefund:
		return "已全额退款"
	default:
		return "未知支付状态"
	}
}

// PaymentMethod 支付方式
type PaymentMethod int8

const (
	PaymentMethodMock   PaymentMethod = 1 // 模拟支付
	PaymentMethodWechat PaymentMethod = 2 // 微信支付
)

func (m PaymentMethod) String() string {
	switch m {
	case PaymentMethodMock:
		return "模拟支付"
	case PaymentMethodWechat:
		return "微信支付"
	default:
		return "未知支付方式"
	}
}

// RefundStatus 退款状态
type RefundStatus int8

const (
	RefundStatusProcessing RefundStatus = 1 // 退款中
	RefundStatusSuccess    RefundStatus = 2 // 退款成功
	RefundStatusFailed     RefundStatus = 3 // 退款失败
)

func (s RefundStatus) String() string {
	switch s {
	case RefundStatusProcessing:
		return "退款中"
	case RefundStatusSuccess:
		return "退款成功"
	case RefundStatusFailed:
		return "退款失败"
	default:
		return "未知退款状态"
	}
}
