// common/httpx/bind.go
package httpx

import (
	"CampusTake/pkg/errors"
	"CampusTake/pkg/validator"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BindAndValidate(r *http.Request, v interface{}) error {
	if err := httpx.Parse(r, v); err != nil {
		logx.Errorf("parse params error: %v | parsed=%+v", err, v)
		return errors.NewParamError("请求参数格式错误")
	}

	if err := validator.ValidateStruct(v); err != nil {
		msg := validator.TranslateValidationError(err)
		logx.Errorf("validate params error: %v | parsed=%+v | msg=%s", err, v, msg)
		return errors.NewParamError(msg)
	}

	return nil
}
