package dao

import "github.com/jinzhu/gorm"

//dao层:对数据进行操作的代码(从service传入的数据和从model层(数据库关联)取出的数据)，其主要功能是获取模型层中的数据以及对传入参数的处理

//连接池----将数据库中的数据实例化
type Dao struct {
	engine *gorm.DB
}

func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}

