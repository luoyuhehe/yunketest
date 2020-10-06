package thinkgo

import (
	"gitee.com/sahara-go/thinkgo/session"
	"gitee.com/sahara-go/thinkgo/session/session_store"
	"gitee.com/sahara-go/thinkgo/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var ContextSessionKey = utils.MD5V([]byte("thinkgo-session"))
var defaultSessionStore = newDefaultSessionStore()

// newDefaultSessionStore 创建默认的session store
func newDefaultSessionStore() session.Store {
	//判断是否开启session
	if !AppConfig.Session.Enable {
		//未开启session
		return nil
	}

	if AppConfig.Session.Store == session.RedisStore {
		store, err := session_store.NewRedisStore(
			AppConfig.Session.Redis.MaxIdleConn,
			AppConfig.Session.Redis.Host,
			AppConfig.Session.Redis.Password,
			[]byte(AppConfig.Session.Redis.GetKey()),
		)
		if err != nil {
			panic(errors.WithMessage(err, "redis session store 创建失败"))
		}
		return store
	}
	return nil
}

// Start 启动会话
func SessionStart(ctx *gin.Context) (session.Session, error) {
	return GetSession(ctx)
}

func GetSession(ctx *gin.Context) (session.Session, error) {
	s, err := session.GetSession(ctx, defaultSessionStore, AppConfig.Session.SessionIdKey)
	if err != nil {
		return nil, err
	}
	return s, nil
}
