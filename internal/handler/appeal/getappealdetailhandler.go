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

func GetAppealDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAppealDetailRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := appeal.NewGetAppealDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetAppealDetail(&req)
		response.Response(r, w, resp, err)
	}
}
