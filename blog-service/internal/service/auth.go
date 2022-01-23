package service

import "errors"

//jwt的入参校验
type AuthRequest struct {
	AppKey    string `form:"app_key" binding:"required"`
	AppSecret string `form:"app_secret" binding:"required"`
}

//JWT，从路由中获取密匙并加以验证
func (svc *Service) CheckAuth(param *AuthRequest) error {
	auth, err := svc.dao.GetAuth(
		param.AppKey,
		param.AppSecret, 
	)
	if err != nil {
		return err
	}

	if auth.ID > 0 {
		return nil
	}

	return errors.New("auth info does not exist.")
}
