package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

//请求超时的中间件

func ContextTime(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(c.Request.Context(),t)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
