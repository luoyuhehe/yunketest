package thinkgo

import (
	"gitee.com/luoyusnnu/thinkgo/log"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// RegisterDefaultRouter 注册路由
func RegisterDefaultRouter(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Debug("register swagger handler")
}
