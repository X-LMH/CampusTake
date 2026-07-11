package model

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID      int64  `gorm:"primaryKey"`
	OrderID int64  `gorm:"uniqueIndex;not null"`
	UserID  int64  `gorm:"not null"`
	RiderID int64  `gorm:"index;not null"`
	Score   int8   `gorm:"not null"`
	Content string `gorm:"size:255"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Review) TableName() string {
	return "review"
}
