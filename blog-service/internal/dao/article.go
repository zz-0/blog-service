package dao

import (
	"blog-service/internal/model"
	"blog-service/pkg/app"
)

type Article struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"create_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

//新建文章---利用model层中的函数在数据库中进行操作
func (d *Dao) CreateArticle(param *Article) (*model.Article, error) {
	article := model.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model: &model.Model{
			CreatedBy: param.CreatedBy,
		},
	}

	//返回model层   执行新建文章函数
	return article.Create(d.engine)
}

//(在数据库中)更新文章
func (d *Dao) UpdateArticle(param *Article) error {
	article := model.Article{
		Model: &model.Model{
			ID: param.ID,
		},
	}
	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}
	if param.Title != "" {
		values["title"] = param.Title
	}
	if param.CoverImageUrl != "" {
		values["cover_image_url"] = param.CoverImageUrl
	}
	if param.Desc != "" {
		values["desc"] = param.Desc
	}
	if param.Content != "" {
		values["content"] = param.Content
	}
	return article.Update(d.engine, values)
}

//(从数据库中)获取文章
func (d *Dao) GetArticle(id uint32, state uint8) (model.Article, error) {
	article := model.Article{
		Model: &model.Model{ID: id},
		State: state,
	}
	return article.Get(d.engine)
}

//(从数据库中)删除文章
func (d *Dao) DeleteArticle(id uint32) error {
	article := model.Article{
		Model: &model.Model{ID: id},
	}
	return article.Delete(d.engine)
}

//(从数据库中)统计文章-依据标签id
func (d *Dao) CountArticleListByTagID(id uint32, state uint8) (int, error) {
	article := model.Article{State: state}
	return article.CountByTagID(d.engine, id)
}

//(从数据库中)获取文章-依据标签id
func (d *Dao) GetArticleListByTagID(id uint32, state uint8, page, pageSize int) ([]*model.ArticleRow, error) {
	article := model.Article{State: state}
	return article.ListByTagID(d.engine, id, app.GetPageOffset(page, pageSize), pageSize)
}
