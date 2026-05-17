package enums

type AppealType int8

const (
	AppealTypeUser  AppealType = 1 // 用户申诉
	AppealTypeRider AppealType = 2 // 骑手申诉
)

// String 返回 AppealType 对应的字符串
func (a AppealType) String() string {
	switch a {
	case AppealTypeUser:
		return "用户申诉"
	case AppealTypeRider:
		return "骑手申诉"
	default:
		return "未知申诉类型"
	}
}

type AppealStatus int8

const (
	AppealStatusDefault   AppealStatus = 0 // 默认
	AppealStatusPending   AppealStatus = 1 // 待处理
	AppealStatusApproved  AppealStatus = 2 // 已通过
	AppealStatusRejected  AppealStatus = 3 // 已驳回
	AppealStatusCancelled AppealStatus = 4 // 已撤销
)

// String 返回 AppealStatus 对应的字符串
func (a AppealStatus) String() string {
	switch a {
	case AppealStatusDefault:
		return "默认"
	case AppealStatusPending:
		return "待处理"
	case AppealStatusApproved:
		return "已通过"
	case AppealStatusRejected:
		return "已驳回"
	case AppealStatusCancelled:
		return "已撤销"
	default:
		return "未知申诉状态"
	}
}

// CanTransferTo 判断状态是否允许流转
func (a AppealStatus) CanTransferTo(to AppealStatus) bool {

	switch a {

	// 默认 -> 待处理
	case AppealStatusDefault:
		return to == AppealStatusPending

	// 待处理 -> 已通过/已驳回/已撤销
	case AppealStatusPending:
		return to == AppealStatusApproved ||
			to == AppealStatusRejected ||
			to == AppealStatusCancelled

	// 已结束状态禁止再次流转
	case AppealStatusApproved,
		AppealStatusRejected,
		AppealStatusCancelled:
		return false

	default:
		return false
	}
}
