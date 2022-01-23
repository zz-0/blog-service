package service

import (
	"blog-service/internal/model"
	"blog-service/pkg/app"
)

//为业务接口中定义的增删改查和统计行为编写的Request结构体

//form---表单的映射字段名   binding----入参校验的规则内容

//处理标签模块的业务逻辑
type CountTagRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type TagListRequest struct {
	Name     string `form:"name" binding:"max=100"`
	//CreateBy string `form:"created_by" binding:"required,min=2.max=100"`
	State    uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreatTagRequest struct {
	Name      string `form:"name" binding:"required,min=3,max=100"`
	CreatedBy string `form:"created_by" binding:"reuired,min=3,max=100"`
	State     uint8  `from:"State,default=1" binding:"oneof=0 1"`
}

type UpdateTagRequest struct {
	ID         uint32 `form:"id" binding:"required,gte=1"`
	Name       string `form:"name" binding:"required,min=3,max=100"`
	State      uint8  `form:"State,default=1" binding:"oneof=0 1"`
	ModifiedBy string `form:"modified_by" bind:"requires,min=2.max=100"`
}

type DeleteTagRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

//定义了request结构体作为接口入参的基准

//获取标签的数量
func (svc *Service) CountTag(param *CountTagRequest) (int, error) {
	return svc.dao.CountTag(param.Name, param.State)
}

//获取标签的详细信息
func (svc *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
	return svc.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}

//新建标签
func (svc *Service) CreateTag(param *CreatTagRequest) error {
	return svc.dao.CreateTag(param.Name, param.State, param.CreatedBy)
}

//更新标签
func (svc *Service) UpdateTag(param *UpdateTagRequest) error {
	return svc.dao.UpdateTag(param.ID, param.Name, param.State, param.ModifiedBy)
}

//删除标签
func (svc *Service) DeleteTag(param *DeleteTagRequest) error {
	return svc.dao.DeleteTag(param.ID)
}
