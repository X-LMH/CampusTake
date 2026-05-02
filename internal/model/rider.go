package model

import (
	"time"

	"gorm.io/gorm"
)

type RiderProfile struct {
	ID                int64  `gorm:"primaryKey"`
	UserID            int64  `gorm:"uniqueIndex;not null"`
	RealName          string `gorm:"size:50;not null"`
	StudentNo         string `gorm:"size:50;not null"`
	IDCardNo          string `gorm:"size:30;not null"`
	DormitoryBuilding string `gorm:"size:100;not null"`
	DormitoryRoom     string `gorm:"size:50;not null"`
	CampusCardPhoto   string `gorm:"size:255"`

	AuditStatus int8   `gorm:"default:0;index"`
	AuditRemark string `gorm:"size:255"`

	RatingAvg      float64 `gorm:"type:decimal(3,2);default:3.00"`
	RatingCount    int     `gorm:"default:0"`
	CompletionRate float64 `gorm:"type:decimal(5,2);default:0"`
	PunctualRate   float64 `gorm:"type:decimal(5,2);default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RiderProfile) TableName() string {
	return "rider_profile"
}

type RiderAuditLog struct {
	ID        int64  `gorm:"primaryKey"`
	RiderID   int64  `gorm:"index;not null"`
	AuditorID int64  `gorm:"index;not null"`
	Result    int8   `gorm:"not null"`
	Remark    string `gorm:"size:255"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RiderAuditLog) TableName() string {
	return "rider_audit_log"
}
