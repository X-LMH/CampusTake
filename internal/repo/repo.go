package repo

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repo struct {
	db  *gorm.DB
	rdb *redis.Client

	user       UserRepo
	address    AddressRepo
	verifyCode VerifyCodeRepo
	token      TokenRepo
}

func NewRepo(db *gorm.DB, rdb *redis.Client) *Repo {
	return &Repo{
		db:  db,
		rdb: rdb,
	}
}

func (r *Repo) User() UserRepo {
	return &userRepo{db: r.db}
}

func (r *Repo) Address() AddressRepo {
	return &addressRepo{db: r.db}
}

func (r *Repo) VerifyCode() VerifyCodeRepo {
	return &verifyCodeRepo{rdb: r.rdb}
}

func (r *Repo) Token() TokenRepo {
	return &tokenRepo{rdb: r.rdb}
}

func (r *Repo) WithTx(ctx context.Context, fn func(r *Repo) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepo(tx, r.rdb)
		return fn(txRepo)
	})
}
