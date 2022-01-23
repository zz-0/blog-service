package service

//service层:主要的业务逻辑层，传入参数(至dao层)

import (
	"blog-service/global"
	"blog-service/internal/dao"
	"context"
)

//将dao层实例化
type Service struct {
	ctx context.Context
	dao *dao.Dao
}

//将service实例化
func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New( global.DBEngine)
	return svc
}
