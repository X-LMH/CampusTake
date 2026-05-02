// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/user"
	"CampusTake/internal/svc"
)

func GetProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetProfile()
		response.Response(r, w, resp, err)
	}
}
