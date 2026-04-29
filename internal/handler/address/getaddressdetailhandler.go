// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/httpxext"
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/address"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func GetAddressDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAddressDetailRequest
		if err := httpxext.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := address.NewGetAddressDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetAddressDetail(&req)
		response.Response(r, w, resp, err)
	}
}
