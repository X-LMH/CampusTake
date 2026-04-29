// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySQLConfig struct {
		User         string
		Password     string
		Host         string
		Port         string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
	}
	RedisConfig struct {
		Addr     string
		Password string
		DB       int
	}
	JwtAuth struct {
		SecretKey string
		Issuer    string
		Expire    int64 // 秒
	}
}
