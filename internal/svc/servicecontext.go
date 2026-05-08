package svc

import (
	"CampusTake/internal/config"
	"CampusTake/internal/middleware"
	"CampusTake/internal/repo"
	"CampusTake/pkg/cache"
	"CampusTake/pkg/db"
	"CampusTake/pkg/jwt"
	"CampusTake/pkg/snowflake"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config            config.Config
	Repo              *repo.Repo
	JwtAuthMiddleware rest.Middleware
	JwtCfg            jwt.JwtConfig
	AdminCheck        rest.Middleware
	RiderCheck        rest.Middleware
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

	jwtCfg := jwt.JwtConfig{
		SecretKey: c.JwtAuth.SecretKey,
		Issuer:    c.JwtAuth.Issuer,
		Expire:    time.Duration(c.JwtAuth.Expire) * time.Second,
	}

	// 初始化雪花算法节点
	err := snowflake.InitSnowFlake(1)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:            c,
		JwtAuthMiddleware: middleware.NewJwtAuthMiddleware(c.JwtAuth.SecretKey).Handle,
		JwtCfg:            jwtCfg,
		Repo:              repo.NewRepo(dbConn, rdb),
		AdminCheck:        middleware.NewAdminCheckMiddleware().Handle,
		RiderCheck:        middleware.NewRiderCheckMiddleware().Handle,
	}
}
