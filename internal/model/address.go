package model

import (
	"CampusTake/common/enum"
	"time"

	"gorm.io/gorm"
)

// Address 用户收货地址模型
type Address struct {
	ID           uint64                  `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserID       uint64                  `gorm:"column:user_id;not null;index;comment:归属用户ID" json:"user_id"`
	ContactName  string                  `gorm:"column:contact_name;size:50;not null;comment:联系人姓名" json:"contact_name"`
	ContactPhone string                  `gorm:"column:contact_phone;size:20;not null;comment:联系电话" json:"contact_phone"`
	Building     string                  `gorm:"column:building;size:50;not null;comment:宿舍楼/教学楼" json:"building"`
	Room         string                  `gorm:"column:room;size:20;not null;comment:具体房间号" json:"room"`
	IsDefault    enum.AddressDefaultType `gorm:"column:is_default;not null;default:0;comment:是否默认地址 0-否 1-是" json:"is_default"`

	// ✅ 保持 string 类型！只加这两个标签
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime;comment:更新时间" json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间" json:"-"`
}

// TableName 强制绑定数据表名
func (Address) TableName() string {
	return "address"
}
