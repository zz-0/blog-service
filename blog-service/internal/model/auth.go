package model

import "github.com/jinzhu/gorm"

//token组成:header+playload+signature
//header:令牌的类型和签名算法---转换成json
//playload:实际传输的数据，承载的具体信息，不建议放入敏感信息，因为此数据可以被前端获取
//signature:由前两部分，利用base64url进行编码，用于验证信息是否被篡改
//完成上述的三步后，再通过header里标注的加密算法对全部的数据进行加密，并且插入规定的密匙

//鉴权作用的结构体，token里的数据
type Auth struct {
	*Model
	AppKey    string `json:"app_key"`    //
	AppSecret string `json:"app_secret"` //密匙
}

func (a Auth) TableName() string {
	return "blog_auth"
}

//此函数用于查询判断客户端传入的信息中的数据是否存在,返回auth表中的鉴权数据(传入dao层中)
func (a Auth) Get(db *gorm.DB) (Auth, error) {
	var auth Auth
	db = db.Where("app_key = ? AND app_secret = ? AND is_del = ?", a.AppKey, a.AppSecret, 0)

	err := db.First(&auth).Error
	//此错误捕获有问题

	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	return auth, nil
}
