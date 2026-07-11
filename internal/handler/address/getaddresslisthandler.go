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

func GetAddressListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAddressListRequest
		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := address.NewGetAddressListLogic(r.Context(), svcCtx)
		resp, err := l.GetAddressList(&req)
		response.Response(r, w, resp, err)
	}
}
