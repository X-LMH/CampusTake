// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/common/httpxext"
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/user/phone"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func VerifyNewPhoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyNewPhoneRequest
		if err := httpxext.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := phone.NewVerifyNewPhoneLogic(r.Context(), svcCtx)
		resp, err := l.VerifyNewPhone(&req)
		response.Response(r, w, resp, err)
	}
}
