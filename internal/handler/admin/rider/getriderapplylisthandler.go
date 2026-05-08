// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/admin/rider"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func GetRiderApplyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRiderApplyListRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := rider.NewGetRiderApplyListLogic(r.Context(), svcCtx)
		resp, err := l.GetRiderApplyList(&req)
		response.Response(r, w, resp, err)
	}
}
