package repo

import (
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

	User       UserRepo
	Address    AddressRepo
	Rider      RiderRepo
	VerifyCode VerifyCodeRepo
	Token      TokenRepo
	Order      OrderRepo
	Payment    PaymentRepo
	Review     ReviewRepo
}

func NewRepo(db *gorm.DB, rdb *redis.Client) *Repo {
	return &Repo{
		db:         db,
		rdb:        rdb,
		User:       NewUserRepo(db),
		Address:    NewAddressRepo(db),
		Rider:      NewRiderRepo(db),
		VerifyCode: NewVerifyCodeRepo(rdb),
		Token:      NewTokenRepo(rdb),
		Order:      NewOrderRepo(db),
		Payment:    NewPaymentRepo(db),
		Review:     NewReviewRepo(db),
	}
}

type RepoTx struct {
	User    UserRepo
	Address AddressRepo
	Rider   RiderRepo
	Order   OrderRepo
	Payment PaymentRepo
	Review  ReviewRepo
}

func (r *Repo) WithTx(ctx context.Context, fn func(txRepo *RepoTx) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &RepoTx{
			User:    NewUserRepo(tx),
			Address: NewAddressRepo(tx),
			Rider:   NewRiderRepo(tx),
			Order:   NewOrderRepo(tx),
			Payment: NewPaymentRepo(tx),
			Review:  NewReviewRepo(tx),
		}
		return fn(txRepo)
	})
}
