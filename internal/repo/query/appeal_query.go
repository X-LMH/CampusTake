package query

import (
	"CampusTake/internal/enums"
	errs "CampusTake/pkg/errors"
	"time"
)

type AppealListQuery struct {
	Status        *enums.AppealStatus
	ApplicantID   *int64
	OrderID       *int64
	ApplicantRole *enums.AppealApplicantRole

	Page int
	Size int
}

type ReapplyAppealParams struct {
	AppealType   int8
	Content      string
	EvidenceUrls []byte
}
type HandleAppealParams struct {
	HandleRemark   string
	RefundAmount   float64
	PunishRider    int8
	TerminateOrder int8
	HandledBy      int64
	HandledAt      time.Time
}

func (p *HandleAppealParams) Validate(
	result int8,
) error {

	// =========================
	// 驳回
	// =========================

	if result == int8(enums.AppealHandleResultFailed) {

		if p.RefundAmount > 0 {
			return errs.ErrRejectedAppealCannotRefund
		}

		if p.TerminateOrder == 1 {
			return errs.ErrRejectedAppealCannotTerminateOrder
		}

		if p.PunishRider == 1 {
			return errs.ErrRejectedAppealCannotPunishRider
		}
	}

	return nil
}
