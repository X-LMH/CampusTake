package enum

// -------------------------- 用户角色枚举 --------------------------
type RoleType int8

const (
	RoleUser  RoleType = 1 // 普通用户
	RoleRider RoleType = 2 // 代取员 / 接单员
	RoleAdmin RoleType = 3 // 管理员
)

func (r RoleType) String() string {
	switch r {
	case RoleUser:
		return "普通用户"
	case RoleRider:
		return "代取员"
	case RoleAdmin:
		return "管理员"
	default:
		return "未知"
	}
}

func (r RoleType) IsAdmin() bool {
	return r == RoleAdmin
}
func (r RoleType) IsRider() bool {
	return r == RoleRider
}

// -------------------------- 用户状态枚举 --------------------------
type UserStatus int8

const (
	UserStatusNormal   UserStatus = 1 // 正常
	UserStatusDisabled UserStatus = 2 // 禁用
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusNormal:
		return "正常"
	case UserStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

type RiderAuditStatus int8

const (
	RiderStatusPending  RiderAuditStatus = 1 // 待审核
	RiderStatusApproved RiderAuditStatus = 2 // 审核通过
	RiderStatusRejected RiderAuditStatus = 3 // 审核拒绝
	RiderStatusCancel   RiderAuditStatus = 4 // 撤销申请
)
