package errx

import "fmt"

// CodeError 自定义错误结构
type CodeError struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
}

// 实现 error 接口
func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.Code, e.Msg)
}

// 构造函数
func NewCodeError(code Code, msg string) *CodeError {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}

// 新建参数错误
func NewParamError(msg string) *CodeError {
	return NewCodeError(RequestParamError, msg)
}

// 默认通用错误
func NewDefaultError(msg string) *CodeError {
	return NewCodeError(ServerCommonError, msg)
}
