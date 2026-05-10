package enums

type PaymentMethod int8

const (
	PaymentMethodVirtual PaymentMethod = 1 // 虚拟支付
)

type PaymentStatus int8

const (
	PaymentStatusUnpaid   PaymentStatus = 0 // 待支付
	PaymentStatusPaid     PaymentStatus = 1 // 已支付
	PaymentStatusRefunded PaymentStatus = 2 // 已退款
)

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusUnpaid:
		return "未支付"
	case PaymentStatusPaid:
		return "已支付"
	case PaymentStatusRefunded:
		return "已退款"
	default:
		return "未知状态"
	}
}
