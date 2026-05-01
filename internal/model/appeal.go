package model

import (
	"time"

	"gorm.io/gorm"
)

type Appeal struct {
	ID          int64 `gorm:"primaryKey"`
	OrderID     int64 `gorm:"index;not null"`
	ApplicantID int64 `gorm:"not null"`
	Type        int8
	Content     string `gorm:"size:255"`
	Status      int8   `gorm:"default:0;index"`

	HandledBy int64
	HandledAt *time.Time

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Appeal) TableName() string {
	return "appeal"
}
