package enums

type PaymentMethod int8

const (
	PaymentMethodVirtual PaymentMethod = 1 // 虚拟支付
)

type PayStatus int8

const (
	PayStatusPending  PayStatus = 0 // 待支付
	PayStatusPaid     PayStatus = 1 // 已支付
	PayStatusRefunded PayStatus = 2 // 已退款
)
