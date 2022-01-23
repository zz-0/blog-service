package global

import "github.com/jinzhu/gorm"

//将数据库连接设置为全局变量
var (
	DBEngine *gorm.DB
)
