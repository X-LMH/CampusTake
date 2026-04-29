package enum

// -------------------------- 用户角色枚举 --------------------------
type RoleType int8

const (
	RoleUser     RoleType = 1 // 普通用户
	RoleReceiver RoleType = 2 // 接单用户
	RoleAdmin    RoleType = 3 // 管理员
)

func (r RoleType) String() string {
	switch r {
	case RoleUser:
		return "普通用户"
	case RoleReceiver:
		return "接单用户"
	case RoleAdmin:
		return "管理员"
	default:
		return "未知"
	}
}
