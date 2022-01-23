package api

import (
	"blog-service/global"
	"blog-service/internal/service"
	"blog-service/pkg/app"
	"blog-service/pkg/errcode"

	"github.com/gin-gonic/gin"
)

//校验和获取入参后，在数据库中进行查询，检查认证信息是否存在
func GetAuth(c *gin.Context) {
	param := service.AuthRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.CheckAuth err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}

	//获取token
	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(c, "app.GenrateToken err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}

	//校验token
	response.ToResponse(gin.H{
		"token": token,
	})
}
