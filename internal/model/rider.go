package model

import (
	"CampusTake/common/enum"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RiderProfile struct {
	ID                int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID            int64  `gorm:"column:user_id;uniqueIndex:uk_user_id;not null" json:"userId"`
	RealName          string `gorm:"column:real_name;type:varchar(50);not null" json:"realName"`
	StudentNo         string `gorm:"column:student_no;type:varchar(50);not null" json:"studentNo"`
	IDCardNo          string `gorm:"column:id_card_no;type:varchar(30);not null" json:"idCardNo"`
	DormitoryBuilding string `gorm:"column:dormitory_building;type:varchar(100);not null" json:"dormitoryBuilding"`
	DormitoryRoom     string `gorm:"column:dormitory_room;type:varchar(50);not null" json:"dormitoryRoom"`

	// 图片路径，存储相对路径
	CampusCardFront string `gorm:"column:campus_card_front;type:varchar(255);not null" json:"campusCardFront"`
	CampusCardBack  string `gorm:"column:campus_card_back;type:varchar(255);not null" json:"campusCardBack"`

	AuditStatus enum.RiderAuditStatus `gorm:"column:audit_status;type:tinyint;default:0;index:idx_audit_status" json:"auditStatus"`
	AuditRemark string                `gorm:"column:audit_remark;type:varchar(255)" json:"auditRemark"`

	// 评分相关字段
	RatingAvg      decimal.Decimal `gorm:"column:rating_avg;type:decimal(3,2);default:3.00" json:"ratingAvg"`
	RatingCount    int32           `gorm:"column:rating_count;type:int;default:0" json:"ratingCount"`
	CompletionRate decimal.Decimal `gorm:"column:completion_rate;type:decimal(5,2);default:0.00" json:"completionRate"`
	PunctualRate   decimal.Decimal `gorm:"column:punctual_rate;type:decimal(5,2);default:0.00" json:"punctualRate"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index:idx_deleted" json:"-"` // 软删除，json不展示
}

// TableName 指定表名
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
