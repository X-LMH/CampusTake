// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/rider"
	"CampusTake/internal/svc"
)

func CancelApplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := rider.NewCancelApplyLogic(r.Context(), svcCtx)
		err := l.CancelApply()
		response.Response(r, w, nil, err)
	}
}
