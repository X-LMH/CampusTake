package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID      int64  `gorm:"primaryKey"`
	OrderNo string `gorm:"size:64;uniqueIndex;not null"`

	UserID  int64  `gorm:"index;not null"`
	RiderID *int64 `gorm:"index"`

	OrderType int8 `gorm:"not null"`

	PickupAddressID   int64 `gorm:"not null"`
	DeliveryAddressID int64 `gorm:"not null"`

	RewardAmount float64 `gorm:"type:decimal(10,2);not null"`

	Status        int8 `gorm:"index;not null"`
	PaymentStatus int8 `gorm:"default:0"`

	Remark       string `gorm:"size:255"`
	CancelReason string `gorm:"size:255"`

	PaidAt     *time.Time
	AcceptedAt *time.Time
	FinishedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Order) TableName() string {
	return "order"
}

type OrderLog struct {
	ID           int64 `gorm:"primaryKey"`
	OrderID      int64 `gorm:"index;not null"`
	FromStatus   int8
	ToStatus     int8
	OperatorType int8
	OperatorID   int64
	Remark       string `gorm:"size:255"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (OrderLog) TableName() string {
	return "order_log"
}
