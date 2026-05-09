// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"CampusTake/pkg/utils"
	"context"
	"mime/multipart"
	"path"
	"path/filepath"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyRiderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyRiderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyRiderLogic {
	return &ApplyRiderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyRiderLogic) ApplyRider(req *types.ApplyRiderRequest,
	frontFile multipart.File, frontHeader *multipart.FileHeader,
	backFile multipart.File, backHeader *multipart.FileHeader) error {

	userID := ctxx.MustUserID(l.ctx)

	// 1. 状态校验：查询现有记录
	// 注意：此处忽略 RecordNotFound 错误，因为新用户申请查不到是正常的
	oldProfile, err := l.svcCtx.Repo.Rider.GetProfileByUserID(l.ctx, userID)
	if err == nil && oldProfile != nil {
		// 【修改点】仅拦截“待审核”和“已通过”，允许“已拒绝(2)”和“已撤回(3)”状态继续向下走
		if oldProfile.AuditStatus == enums.RiderStatusPending {
			return errors.ErrApplyRiderDuplicate
		}
		if oldProfile.AuditStatus == enums.RiderStatusApproved {
			return errors.ErrApplyRiderAlready
		}
	}

	// 2. 存储图片
	saveDir := filepath.Join(l.svcCtx.Config.Upload.CampusCardPath, strconv.FormatInt(userID, 10))
	frontFileName, err := utils.SaveFileToLocal(frontFile, frontHeader, saveDir, "front")
	if err != nil {
		return err
	}
	backFileName, err := utils.SaveFileToLocal(backFile, backHeader, saveDir, "back")
	if err != nil {
		return err
	}

	// 3. 构建模型
	// 【修改点】使用 path.Join (跨平台 URL 兼容) 而非 filepath.Join
	cardFrontURL := path.Join(l.svcCtx.Config.Upload.CampusCardPathPrefix, strconv.FormatInt(userID, 10), frontFileName)
	cardBackURL := path.Join(l.svcCtx.Config.Upload.CampusCardPathPrefix, strconv.FormatInt(userID, 10), backFileName)

	profile := &model.RiderProfile{
		UserID:            userID,
		RealName:          req.RealName,
		StudentNo:         req.StudentNo,
		IDCardNo:          req.IdCardNo,
		DormitoryBuilding: req.DormitoryBuilding,
		DormitoryRoom:     req.DormitoryRoom,
		CampusCardFront:   cardFrontURL,
		CampusCardBack:    cardBackURL,
		AuditStatus:       enums.RiderStatusPending,
	}

	// 4. 【修改点】调用 UpsertProfile 自动识别插入或更新
	err = l.svcCtx.Repo.Rider.UpsertProfile(l.ctx, profile)
	if err != nil {
		return err
	}

	return nil
}
