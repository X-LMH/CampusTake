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
	rdb *redis.Client

	User       impl.UserRepo
	Address    impl.AddressRepo
	Rider      impl.RiderRepo
	VerifyCode impl.VerifyCodeRepo
	Token      impl.TokenRepo
	Order      impl.OrderRepo
	Payment    impl.PaymentRepo
	Review     impl.ReviewRepo
	Appeal     impl.AppealRepo
}

func NewRepo(db *gorm.DB, rdb *redis.Client) *Repo {
	return &Repo{
		db:         db,
		rdb:        rdb,
		User:       impl.NewUserRepo(db),
		Address:    impl.NewAddressRepo(db),
		Rider:      impl.NewRiderRepo(db),
		VerifyCode: impl.NewVerifyCodeRepo(rdb),
		Token:      impl.NewTokenRepo(rdb),
		Order:      impl.NewOrderRepo(db),
		Payment:    impl.NewPaymentRepo(db),
		Review:     impl.NewReviewRepo(db),
		Appeal:     impl.NewAppealRepo(db),
	}
}

type RepoTx struct {
	User    impl.UserRepo
	Address impl.AddressRepo
	Rider   impl.RiderRepo
	Order   impl.OrderRepo
	Payment impl.PaymentRepo
	Review  impl.ReviewRepo
	Appeal  impl.AppealRepo
}

func (r *Repo) WithTx(ctx context.Context, fn func(txRepo *RepoTx) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &RepoTx{
			User:    impl.NewUserRepo(tx),
			Address: impl.NewAddressRepo(tx),
			Rider:   impl.NewRiderRepo(tx),
			Order:   impl.NewOrderRepo(tx),
			Payment: impl.NewPaymentRepo(tx),
			Review:  impl.NewReviewRepo(tx),
			Appeal:  impl.NewAppealRepo(tx),
		}
		return fn(txRepo)
	})
}
