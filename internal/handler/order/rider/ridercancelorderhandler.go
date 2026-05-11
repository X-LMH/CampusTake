// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/order/rider"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func RiderCancelOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RiderCancelOrderRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := rider.NewRiderCancelOrderLogic(r.Context(), svcCtx)
		err := l.RiderCancelOrder(&req)
		response.Response(r, w, nil, err)
	}
}
