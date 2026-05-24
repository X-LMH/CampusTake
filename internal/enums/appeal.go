package enums

type AppealApplicantRole int8

const (
	AppealApplicantRoleUser  AppealApplicantRole = 1 // 用户申诉
	AppealApplicantRoleRider AppealApplicantRole = 2 // 骑手申诉
)

func (a AppealApplicantRole) String() string {
	switch a {
	case AppealApplicantRoleUser:
		return "用户"
	case AppealApplicantRoleRider:
		return "骑手"
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

type AppealHandleResult int8

const (
	AppealHandleResultDefault AppealHandleResult = 0 // 默认
	AppealHandleResultSuccess AppealHandleResult = 1 // 处理成功
	AppealHandleResultFailed  AppealHandleResult = 2 // 处理失败
)

func (a AppealHandleResult) String() string {
	switch a {
	case AppealHandleResultDefault:
		return "默认"
	case AppealHandleResultSuccess:
		return "申诉通过"
	case AppealHandleResultFailed:
		return "申诉驳回"
	default:
		return "未知处理结果"
	}
}
