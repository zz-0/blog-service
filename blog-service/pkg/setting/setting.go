package setting

import "github.com/spf13/viper"

//此模块作用是对服务端的配置进行设置

//viper模块解析yaml文件
type Setting struct {
	vp *viper.Viper
}

//初始化本项目的配置的基础属性
func NewSetting() (*Setting, error) {
	vp := viper.New()

	//1.读取设定配置文件的名称为config
	vp.SetConfigName("config")

	//2.添加读取设置其配置路径相对路径为configs/   注:viper是支持配置多个路径的
	vp.AddConfigPath("configs/")

	//3.配置类型为yaml
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp}, nil
}



 

  