package model

import (
	"CampusTake/internal/enums"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID                int64                    `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID"`
	OrderNo           string                   `gorm:"column:order_no;size:64;uniqueIndex;not null;comment:订单号"`
	UserID            int64                    `gorm:"column:user_id;index;not null;comment:下单用户ID"`
	RiderID           *int64                   `gorm:"column:rider_id;index;comment:骑手ID"`
	OrderType         enums.OrderType          `gorm:"column:order_type;not null;comment:订单类型：1=快递代取 2=外卖代取"`
	PickupAddressID   int64                    `gorm:"column:pickup_address_id;not null;comment:取件地址ID"`
	DeliveryAddressID int64                    `gorm:"column:delivery_address_id;not null;comment:送达地址ID"`
	RewardAmount      float64                  `gorm:"column:reward_amount;type:decimal(10,2);not null;comment:悬赏金额"`
	Status            enums.OrderStatus        `gorm:"column:status;index;not null;comment:订单状态"`
	PaymentStatus     enums.OrderPaymentStatus `gorm:"column:payment_status;default:0;not null;comment:支付状态：0=未支付 1=已支付 2=已退款"`
	AppealStatus      enums.OrderAppealStatus  `gorm:"column:appeal_status;default:0;not null;comment:0 无申诉 1 申诉中 2 申诉通过 3 申诉驳回"`
	CanReassign       enums.OrderCanReassign   `gorm:"column:can_reassign;default:0;not null;comment:骑手取消后是否可重新派单"`
	Remark            string                   `gorm:"column:remark;size:255;comment:订单备注"`
	CancelReason      string                   `gorm:"column:cancel_reason;size:255;comment:取消原因"`
	CreatedAt         time.Time                `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt         time.Time                `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	DeletedAt         gorm.DeletedAt           `gorm:"column:deleted_at;index;comment:软删除时间"`
	PaidAt            *time.Time               `gorm:"column:paid_at;comment:支付时间"`
	AcceptedAt        *time.Time               `gorm:"column:accepted_at;comment:骑手接单时间"`
	PickedUpAt        *time.Time               `gorm:"column:picked_up_at;comment:骑手取件时间"`
	DeliveredAt       *time.Time               `gorm:"column:delivered_at;comment:骑手送达时间"`
	CompletedAt       *time.Time               `gorm:"column:completed_at;comment:订单完成时间"`
	CancelledAt       *time.Time               `gorm:"column:cancelled_at;comment:取消时间"`
	RefundedAt        *time.Time               `gorm:"column:refunded_at;comment:退款时间"`
}

func (Order) TableName() string {
	return "order"
}

type OrderLog struct {
	ID           int64              `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID      int64              `gorm:"column:order_id;index;not null"`
	FromStatus   enums.OrderStatus  `gorm:"column:from_status"`
	ToStatus     enums.OrderStatus  `gorm:"column:to_status"`
	OperatorType enums.OperatorType `gorm:"column:operator_type"`
	OperatorID   int64              `gorm:"column:operator_id"`
	Remark       string             `gorm:"column:remark;size:255"`
	CreatedAt    time.Time          `gorm:"column:created_at;autoCreateTime"`
	DeletedAt    gorm.DeletedAt     `gorm:"column:deleted_at;index"`
}

func (OrderLog) TableName() string {
	return "order_log"
}
