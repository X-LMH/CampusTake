package model

import (
	"CampusTake/internal/enums"
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID     int64             `gorm:"primaryKey"`
	UserID int64             `gorm:"index;not null"`
	Type   enums.AddressType `gorm:"not null"` // 1收货 2取件

	ContactName  string `gorm:"size:50;not null"`
	ContactPhone string `gorm:"size:20;not null"`

	Building string `gorm:"size:100;not null"`
	Room     string `gorm:"size:50"`
	Detail   string `gorm:"size:255"`

	IsDefault enums.AddressDefaultType `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Address) TableName() string {
	return "address"
}
