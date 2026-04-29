package response

import (
	"errors"
	"net/http"

	"CampusTake/common/errx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 统一响应结构
type Body struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Response 统一出口
func Response(r *http.Request, w http.ResponseWriter, res interface{}, err error) {
	if err != nil {
		handleError(r, w, err)
		return
	}

	httpx.WriteJson(w, http.StatusOK, Body{
		Code: int(errx.Success),
		Data: res,
		Msg:  "success",
	})
}

// handleError 内部错误处理
func handleError(r *http.Request, w http.ResponseWriter, err error) {
	code := int(errx.ServerCommonError)
	msg := "服务器开小差了，请稍后再试"

	var e *errx.CodeError
	if errors.As(err, &e) {
		code = int(e.Code)
		if e.Code != errx.ServerCommonError {
			msg = e.Msg
		}
	}

	httpx.WriteJson(w, http.StatusOK, Body{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
