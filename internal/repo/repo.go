package repo

import (
	"CampusTake/internal/repo/impl"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ---------------------
// 主 Repo
// ---------------------

type Repo struct {
	db  *gorm.DB
	rdb redis.Cmdable

	User         impl.UserRepo
	Address      impl.AddressRepo
	Rider        impl.RiderRepo
	VerifyCode   impl.VerifyCodeRepo
	Token        impl.TokenRepo
	Order        impl.OrderRepo
	Payment      impl.PaymentRepo
	Review       impl.ReviewRepo
	Appeal       impl.AppealRepo
	AppealHandle impl.AppealHandleRepo
	Refund       impl.RefundRepo
}

func NewRepo(db *gorm.DB, rdb redis.Cmdable) *Repo {
	return &Repo{
		db:  db,
		rdb: rdb,

		User:         impl.NewUserRepo(db, rdb),
		Address:      impl.NewAddressRepo(db, rdb),
		Rider:        impl.NewRiderRepo(db, rdb),
		VerifyCode:   impl.NewVerifyCodeRepo(db, rdb),
		Token:        impl.NewTokenRepo(db, rdb),
		Order:        impl.NewOrderRepo(db, rdb),
		Payment:      impl.NewPaymentRepo(db, rdb),
		Review:       impl.NewReviewRepo(db, rdb),
		Appeal:       impl.NewAppealRepo(db, rdb),
		AppealHandle: impl.NewAppealHandleRepo(db, rdb),
		Refund:       impl.NewRefundRepo(db, rdb),
	}
}

// ---------------------
// 事务 Repo
// ---------------------

type RepoTx struct {
	User         impl.UserRepo
	Address      impl.AddressRepo
	Rider        impl.RiderRepo
	Order        impl.OrderRepo
	Payment      impl.PaymentRepo
	Review       impl.ReviewRepo
	Appeal       impl.AppealRepo
	AppealHandle impl.AppealHandleRepo
	Refund       impl.RefundRepo
}

// WithTx 开启事务
func (r *Repo) WithTx(
	ctx context.Context,
	fn func(txRepo *RepoTx) error,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txRepo := &RepoTx{
			User:         impl.NewUserRepo(tx, r.rdb),
			Address:      impl.NewAddressRepo(tx, r.rdb),
			Rider:        impl.NewRiderRepo(tx, r.rdb),
			Order:        impl.NewOrderRepo(tx, r.rdb),
			Payment:      impl.NewPaymentRepo(tx, r.rdb),
			Review:       impl.NewReviewRepo(tx, r.rdb),
			Appeal:       impl.NewAppealRepo(tx, r.rdb),
			AppealHandle: impl.NewAppealHandleRepo(tx, r.rdb),
			Refund:       impl.NewRefundRepo(tx, r.rdb),
		}

		return fn(txRepo)
	})
}
