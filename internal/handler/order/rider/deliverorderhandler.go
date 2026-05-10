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

func DeliverOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeliverOrderRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := rider.NewDeliverOrderLogic(r.Context(), svcCtx)
		err := l.DeliverOrder(&req)
		response.Response(r, w, nil, err)
	}
}
