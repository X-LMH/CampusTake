package model

import (
	"CampusTake/internal/enums"
	"encoding/json"
	"time"
)

type Appeal struct {
	ID            int64                     `gorm:"primaryKey;comment:申诉ID"`
	OrderID       int64                     `gorm:"not null;index;comment:订单ID"`
	ApplicantID   int64                     `gorm:"not null;comment:申诉人ID"`
	ApplicantRole enums.AppealApplicantRole `gorm:"not null;comment:1用户 2骑手"`
	AppealType    int8                      `gorm:"not null;comment:申诉类型"`
	Content       string                    `gorm:"type:varchar(500);not null;comment:申诉内容"`
	EvidenceUrls  json.RawMessage           `gorm:"type:json;comment:证据图片"`
	Status        enums.AppealStatus        `gorm:"not null;default:1;comment:申诉状态：1待处理 2已通过 3已驳回 4已撤销"`
	HandledBy     *int64                    `gorm:"comment:处理管理员ID"`
	HandledAt     *time.Time                `gorm:"comment:处理时间"`
	CreatedAt     time.Time                 `gorm:"autoCreateTime"`
	UpdatedAt     time.Time                 `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time                `gorm:"index"`
}

func (Appeal) TableName() string {
	return "appeal"
}

type AppealHandle struct {
	ID             int64                    `gorm:"primaryKey;comment:处理记录ID"`
	AppealID       int64                    `gorm:"not null;comment:申诉ID"`
	HandlerID      int64                    `gorm:"not null;comment:管理员ID"`
	Result         enums.AppealHandleResult `gorm:"not null;comment:1通过 2驳回"`
	Remark         string                   `gorm:"type:varchar(255);comment:处理备注"`
	RefundAmount   float64                  `gorm:"type:decimal(10,2);default:0.00;comment:退款金额"`
	PunishRider    int8                     `gorm:"default:0;comment:是否处罚骑手"`
	PunishUser     int8                     `gorm:"default:0;comment:是否处罚用户"`
	TerminateOrder int8                     `gorm:"default:0;comment:是否终止订单"`
	CreatedAt      time.Time                `gorm:"autoCreateTime;comment:创建时间"`
}

func (AppealHandle) TableName() string {
	return "appeal_handle"
}
