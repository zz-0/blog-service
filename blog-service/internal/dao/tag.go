package dao

import (
	"blog-service/internal/model"
	"blog-service/pkg/app"
)

//在dao层进行了数据访问对象的封装,并对业务所需的字段进行处理

//(从数据库中)统计标签数量
func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.Tag{
		Name:  name,
		State: state,
	}
	return tag.Count(d.engine)
}

//(从数据库中)获取标签
func (d *Dao) GetTag(id uint32, state uint8) (model.Tag, error) {
	tag := model.Tag{Model: &model.Model{ID: id}, State: state}
	return tag.Get(d.engine)
}

//(从数据库中)获取标签列表
func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.Tag, error) {
	tag := model.Tag{Name: name, State: state}
	pageOffset := app.GetPageOffset(page, pageSize)
	return tag.List(d.engine, pageOffset, pageSize)
}

//(从数据库中)新建标签
func (d *Dao) CreateTag(name string, state uint8, createdBy string) error {
	tag := model.Tag{
		Name:  name,
		State: state,
		Model: &model.Model{
			CreatedBy: createdBy,
		},
	}
	return tag.Create(d.engine)
}

//(从数据库中)更新标签
func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
	tag := model.Tag{
		Model: &model.Model{
			ID: id,
		},
	}
	values := map[string]interface{}{
		"state":       state,
		"modified_by": modifiedBy,
	}
	if name != "" {
		values["name"] = name
	}

	return tag.Update(d.engine, values)
}

//(从数据库中)删除标签
func (d *Dao) DeleteTag(id uint32) error {
	tag := model.Tag{Model: &model.Model{ID: id}}
	return tag.Delete(d.engine)
}
