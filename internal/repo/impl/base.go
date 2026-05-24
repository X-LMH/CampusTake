package impl

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RepoBase struct {
	db  *gorm.DB
	rdb redis.Cmdable
}
