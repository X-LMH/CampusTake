// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/admin/order"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func GetAppealListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAppealListRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := order.NewGetAppealListLogic(r.Context(), svcCtx)
		resp, err := l.GetAppealList(&req)
		response.Response(r, w, resp, err)
	}
}
