package dao

import "blog-service/internal/model"

//编写对文章标签的增删改查操作

//根据文章id(从数据库中)获取文章标签
func (d *Dao) GetArticleTagByAID(articleID uint32) (model.ArticleTag, error) {
	//取出文章id为article的文章标签
	articleTag := model.ArticleTag{ArticleID: articleID}
	return articleTag.GetByID(d.engine)
}

//(从数据库中)根据标签id获取文章标签
func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{TagID: tagID}
	return articleTag.ListByTID(d.engine)
}

//根据标签id获取文章标签列表
func (d *Dao) GetArticleTagListByAIDs(articleIDs []uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{}
	return articleTag.ListByAIDs(d.engine, articleIDs)
}

//(从数据库中)新建文章标签
func (d *Dao) CreateArticleTag(articleID, tagID uint32, createdBy string) error {
	articleTag := model.ArticleTag{
		Model: &model.Model{
			CreatedBy: createdBy,
		},
		ArticleID: articleID,
		TagID:     tagID,
	}

	//返回数据库中创建新的文章标签
	return articleTag.Create(d.engine)
}

//(从数据库中)更新文章标签
func (d *Dao) UpdateArticleTag(articleID, tagID uint32, modifiedBy string) error {
	articleTag := model.ArticleTag{ArticleID: articleID}
	values := map[string]interface{}{
		"article_id":  articleID,
		"tag_id":      tagID,
		"modified_by": modifiedBy,
	}
	return articleTag.UpdateOne(d.engine, values)
}

//(从数据库中)删除文章标签
func (d *Dao) DeleteArticleTag(articleID uint32) error {
	articleTag := model.ArticleTag{ArticleID: articleID}
	return articleTag.DeleteOne(d.engine)
}
