package svc

import (
	"github.com/go-redis/redis/v8"
	"wszero/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Cache[0].Host,
		Password: c.Cache[0].Pass, // no password set
		DB:       0,               // use default DB
	})
	return &ServiceContext{
		Config:      c,
		RedisClient: rdb,
	}
}
