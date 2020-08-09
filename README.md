# thinkgo

thinkgo 是快速开发框架，包含了session、日志、配置管理、rbac鉴权、中间件和验证码等常用组件的封装


## Quick Start

#### Create `hello` directory, cd `hello` directory

    mkdir hello
    cd hello
 
#### Init module

    go mod init

#### Download and install

    go get github.com/sahara-gopher/thinkgo
    
#### Create file `config/app.yaml`

```
# 应用配置
app_name: 'gofast-admin'
base_router_group: 'api'
auto_log_resp: 'all' # all|error|success

# 数据库连接配置
database:
  type: 'mysql'
  mysql:
    username: root
    password: '123456'
    host: 'xxx.xxx.xxx.xxx'
    port: '3366'
    db_name: 'gofast_admin'
    charset: 'utf8'
    idle_num: 10
    pool_num: 10
    log_mode: true
    multi_statements: true

# redis配置
redis:
  host: 'xxx.xxx.xxx.xxx:6379'
  password: ''
  db: 0
  pool_num: 10
  idel_timeout: 30

# server配置
server:
  use_multipoint: false
  env: 'public'  # Change to "develop" to skip authentication for development mode
  addr: '127.0.0.1'
  port: 8888
  read_timeout: 10
  write_timeout: 10

session:
  enable: true
  disable_cookie: false #是否禁用从cookie获取sessionid
  sessionid_key: "thinkgo_sessionid" #sessionid键名
  store: 'redis'
  redis:
    host: 'xxx.xxx.xxx.xxx:6379'
    max_idle_conn: 10
    password: ''
    key: ''

# captcha配置
captcha:
  store: 'redis'
  redis:
    host: 'xxx.xxx.xxx.xxx:6379'
    password: ''
    db: 0
    pool_num: 10
    idel_timeout: 30
    key_prefix: 'captcha:code:'
    expire: 86400


#日志配置
log:
  file:
    prefix: ''
    enable: true
    path: '/var/log/gofast-admin/'
    level: 'debug'
    max_size: 100 # 单位是MB
    max_backups: 100
    max_age: 30
    compress: false
  console:
    prefix: ''
    enable: true
    level: 'debug'
  kafka:
    prefix: ''
    enable: false
    level: 'debug'
    addr:
      - '127.0.0.1'
    topic: ''

#ssl 配置
ssl:
  key_file: ''
  cert_file: ''

rbac:
  enable: false

```

#### Create file `hello.go`
```go
package main

import "github.com/sahara-gopher/thinkgo"

func main(){
    package main

import (
	"gitee.com/luoyusnnu/thinkgo"
	"github.com/gin-gonic/gin"
	"gofast-admin/router"
)

func main() {
	//创建http应用
	httpApp := &thinkgo.HttpApp{
		RegisterRouter: func(engine *gin.Engine) (err error) {
			router.RegisterStaticFileRouter(engine)
			// 方便统一添加路由组前缀 多服务器上线使用
			BaseGroup := engine.Group(thinkgo.AppConfig.BaseRouterGroup)
			//router.POST("register", v1.Register)
			BaseGroup.Group("admin").POST("user/login", func(ctx *gin.Context) {
			})
			thinkgo.GetLogger().Info("router register success")
			return nil
		},
		RegisterMiddleware: func(engine *gin.Engine) (err error) {
			return nil
		},
		StartUp: func(engine *gin.Engine) (err error) {
			return nil
		},
	}
	//启动http应用
	httpApp.Run()
}
```
#### Build and run

    go build hello.go
    ./hello

#### Go to [http://localhost:8888](http://localhost:8888)
