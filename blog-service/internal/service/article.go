package service

import (
	"blog-service/internal/dao"
	"blog-service/internal/model"
	"blog-service/pkg/app"
)

//创建文章结构体
type Article struct {
	ID            uint32     `json:"id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	Content       string     `json:"content"`
	CoverImageUrl string     `json:"cover_image_url"`
	State         uint8      `json:"state"`
	Tag           *model.Tag `json:"tag"`
}

//文章请求结构体(入参校验)
type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

//文章列表(入参校验)
type ArticleListRequest struct {
	TagID uint32 `form:"tag_id" binding:"gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

//新建标签请求
type CreateArticleRequest struct {
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"required,min=2,max=100"`
	Desc          string `form:"desc" binding:"required,min=2,max=255"`
	Content       string `form:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"required,url"`
	CreatedBy     string `form:"created_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

//更新标签校验
type UpdateArticleRequest struct {
	ID            uint32 `form:"id" binding:"required,gte=1"`
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"min=2,max=100"`
	Desc          string `form:"desc" binding:"min=2,max=255"`
	Content       string `form:"content" binding:"min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"url"`
	ModifiedBy    string `form:"modified_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

//删除文章参数
type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

//获取单个文章---传入文章请求，获取文章的指针
func (svc *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	//param包含从路由传入的文章id
	//从dao层获取文章(依据传入的id)
	article, err := svc.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}

	//文章标签
	articleTag, err := svc.dao.GetArticleTagByAID(article.ID)
	if err != nil {
		return nil, err
	}

	//标签
	tag, err := svc.dao.GetTag(articleTag.TagID, model.STATE_OPEN)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		State:         article.State,
		Tag:           &tag,
	}, nil
}

//获取文章列表
func (svc *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {
	//统计文章
	articleCount, err := svc.dao.CountArticleListByTagID(param.TagID, param.State)
	if err != nil {
		return nil, 0, err
	}

	//利用dao层获取文章
	articles, err := svc.dao.GetArticleListByTagID(param.TagID, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var articleList []*Article
	//遍历文章
	for _, article := range articles {
		articleList = append(articleList, &Article{
			ID:            article.ArticleID,
			Title:         article.ArticletItle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: article.CoverImageUrl,
			Tag: &model.Tag{
				Model: &model.Model{ID: article.TagID},
				Name:  article.TagName,
			},
		})
	}

	return articleList, articleCount, nil
}

//创建文章
func (svc *Service) CreateArticle(param *CreateArticleRequest) error {
	//先将输入的param赋给article
	article, err := svc.dao.CreateArticle(&dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		CreatedBy:     param.CreatedBy,
	})
	if err != nil {
		return err
	}

	//再利用获取到的article的内容创建新的文章
	err = svc.dao.CreateArticleTag(article.ID, param.TagID, param.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

//操作与创建文章相似
//更新文章
func (svc *Service) UpdateArticle(param *UpdateArticleRequest) error {
	err := svc.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		ModifiedBy:    param.ModifiedBy,
	})
	if err != nil {
		return nil
	}

	err = svc.dao.UpdateArticleTag(param.ID, param.TagID, param.ModifiedBy)
	if err != nil {
		return nil
	}

	return nil
}

//删除文章
func (svc *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := svc.dao.DeleteArticle(param.ID)
	if err != nil {
		return err
	}

	return nil
}
