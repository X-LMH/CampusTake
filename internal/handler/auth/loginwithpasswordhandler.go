// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/common/httpxext"
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/auth"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func LoginWithPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginWithPasswordRequest
		if err := httpxext.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := auth.NewLoginWithPasswordLogic(r.Context(), svcCtx)
		resp, err := l.LoginWithPassword(&req)
		response.Response(r, w, resp, err)
	}
}
