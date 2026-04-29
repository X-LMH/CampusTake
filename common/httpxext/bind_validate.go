// common/httpxext/bind_validate.go
package httpxext

import (
	"CampusTake/common/validatorx"
	"net/http"

	"CampusTake/common/errx"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BindAndValidate(r *http.Request, v interface{}) error {
	if err := httpx.Parse(r, v); err != nil {
		logx.Errorf("parse params error: %v | parsed=%+v", err, v)
		return errx.NewParamError("请求参数格式错误")
	}

	if err := validatorx.ValidateStruct(v); err != nil {
		msg := validatorx.TranslateValidationError(err)
		logx.Errorf("validate params error: %v | parsed=%+v | msg=%s", err, v, msg)
		return errx.NewParamError(msg)
	}

	return nil
}
