package model

import (
	"CampusTake/internal/enums"
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID      int64               `gorm:"primaryKey"`
	OrderID int64               `gorm:"index;not null"`
	PayNo   string              `gorm:"size:64;uniqueIndex;not null"`
	Amount  float64             `gorm:"type:decimal(10,2);not null"`
	Status  enums.PayStatus     `gorm:"index;not null"`
	Method  enums.PaymentMethod `gorm:"default:1"`

	PaidAt     *time.Time
	RefundedAt *time.Time

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Payment) TableName() string {
	return "payment"
}
