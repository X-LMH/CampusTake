package query

import "CampusTake/internal/enums"

type AppealQuery struct {
	Status      *enums.AppealStatus
	ApplicantID *int64
	OrderID     *int64
	Type        *enums.AppealType

	Page int
	Size int
}
