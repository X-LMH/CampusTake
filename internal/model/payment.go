package model

import (
	"CampusTake/internal/enums"
	"time"
)

// Payment 支付表模型
type Payment struct {
	ID           int64               `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID      int64               `gorm:"column:order_id;not null"`
	PayNo        string              `gorm:"column:pay_no;size:64;unique;not null"`
	Amount       float64             `gorm:"column:amount;type:decimal(10,2);not null"`
	RefundAmount float64             `gorm:"column:refund_amount;type:decimal(10,2);default:0;not null"`
	Status       enums.PaymentStatus `gorm:"column:status;not null"`
	Method       enums.PaymentMethod `gorm:"column:method;default:1"`
	PaidAt       *time.Time          `gorm:"column:paid_at"`
	CreatedAt    time.Time           `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time           `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time          `gorm:"column:deleted_at"`
}

func (Payment) TableName() string {
	return "payment"
}

// PaymentRefund 退款表模型
type PaymentRefund struct {
	ID         int64              `gorm:"column:id;primaryKey;autoIncrement"`
	PaymentID  int64              `gorm:"column:payment_id;not null"`
	RefundNo   string             `gorm:"column:refund_no;size:64;unique;not null"`
	Amount     float64            `gorm:"column:amount;type:decimal(10,2);not null"`
	Status     enums.RefundStatus `gorm:"column:status;not null"`
	Reason     string             `gorm:"column:reason;size:255"`
	OperatorID int64              `gorm:"column:operator_id"`
	RefundedAt *time.Time         `gorm:"column:refunded_at"`
	CreatedAt  time.Time          `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time          `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  *time.Time         `gorm:"column:deleted_at"`
}

func (PaymentRefund) TableName() string {
	return "payment_refund"
}
