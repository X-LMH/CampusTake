// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/address"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func DeleteAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAddressRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := address.NewDeleteAddressLogic(r.Context(), svcCtx)
		err := l.DeleteAddress(&req)
		response.Response(r, w, nil, err)
	}
}
