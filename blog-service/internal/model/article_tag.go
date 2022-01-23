package model

import "github.com/jinzhu/gorm"

//文章标签关联表的结构体
type ArticleTag struct {
	*Model           //公共字段
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_artcle_tag"
}

//依靠ID获取文章标签(文章标签关联模型)
func (a ArticleTag) GetByID(db *gorm.DB) (ArticleTag, error) {
	var articleTag ArticleTag
	err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).First(&articleTag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return articleTag, err
	}

	return articleTag, nil
}

//根据标签id来获取文章标签列表--数组
func (a ArticleTag) ListByTID(db *gorm.DB) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	if err := db.Where("tag_id = ? AND is_del = ?", a.TagID, 0).Find(&articleTags).Error; err != nil {
		return nil, err
	}
	return articleTags, nil
}

//根据文章id来获取文章标签列表
func (a ArticleTag) ListByAIDs(db *gorm.DB, ArticleIDs []uint32) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	err := db.Where("article_id IN (?) AND is_del = ?", ArticleIDs, 0).Find(&articleTags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articleTags, nil
}

//新建文章标签
func (a ArticleTag) Create(db *gorm.DB) error {
	if err := db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

//更新
func (a ArticleTag) UpdateOne(db *gorm.DB, values interface{}) error {
	if err := db.Model(&a).Where("article_id = ? AND is_del  = ?", a.ArticleID, 0).Limit(1).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

//依据id主键删除
func (a ArticleTag) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}

	return nil
}

//依据文章id删除
func (a ArticleTag) DeleteOne(db *gorm.DB) error {
	if err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).Delete(&a).Limit(1).Error; err != nil {
		return err
	}

	return nil
}
