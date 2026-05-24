// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/appeal"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func GetMyAppealListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMyAppealListRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := appeal.NewGetMyAppealListLogic(r.Context(), svcCtx)
		resp, err := l.GetMyAppealList(&req)
		response.Response(r, w, resp, err)
	}
}
