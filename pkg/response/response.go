package response

import (
	errx2 "CampusTake/pkg/errors"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
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
		Code: int(errx2.Success),
		Data: res,
		Msg:  "success",
	})
}

// handleError 内部错误处理
func handleError(_ *http.Request, w http.ResponseWriter, err error) {
	code := int(errx2.ServerCommonError)
	msg := "服务器开小差了，请稍后再试"

	var e *errx2.CodeError
	if errors.As(err, &e) {
		code = int(e.Code)
		if e.Code != errx2.ServerCommonError {
			msg = e.Msg
		}
	}

	if e != nil {
		if e.Code == errx2.ServerCommonError {
			// 若是通用服务错误，记录 CodeError（含 code/msg）
			logx.Error(e)
		}
	} else {
		// 非 CodeError，视为内部错误，记录原始错误以便排查
		logx.Error(err)
	}

	httpx.WriteJson(w, http.StatusOK, Body{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
