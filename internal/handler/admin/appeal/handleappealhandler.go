// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/admin/appeal"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func HandleAppealHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HandleAppealRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := appeal.NewHandleAppealLogic(r.Context(), svcCtx)
		err := l.HandleAppeal(&req)
		response.Response(r, w, nil, err)
	}
}
