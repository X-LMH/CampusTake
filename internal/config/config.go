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
	Upload struct {
		UrlPrefix            string
		AvatarPath           string
		AvatarPathPrefix     string
		CampusCardPath       string
		CampusCardPathPrefix string
	}
	RabbitMQConfig struct {
		User                  string
		Password              string
		Host                  string
		Port                  string
		VirtualHost           string
		OrderDelayExchange    string
		OrderDelayQueue       string
		OrderDelayRoutingKey  string
		OrderCancelExchange   string
		OrderCancelQueue      string
		OrderCancelRoutingKey string
		TTL                   int32 // 毫秒数，建议用 int32 方便后续传参
	}
}
