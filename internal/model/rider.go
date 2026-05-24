package model

import (
	"CampusTake/internal/enums"
	"time"

	"gorm.io/gorm"
)

type RiderProfile struct {
	ID                  int64                  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID              int64                  `gorm:"column:user_id;uniqueIndex;not null"`
	RealName            string                 `gorm:"column:real_name;type:varchar(50);not null"`
	StudentNo           string                 `gorm:"column:student_no;type:varchar(50);not null"`
	IDCardNo            string                 `gorm:"column:id_card_no;type:varchar(30);not null"`
	DormitoryBuilding   string                 `gorm:"column:dormitory_building;type:varchar(100);not null"`
	DormitoryRoom       string                 `gorm:"column:dormitory_room;type:varchar(50);not null"`
	CampusCardFront     string                 `gorm:"column:campus_card_front;type:varchar(255);not null"`
	CampusCardBack      string                 `gorm:"column:campus_card_back;type:varchar(255);not null"`
	AuditStatus         enums.RiderAuditStatus `gorm:"column:audit_status;type:tinyint;default:1;index"`
	AuditRemark         string                 `gorm:"column:audit_remark;type:varchar(255)"`
	RatingAvg           float64                `gorm:"column:rating_avg;type:decimal(3,2);default:3.00"`
	RatingCount         int32                  `gorm:"column:rating_count;type:int;default:0"`
	AcceptedOrderCount  int32                  `gorm:"column:accepted_order_count;type:int;default:0"`                // 接取订单数
	CompletedOrderCount int32                  `gorm:"column:completed_order_count;type:int;default:0"`               // 完成订单数
	CompletionRate      float64                `gorm:"column:completion_rate;type:decimal(5,2);default:0.00;->;<-:-"` // 关键：只读，不写入
	CreatedAt           time.Time              `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time              `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt           gorm.DeletedAt         `gorm:"column:deleted_at;index"`
}

// TableName 指定表名
func (RiderProfile) TableName() string {
	return "rider_profile"
}

type RiderAuditLog struct {
	ID        int64                  `gorm:"primaryKey"`
	RiderID   int64                  `gorm:"index;not null"`
	AuditorID int64                  `gorm:"index;not null"`
	Result    enums.AdminAuditResult `gorm:"not null"`
	Remark    string                 `gorm:"size:255"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RiderAuditLog) TableName() string {
	return "rider_audit_log"
}
