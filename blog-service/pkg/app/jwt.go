package app

/*鉴权的过程：
1.用户端输入账号密码传入服务端
2.服务端生成tkoen并返回给用户端
3.用户端发送请求(携带token)给服务端
4.服务端进行验证(验证签名)
5.返回用户端请求的数据
*/
import (
	"blog-service/global"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-programming-tour-book/blog-service/pkg/util"
)

type Claims struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	jwt.StandardClaims
}

//获取该项目中的jwt secret(yaml文件中的secret，密匙，不能随便暴露)
func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

//传入key和secret，生成jwt token
func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire)
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

//服务端对token进行解析和校验
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return GetJWTSecret(), nil
		})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
