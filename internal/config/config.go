// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"CampusTake/pkg/mq"

	"github.com/zeromicro/go-zero/rest"
)

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
		Addr         string
		Password     string
		DB           int
		PoolSize     int
		MinIdleConns int
	}
	Performance struct {
		HTTPLog bool
	}
	JwtAuth struct {
		SecretKey string
		Issuer    string
		Expire    int64 // 秒
	}
	UploadConfig   UploadConfig
	RabbitMQConfig struct {
		User         string
		Password     string
		Host         string
		Port         string
		VirtualHost  string
		OrderCancel  mq.DelayQueueConfig
		OrderConfirm mq.DelayQueueConfig
		RiderCheck   mq.DelayQueueConfig
	}
}

type UploadConfig struct {
	UrlPrefix      string
	MaxImageSizeMB int64

	Avatar     UploadPathConfig
	Appeal     UploadPathConfig
	CampusCard UploadPathConfig
}
type UploadPathConfig struct {
	Path       string
	PathPrefix string
}
