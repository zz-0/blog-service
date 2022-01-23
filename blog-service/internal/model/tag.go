package model

import (
	"blog-service/pkg/app"

	"github.com/jinzhu/gorm"
)

//标签表的字段
type Tag struct {
	*Model        //公用部分
	Name   string `json:"name"`  //标签名称
	State  uint8  `json:"state"` //状态管理标识
}

//返回表的标签表名
func (a Tag) TableName() string {
	return "blog_tag"
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

//统计标签数量
func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

//标签列表
func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	//查找语句
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

//按照ID排列
func (t Tag) ListByIDs(db *gorm.DB, ids []uint32) ([]*Tag, error) {
	var tags []*Tag
	db = db.Where("state = ? AND is_del = ?", t.State, 0)
	err := db.Where("id IN (?)", ids).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

//(从数据库中)获取标签
func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	err := db.Where("id = ? AND is_del = ? AND state = ?", t.ID, 0, t.State).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}

	return tag, nil
}

//(从数据库中)新建标签
func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

//(从数据库中)更新标签
func (t Tag) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Updates(values).Where("id = ? AND is_del = ?", t.ID).Error; err != nil {
		return err
	}
	return nil
}

//(从数据库中)删除标签
func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}
