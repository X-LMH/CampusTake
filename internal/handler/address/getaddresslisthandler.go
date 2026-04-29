// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/address"
	"CampusTake/internal/svc"
)

func GetAddressListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewGetAddressListLogic(r.Context(), svcCtx)
		resp, err := l.GetAddressList()
		response.Response(r, w, resp, err)
	}
}
