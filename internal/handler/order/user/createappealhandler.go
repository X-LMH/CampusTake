// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/order/user"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func CreateAppealHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateAppealRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := user.NewCreateAppealLogic(r.Context(), svcCtx)
		err := l.CreateAppeal(&req)
		response.Response(r, w, nil, err)
	}
}
