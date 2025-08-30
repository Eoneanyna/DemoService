package data

import (
	"demoserveice/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedis(c *conf.Data) (*Redis, func(), error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     30,
		PoolTimeout:  30 * time.Second,
		MinIdleConns: 10,
		DB:           11,
	})
	//初始化协程池
	d := &Redis{
		rdb: RedisClient,
	}
	cleanup := func() {
		log.Info("message", "closing the redis resources")
		if err := d.rdb.Close(); err != nil {
			log.Error(err)
		}
	}

	return d, cleanup, nil
}
