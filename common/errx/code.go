package errx

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
	ErrUserNotFound          = NewCodeError(1000, "用户不存在")
	ErrUserExist             = NewCodeError(1001, "用户已存在")
	ErrPasswordWrong         = NewCodeError(1002, "用户名或密码错误")
	ErrUserForbidden         = NewCodeError(1003, "用户被禁用")
	ErrAddressNotFound       = NewCodeError(1100, "地址不存在")
	ErrVerifyCodeTooFrequent = NewCodeError(1200, "验证码请求过于频繁")
	ErrVerifyCodeWrong       = NewCodeError(1201, "验证码错误")
	ErrPasswordNoChange      = NewCodeError(1300, "新密码不能与旧密码相同")
	ErrVerifyCodeNotFound    = NewCodeError(1202, "验证码不存在或已过期")
)
