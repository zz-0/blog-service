package dao

import "blog-service/internal/model"

//鉴权操作
//传入密匙和具体的信息，之后获取到数据库中的结构体
func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
	auth := model.Auth{AppKey: appKey, AppSecret: appSecret}
	return auth.Get(d.engine)
}