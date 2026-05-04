// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"CampusTake/common/httpxext"
	"CampusTake/common/response"
	"CampusTake/internal/logic/user/profile"
	"net/http"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func GetAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAvatarRequest
		if err := httpxext.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err) // 参数错误继续走 JSON 错误通道
			return
		}

		l := profile.NewGetAvatarLogic(r.Context(), svcCtx)
		avatarPath, err := l.GetAvatar(&req)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		http.ServeFile(w, r, avatarPath)
	}
}
