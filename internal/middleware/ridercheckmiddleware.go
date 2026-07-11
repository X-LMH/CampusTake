// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"CampusTake/pkg/response"
	"net/http"
)

type RiderCheckMiddleware struct {
}

func NewRiderCheckMiddleware() *RiderCheckMiddleware {
	return &RiderCheckMiddleware{}
}

func (m *RiderCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 直接通过 context 获取解析好的 role (基于你 ctxx 中封装的方法)
		role := ctxx.GetRole(r.Context())
		// 2. 校验是否为管理员
		if !role.IsRider() {
			response.Response(r, w, nil, errors.ErrUserPermissionDenied)
			return
		}
		next(w, r)
	}
}
