package query

type OrderStatusUpdateQuery struct {
	OrderID int64

	UserID  *int64
	RiderID *int64
}
