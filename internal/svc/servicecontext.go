package svc

import (
	"CampusTake/common/cache"
	"CampusTake/common/db"
	"CampusTake/common/jwtx"
	"CampusTake/internal/config"
	"CampusTake/internal/middleware"
	"CampusTake/internal/repo"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config            config.Config
	Repo              *repo.Repo
	JwtAuthMiddleware rest.Middleware
	JwtCfg            jwtx.JwtConfig
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQLConfig.User,
		c.MySQLConfig.Password,
		c.MySQLConfig.Host,
		c.MySQLConfig.Port,
		c.MySQLConfig.Database,
	)

	dbConn := db.NewMySQLConnection(
		dsn,
		c.MySQLConfig.MaxOpenConns,
		c.MySQLConfig.MaxIdleConns,
	)

	rdb := cache.NewRedisClient(
		c.RedisConfig.Addr,
		c.RedisConfig.Password,
		c.RedisConfig.DB,
	)

	jwtCfg := jwtx.JwtConfig{
		SecretKey: c.JwtAuth.SecretKey,
		Issuer:    c.JwtAuth.Issuer,
		Expire:    time.Duration(c.JwtAuth.Expire) * time.Second,
	}

	return &ServiceContext{
		Config:            c,
		JwtAuthMiddleware: middleware.NewJwtAuthMiddleware(c.JwtAuth.SecretKey).Handle,
		JwtCfg:            jwtCfg,
		Repo:              repo.NewRepo(dbConn, rdb),
	}
}
