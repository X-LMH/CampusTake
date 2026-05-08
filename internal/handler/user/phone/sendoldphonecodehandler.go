// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/user/phone"
	"CampusTake/internal/svc"
)

func SendOldPhoneCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := phone.NewSendOldPhoneCodeLogic(r.Context(), svcCtx)
		err := l.SendOldPhoneCode()
		response.Response(r, w, nil, err)
	}
}
