// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package upload

import (
	"CampusTake/pkg/httpx"
	"CampusTake/pkg/response"
	"net/http"

	"CampusTake/internal/logic/upload"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

// 上传图片
func UploadImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadImageRequest

		if err := httpx.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := upload.NewUploadImageLogic(r.Context(), svcCtx)

		resp, err := l.UploadImage(r, &req)
		response.Response(r, w, resp, err)
	}
}
