package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"tic3001-go-server/config"
)

var (
	Client *redis.Client
)

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr: config.Config.MustString("redis.cache.host", "127.0.0.1") + ":" +
			config.Config.MustString("redis.cache.port", "5379"),
		Password: "",
		DB:       0,
	})
	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("fail to init redis:%s", err.Error())
	}
}
