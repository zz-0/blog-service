package global

import (
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
)

//此文件是将pkg/setting下的配置信息与应用程序关联起来

//对最初预估的三个区段进行配置并声明全局变量
var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger //定义Logger对象
	JWTSetting      *setting.JWTSettingS
	EmailSetting    *setting.EmailSettingS
)

//邮箱配置项
type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	IsSSL    bool
	From     string
	To       []string
}
