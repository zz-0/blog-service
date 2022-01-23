package middleware

import (
	"blog-service/pkg/app"
	"blog-service/pkg/errcode"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//校验token
func JWT() gin.HandlerFunc {
	//返回值
	return func(c *gin.Context) {
		var (
			token string  
			ecode = errcode.Success
		)

		//获取token
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}

		//判断token
		if token == "" {
			ecode = errcode.InvalidParams
		} else {
			_, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			}
		}
		if ecode != errcode.Success {
			response := app.NewResponse(c)
			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		c.Next()
	}
}
