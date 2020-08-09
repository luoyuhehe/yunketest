package thinkgo

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// corsMiddleware 处理跨域请求,支持options访问
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		GetLogger().Debug("cors middleware load")
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// Recovery returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func recoveryMiddleware(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recovery middleware recovers from any panics and writes a 500 if there was one.
		//appConfigS.Engine.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		//	if err, ok := recovered.(string); ok {
		//		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		//	}
		//	c.AbortWithStatus(http.StatusInternalServerError)
		//}))
		GetLogger().Debug("recovery middleware load")
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					GetLogger().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					//c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					GetLogger().Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					GetLogger().Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// sessionMiddleware 会话中间件
func sessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		GetLogger().Debug("session middleware load")
		//判断是否开启session
		if !AppConfig.Session.Enable {
			//未开启session
			return
		}

		if _, err := SessionStart(c); err != nil {
			log.Panic(errors.Wrap(err, "session中间件启动发生错误：%s"))
		}
		c.Next()
	}
}

// tlsMiddleware tls中间件
func tlsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		GetLogger().Debug("tls middleware load")
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			GetLogger().Error(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}

// logMiddleware 接收gin框架默认的日志
func logMiddleware(timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		//appConfigS.Engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		//	// your custom format
		//	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		//		param.ClientIP,
		//		param.TimeStamp.Format(time.RFC1123),
		//		param.Method,
		//		param.Path,
		//		param.Request.Proto,
		//		param.StatusCode,
		//		param.Latency,
		//		param.Request.UserAgent(),
		//		param.ErrorMessage,
		//	)
		//}))
		// Recovery middleware recovers from any panics and writes a 500 if there was one.
		GetLogger().Debug("log middleware load")
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				GetLogger().Error(e)
			}
		} else {
			GetLogger().Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
			)
		}
	}
}

// rbacMiddleware 权限控制中间件
func rbacMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		GetLogger().Debug("rbac middleware load")
		c.Next()
	}
}
