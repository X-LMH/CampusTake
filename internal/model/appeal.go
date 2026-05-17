package model

import (
	"CampusTake/internal/enums"
	"time"
)

type Appeal struct {
	ID int64 `gorm:"primaryKey;comment:主键ID"`

	OrderID int64 `gorm:"not null;index:idx_order_id;comment:订单ID"`

	ApplicantID int64 `gorm:"not null;index:idx_applicant_id;comment:申诉人ID"`

	// 1用户申诉 2骑手申诉
	Type enums.AppealType `gorm:"not null;comment:申诉类型：1用户申诉 2骑手申诉"`

	Content string `gorm:"type:varchar(500);not null;comment:申诉内容"`

	// 0待处理 1已通过 2已驳回
	Status enums.AppealStatus `gorm:"not null;default:0;index:idx_status;comment:申诉状态"`

	// 管理员处理备注
	HandleRemark string `gorm:"type:varchar(255);default:'';comment:处理备注"`

	// 退款金额
	RefundAmount float64 `gorm:"type:decimal(10,2);default:0;comment:退款金额"`

	// 是否处罚骑手：0否 1是
	PunishRider int8 `gorm:"default:0;comment:是否处罚骑手：0否 1是"`

	// 是否终止订单：0否 1是
	TerminateOrder int8 `gorm:"default:0;comment:是否终止订单：0否 1是"`

	HandledBy *int64 `gorm:"comment:处理管理员ID"`

	HandledAt *time.Time `gorm:"comment:处理时间"`

	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间"`

	DeletedAt *time.Time `gorm:"index:idx_deleted_at;comment:软删除时间"`
}

func (Appeal) TableName() string {
	return "appeal"
}
