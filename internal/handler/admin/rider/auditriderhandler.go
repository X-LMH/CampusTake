// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/common/httpxext"
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/admin/rider"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func AuditRiderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuditRiderRequest
		if err := httpxext.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := rider.NewAuditRiderLogic(r.Context(), svcCtx)
		err := l.AuditRider(&req)
		response.Response(r, w, nil, err)
	}
}
