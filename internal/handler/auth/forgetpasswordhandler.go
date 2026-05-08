// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/auth"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func ForgetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ForgetPasswordRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := auth.NewForgetPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ForgetPassword(&req)
		response.Response(r, w, resp, err)
	}
}
