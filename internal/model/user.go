package model

import (
	"CampusTake/common/enum"
	"time"

	"gorm.io/gorm"
)

// User 用户基础信息表
type User struct {
	ID        int64           `gorm:"primaryKey;comment:用户ID"`
	Phone     string          `gorm:"size:20;uniqueIndex;not null;comment:手机号"`
	Password  string          `gorm:"size:100;not null;comment:密码"`
	Nickname  string          `gorm:"size:50;default:'';comment:昵称"`
	Avatar    string          `gorm:"size:255;default:'';comment:头像"`
	Gender    int8            `gorm:"not null;default:0;comment:性别 0未知 1男 2女"`
	Role      enum.RoleType   `gorm:"not null;default:1;comment:角色"`
	Status    enum.UserStatus `gorm:"not null;default:1;comment:状态"`
	CreatedAt time.Time       `gorm:"comment:创建时间"`
	UpdatedAt time.Time       `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt  `gorm:"index;comment:软删除时间"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
