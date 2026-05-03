// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/common/httpxext"
	"CampusTake/common/response"
	"net/http"

	"CampusTake/internal/logic/rider"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
)

func ApplyRiderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyRiderRequest
		if err := httpxext.BindAndValidate(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 2. 获取正面照片
		frontFile, frontHeader, err := r.FormFile("card_front")
		if err != nil {
			response.Response(r, w, nil, err) // 这里可以封装一个“请上传正面照”的错误
			return
		}
		defer frontFile.Close()

		// 3. 获取反面照片
		backFile, backHeader, err := r.FormFile("card_back")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		defer backFile.Close()

		l := rider.NewApplyRiderLogic(r.Context(), svcCtx)
		err = l.ApplyRider(&req, frontFile, frontHeader, backFile, backHeader)
		response.Response(r, w, nil, err)
	}
}
