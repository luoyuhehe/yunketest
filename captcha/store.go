package captcha

import (
	"fmt"
	"gitee.com/sahara-gopher/thinkgo"
	"gitee.com/sahara-gopher/thinkgo/log"
	"gitee.com/sahara-gopher/thinkgo/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/mojocn/base64Captcha"
	"github.com/pkg/errors"
	"sync"
	"time"
)

const (
	SessionStoreType = "session"
	RedisStoreType   = "redis"
)

var initRedisStoreOnce sync.Once
var defaultCaptchaStore base64Captcha.Store

type RedisStoreOptions struct {
	Host      string
	Password  string
	DB        int
	KeyPrefix string
	expire    int
	PoolNum   int
	Timeout   int
}

type RedisStore struct {
	engine  *redis.Pool
	options *RedisStoreOptions
}

func getDefaultCaptchaStore() (base64Captcha.Store, error) {
	c := thinkgo.AppConfig.Captcha
	if !utils.InStringArray(c.Store, SessionStoreType, RedisStoreType) {
		return nil, errors.New(fmt.Sprintf("store类型%s不合法", c.Store))
	}
	initRedisStoreOnce.Do(func() {
		if c.Store == RedisStoreType {
			store, err := NewRedisStore(&RedisStoreOptions{
				Host:      c.Redis.Host,
				Password:  c.Redis.Password,
				DB:        c.Redis.DB,
				KeyPrefix: c.Redis.KeyPrefix,
				expire:    c.Redis.Expire,
				PoolNum:   c.Redis.PoolNum,
				Timeout:   c.Redis.IdleTimeout,
			})
			if err != nil {
				panic(err)
			}
			defaultCaptchaStore = store
		}
	})
	return defaultCaptchaStore, nil
}

func NewRedisStore(options *RedisStoreOptions) (RedisStore, error) {
	pool := &redis.Pool{
		MaxIdle:     options.PoolNum,
		MaxActive:   options.PoolNum,
		IdleTimeout: time.Duration(options.Timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", options.Host)
			if err != nil {
				return nil, err
			}
			if options.Password != "" {
				if _, err := c.Do("AUTH", options.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			rs, err := c.Do("PING")
			log.Infof("redis ping :%+v, error:%+v", rs, err)
			return err
		},
	}

	return RedisStore{engine: pool, options: options}, nil
}

func (s RedisStore) getKey(id string) string {
	return fmt.Sprintf("%s_%s", s.options.KeyPrefix, id)
}

// Set sets the digits for the captcha id.
func (s RedisStore) Set(id string, value string) {
	conn := s.engine.Get()
	if _, err := conn.Do("SETEX", s.getKey(id), s.options.expire, value); err != nil {
		panic(errors.WithMessage(err, "验证码设置失败"))
	}
	defer func() {
		_ = conn.Close()
	}()
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (s RedisStore) Get(id string, clear bool) string {
	conn := s.engine.Get()
	defer func() {
		_ = conn.Close()
	}()
	value, err := redis.String(conn.Do("GET", s.getKey(id)))
	if err == redis.ErrNil {
		return ""
	} else if err != nil {
		panic(err)
	}

	if clear {
		_, _ = s.engine.Get().Do("DEL", s.getKey(id))
	}
	return value
}

//Verify captcha's answer directly
func (s RedisStore) Verify(id string, answer string, clear bool) bool {
	conn := s.engine.Get()
	defer func() {
		_ = conn.Close()
	}()

	value, err := redis.String(conn.Do("GET", s.getKey(id)))
	log.Infof("captcha verify, value:%+v, error:%+v", value, err)
	var rs = false
	if err != nil {
		if err != redis.ErrNil {
			panic(err)
		}
	} else if value == answer {
		rs = true
	}

	//判断是否删除
	if clear {
		_, _ = s.engine.Get().Do("DEL", s.getKey(id))
	}
	return rs
}
