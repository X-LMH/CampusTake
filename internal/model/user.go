package model

import (
	"CampusTake/common/enum"
	"time"
)

// User 对应数据库 user 表
type User struct {
	ID        uint64        `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID"`
	Phone     string        `gorm:"column:phone;type:varchar(20);not null;uniqueIndex:uk_phone;comment:手机号"`
	Password  string        `gorm:"column:password;type:varchar(255);not null;comment:密码"`
	Nickname  string        `gorm:"column:nickname;type:varchar(50);not null;default:'默认昵称';comment:用户昵称"`
	Avatar    string        `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像URL"`
	Status    int8          `gorm:"column:status;type:tinyint;not null;default:1;comment:状态: 1-正常, 2-封禁"`
	Role      enum.RoleType `gorm:"column:role;type:tinyint;not null;default:1;comment:角色: 1-普通用户, 2-接单用户, 3-管理员"`
	CreatedAt time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt *time.Time    `gorm:"column:deleted_at;type:timestamp;default:null;index:idx_deleted_at;comment:删除时间"`
}

func (User) TableName() string {
	return "user"
}
