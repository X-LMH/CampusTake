// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"CampusTake/common/response"
	"CampusTake/internal/logic/user/profile"
	"net/http"

	"CampusTake/internal/svc"
)

func UpdateAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1️⃣ 取文件
		file, header, err := r.FormFile("avatar")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		defer file.Close()

		l := profile.NewUpdateAvatarLogic(r.Context(), svcCtx)

		resp, err := l.UpdateAvatar(file, header)
		response.Response(r, w, resp, err)
	}
}
