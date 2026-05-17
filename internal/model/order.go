package model

import (
	"CampusTake/internal/enums"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID                int64                    `gorm:"primaryKey"`
	OrderNo           string                   `gorm:"size:64;uniqueIndex;not null"`
	UserID            int64                    `gorm:"index;not null"`
	RiderID           *int64                   `gorm:"index"`
	OrderType         enums.OrderType          `gorm:"not null"`
	PickupAddressID   int64                    `gorm:"not null"`
	DeliveryAddressID int64                    `gorm:"not null"`
	RewardAmount      float64                  `gorm:"type:decimal(10,2);not null"`
	Status            enums.OrderStatus        `gorm:"index;not null"`
	PaymentStatus     enums.OrderPaymentStatus `gorm:"default:0;not null"`
	Remark            string                   `gorm:"size:255"`
	CancelReason      string                   `gorm:"size:255"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	PaidAt            *time.Time
	AcceptedAt        *time.Time
	PickedUpAt        *time.Time
	DeliveredAt       *time.Time
	CompletedAt       *time.Time
	CancelledAt       *time.Time
	RefundedAt        *time.Time
}

func (Order) TableName() string {
	return "order"
}

type OrderLog struct {
	ID           int64 `gorm:"primaryKey"`
	OrderID      int64 `gorm:"index;not null"`
	FromStatus   enums.OrderStatus
	ToStatus     enums.OrderStatus
	OperatorType enums.OperatorType
	OperatorID   int64
	Remark       string `gorm:"size:255"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (OrderLog) TableName() string {
	return "order_log"
}
