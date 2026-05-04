// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"CampusTake/common/response"
	"CampusTake/internal/logic/user/profile"
	"net/http"

	"CampusTake/internal/svc"
)

func GetProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := profile.NewGetProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetProfile()
		response.Response(r, w, resp, err)
	}
}
