package thinkgo

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

func NewDefaultRedisClient() {
	redisCfg := AppConfig.Redis
	RedisPool = &redis.Pool{
		MaxIdle:     redisCfg.PoolNum,
		MaxActive:   redisCfg.PoolNum,
		IdleTimeout: time.Duration(redisCfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisCfg.Host)
			if err != nil {
				return nil, err
			}
			if redisCfg.Password != "" {
				if _, err := c.Do("AUTH", redisCfg.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
