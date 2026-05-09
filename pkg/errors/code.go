package errors

// Code 类型
type Code int

// 全局基础错误码
const (
	Success           Code = 0
	ServerCommonError Code = 10001
	RequestParamError Code = 10002
	TokenExpireError  Code = 10003
	TokenInvalidError Code = 10004
	TokenMissingError Code = 10005
)

// 用户相关错误 1000+ 模块
var (
	ErrInvalidParam = NewCodeError(RequestParamError, "参数错误")

	ErrUserNotFound         = NewCodeError(1000, "用户不存在")
	ErrUserExist            = NewCodeError(1001, "用户已存在")
	ErrPasswordWrong        = NewCodeError(1002, "用户名或密码错误")
	ErrUserForbidden        = NewCodeError(1003, "用户被禁用")
	ErrUserPermissionDenied = NewCodeError(1004, "用户权限不足")
	ErrPhoneAlreadyBound    = NewCodeError(1005, "手机号已绑定")
	ErrPhoneSameWithOld     = NewCodeError(1006, "新手机号与旧手机号一致")

	ErrAddressNotFound = NewCodeError(1100, "地址不存在")

	ErrVerifyCodeTooFrequent = NewCodeError(1200, "验证码请求过于频繁")
	ErrVerifyCodeWrong       = NewCodeError(1201, "验证码错误")
	ErrVerifyCodeNotFound    = NewCodeError(1202, "验证码不存在或已过期")
	ErrVerifyTokenNotFound   = NewCodeError(1203, "验证令牌不存在或已过期")

	ErrPasswordNoChange = NewCodeError(1300, "新密码不能与旧密码相同")

	ErrApplyRiderDuplicate = NewCodeError(1400, "申请审核中，请勿重复提交")
	ErrApplyRiderAlready   = NewCodeError(1401, "您已是骑手，无需再次申请")
	ErrCannotCancelStatus  = NewCodeError(1402, "当前状态无法撤销申请")
)

// 文件相关错误 2000+ 模块
var (
	ErrFileNotFound    = NewCodeError(2000, "文件不存在")
	ErrFileFormatError = NewCodeError(2001, "文件格式错误")
	ErrFileTooLarge    = NewCodeError(2002, "文件过大")
)

var (
	ErrOrderNotFound      = NewCodeError(3000, "订单不存在")
	ErrOrderStatusInvalid = NewCodeError(3001, "订单状态不合法")
)

var (
	ErrPaymentNotFound      = NewCodeError(4000, "支付记录不存在")
	ErrPaymentStatusInvalid = NewCodeError(4001, "支付状态不合法")
)
