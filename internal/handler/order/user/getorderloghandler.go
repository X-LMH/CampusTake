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

func GetOrderLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrderLogRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := user.NewGetOrderLogLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderLog(&req)
		response.Response(r, w, resp, err)
	}
}
