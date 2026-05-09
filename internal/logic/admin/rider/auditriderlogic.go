// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditRiderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditRiderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditRiderLogic {
	return &AuditRiderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *AuditRiderLogic) AuditRider(req *types.AuditRiderRequest) error {
	// 1. 获取管理员ID（用于记录日志）
	adminID := ctxx.MustUserID(l.ctx)

	// 2. 转换审核结果
	// 1通过，2拒绝
	result := enums.AdminAuditResult(req.Result)
	var targetStatus enums.RiderAuditStatus
	if result.Agree() {
		targetStatus = enums.RiderStatusApproved
	} else {
		targetStatus = enums.RiderStatusRejected
		if req.Remark == "" {
			return errors.NewParamError("拒绝时备注不能为空")
		}
	}

	// 3. 开启事务：先查询校验，后执行更新
	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		profile, err := tx.Rider.GetProfileByUserID(l.ctx, req.UserID)
		if err != nil {
			return err
		}

		// --- 步骤 2: 状态机校验 ---
		// 如果已经是目标状态（例如重复点击“通过”），则视为幂等，直接返回成功
		if profile.AuditStatus == targetStatus {
			l.Infof("申请 ID: %d 已经是目标状态 %v，无需重复操作", req.UserID, targetStatus)
			return nil
		}

		// 如果当前不是“待审核”状态（例如已经是“撤回”或“拒绝”），则禁止修改
		if profile.AuditStatus != enums.RiderStatusPending {
			return errors.NewParamError("只有待审核的申请才能进行此操作")
		}

		// --- 步骤 3: 更新主表状态和备注 ---
		err = tx.Rider.UpdateStatusAndRemarkByUserID(l.ctx, req.UserID, targetStatus, req.Remark)
		if err != nil {
			l.Errorf("更新骑手状态失败: %v", err)
			return err
		}

		// 更新一下用户表中的骑手状态，保持数据一致
		err = tx.User.UpdateRoleByID(l.ctx, req.UserID, enums.RoleRider)
		if err != nil {
			l.Errorf("更新用户角色失败: %v", err)
			return err
		}

		// --- 步骤 4: 计入审核日志 ---
		auditLog := &model.RiderAuditLog{
			RiderID:   profile.ID,
			AuditorID: adminID,
			Result:    result,
			Remark:    req.Remark,
		}
		err = tx.Rider.CreateLog(l.ctx, auditLog)
		if err != nil {
			l.Errorf("记录审核日志失败: %v", err)
			return err
		}

		return nil
	})
}
