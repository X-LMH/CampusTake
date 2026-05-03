// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/common/ctxx"
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
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
	result := enum.AdminAuditResult(req.Result)
	var targetStatus enum.RiderAuditStatus
	if result.Agree() {
		targetStatus = enum.RiderStatusApproved
	} else {
		targetStatus = enum.RiderStatusRejected
	}

	// 3. 开启事务：先查询校验，后执行更新
	return l.svcCtx.Repo.WithTx(l.ctx, func(txRepo *repo.Repo) error {

		// --- 步骤 1: 查询当前申请记录并加锁 (使用 txRepo) ---
		// 建议使用事务内的查询，确保数据一致性
		profile, err := txRepo.Rider().GetProfileByUserID(l.ctx, req.RiderID)
		if err != nil {
			return err
		}

		// --- 步骤 2: 状态机校验 ---
		// 如果已经是目标状态（例如重复点击“通过”），则视为幂等，直接返回成功
		if profile.AuditStatus == targetStatus {
			l.Infof("申请 ID: %d 已经是目标状态 %v，无需重复操作", req.RiderID, targetStatus)
			return nil
		}

		// 如果当前不是“待审核”状态（例如已经是“撤回”或“拒绝”），则禁止修改
		if profile.AuditStatus != enum.RiderStatusPending {
			return errx.NewParamError("只有待审核的申请才能进行此操作")
		}

		// --- 步骤 3: 更新主表状态和备注 ---
		err = txRepo.Rider().UpdateStatusAndRemarkByUserID(l.ctx, req.RiderID, targetStatus, req.Remark)
		if err != nil {
			l.Errorf("更新骑手状态失败: %v", err)
			return err
		}

		// --- 步骤 4: 计入审核日志 ---
		auditLog := &model.RiderAuditLog{
			RiderID:   req.RiderID,
			AuditorID: adminID,
			Result:    enum.AdminAuditResult(int8(result)),
			Remark:    req.Remark,
		}
		err = txRepo.Rider().CreateLog(l.ctx, auditLog)
		if err != nil {
			l.Errorf("记录审核日志失败: %v", err)
			return err
		}

		return nil
	})
}
