package middleware

// 拦截器
//func CasbinHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		log.Debug("casbin middleware load")
//		// 获取请求的URI
//		obj := c.Request.URL.RequestURI()
//		// 获取请求方法
//		act := c.Request.Method
//		// 获取用户的角色
//		sub := "waitUse.AuthorityId"
//		e := casbin.Casbin(thinkgo.DB)
//		// 判断策略中是否存在
//		if config.AppConfig.Server.Env == "develop" || e.Enforce(sub, obj, act) {
//			c.Next()
//			return
//		}
//
//		response.New(c).Json(-1, gin.H{}, "权限不足")
//		c.Abort()
//		return
//	}
//}
