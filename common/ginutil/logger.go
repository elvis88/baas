package ginutil

import (
	"time"

	"github.com/gin-gonic/gin"
)

// UseLogger logger
func UseLogger(router *gin.Engine, print func(format string, a ...interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		//开始时间
		start := time.Now()
		//处理请求
		c.Next()
		//结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)
		//path
		path := c.Request.URL.Path
		//ip
		clientIP := c.ClientIP()
		//方法
		method := c.Request.Method
		//状态
		statusCode := c.Writer.Status()
		print("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
