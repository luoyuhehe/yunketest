package thinkgo

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"github.com/sahara-go/thinkgo/event"
	"github.com/sahara-go/thinkgo/validate"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App interface {
	Run() (err error)
}

type HttpApp struct {
	engine               *gin.Engine
	httpServer           *http.Server
	RegisterRouter       func(engine *gin.Engine) (err error)
	RegisterMiddleware   func(engine *gin.Engine) (err error)
	RegisterEventHandler func(event event.Event) error
	StartUp              func(engine *gin.Engine) (err error)
}

const (
	defaultReadTimeout  int64 = 10
	defaultWriteTimeout int64 = 10
)

func (app *HttpApp) Run() {
	ctx := context.Background()
	// 监听配置文件
	app.WatchConfig()
	//初始化数据库
	if err := app.InitDB(); err != nil {
		panic(err)
	}
	// 初始化redis服务
	NewDefaultRedisClient()
	// 程序结束前关闭数据库链接
	defer app.deferHandle()
	// 创建gin.Engine
	if err := app.createEngine(); err != nil {
		panic(err)
	}
	//注册默认路由
	RegisterDefaultRouter(app.engine)
	//注册自定义路由
	if err := app.RegisterRouter(app.engine); err != nil {
		panic(fmt.Errorf("注册路由失败：%s", err))
	}
	GetLogger().Info("router register success")
	// 注册插件路由
	//InstallPlugs(Router)
	// 初始化server
	app.initHttpServer()
	// 注册数据验证引擎
	app.RegisterValidator()
	// Rbac权限初始化
	app.InitRbac()
	// 执行用户自定义初始化操作
	if err := app.StartUp(app.engine); err != nil {
		panic(fmt.Errorf("start up handle error：%s", err))
	}
	go GoroutineAttachPanicHandle(func() {
		// 服务连接
		if AppConfig.Server.IsHttps {
			// https
			if err := app.httpServer.ListenAndServeTLS(AppConfig.SSL.CertFile, AppConfig.SSL.KeyFile); err != nil && err != http.ErrServerClosed {
				panic(fmt.Errorf("start server error: %s", err))
			}
			return
		}

		// http
		if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("start server error: %s", err))
		}
	})

	GetLogger().Infof("server run success on %s:%d", AppConfig.Server.Addr, AppConfig.Server.Port)
	GetLogger().Infof("Server Run http://%s:%d/", AppConfig.Server.Addr, AppConfig.Server.Port)
	GetLogger().Infof("Swagger URL http://%s:%d/swagger/index.html", AppConfig.Server.Addr, AppConfig.Server.Port)
	GetLogger().Infof("Enter Control + C Shutdown Server")

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	GetLogger().Info("Shutdown Server")
	if err := app.httpServer.Shutdown(ctx); err != nil {
		panic(fmt.Errorf("server shutdown error:%s", err))
	}
}

func (app *HttpApp) RegisterValidator() {
	binding.Validator = validate.GetDefaultValidator()
}

// InitRbac 权限控制初始化
func (app *HttpApp) InitRbac() {
	if AppConfig.Rbac.Enable {
		createDefaultRbac()
	}
}

func (app *HttpApp) WatchConfig() {
	//监听变化
	appConfig.WatchConfig()
	appConfig.OnConfigChange(func() {
		if err := appConfig.Unmarshal(AppConfig); err != nil {
			GetLogger().Error(errors.WithMessage(err, "应用配置解析失败"))
		}

		if err := checkAppConfig(AppConfig); err != nil {
			GetLogger().Error(errors.WithMessage(err, "应用配置检查不通过"))
		}
	})
}

// initHttpServer 初始化httpServer
func (app *HttpApp) initHttpServer() {
	readTimeout := AppConfig.Server.ReadTimeOut
	if readTimeout == 0 {
		readTimeout = defaultReadTimeout
	}

	writeTimeout := AppConfig.Server.ReadTimeOut
	if writeTimeout == 0 {
		writeTimeout = defaultWriteTimeout
	}

	address := fmt.Sprintf("%s:%d", AppConfig.Server.Addr, AppConfig.Server.Port)
	app.httpServer = &http.Server{
		Addr:           address,
		Handler:        app.engine,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (app *HttpApp) createEngine() (err error) {
	app.engine = gin.New()
	//设置模式
	if AppConfig.IsProdMode() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	// 注册异常处理中间件
	app.engine.Use(recoveryMiddleware(true))
	// 注册日志中间件
	app.engine.Use(logMiddleware(time.RFC3339, true))
	// 注册ssl中间件
	if AppConfig.Server.IsHttps {
		app.engine.Use(tlsMiddleware()) // 打开就能玩https了
	}
	GetLogger().Debug("use middleware logger")
	// 注册跨域处理中间件
	app.engine.Use(corsMiddleware())
	if err := app.RegisterMiddleware(app.engine); err != nil {
		GetLogger().Errorf("注册中间件发生错误:%s", err)
		return err
	}
	// 注册权限控制中间件
	if AppConfig.Rbac.Enable {
		app.engine.Use(rbacMiddleware())
	}
	// 注册session中间件
	app.engine.Use(sessionMiddleware())
	GetLogger().Debug("use middleware cors")
	return nil
}

// InitDB 初始化数据库连接
func (app *HttpApp) InitDB() (err error) {
	createDatabaseConnPool()
	return nil
}

// deferHandle 退出时进行的操作
func (app *HttpApp) deferHandle() {
	//关闭数据库连接
	if err := DB.Close(); err != nil {
		GetLogger().Error(errors.WithMessage(err, "db close failed"))
	}
	//关闭session仓储
	if err := defaultSessionStore.Close(); err != nil {
		GetLogger().Error(errors.WithMessage(err, "default session store close failed"))
	}

	// 日志缓存刷回
	if err := GetLogger().Sync(); err != nil {
		GetLogger().Error(errors.WithMessage(err, "logger sync failed"))
	}

	if err := recover(); err != nil {
		GetLogger().Error(err)
	}
}
