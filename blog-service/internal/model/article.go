package model

import (
	"blog-service/pkg/app"

	"github.com/jinzhu/gorm"
)

//文章表结构体
type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

//返回表名
func (a Article) TableName() string {
	return "blog_article"
}

//新建文章
func (a Article) Create(db *gorm.DB) (*Article, error) {
	//传入的是数据库的连接，获取到文章的地址
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

//更新文章
func (a Article) Update(db *gorm.DB, values interface{}) error {
	//在数据库中进行操作
	if err := db.Model(&a).Updates(values).Where("id = ? AND is_del = ?", a.ID).Error; err != nil {
		return err
	}
	return nil
}

//获取文章
func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND state = ? AND is_del = ?", a.ID, a.State, 0)

	//first函数:主键作为第一条件来查找
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}
	return article, nil
}

//删除文章
func (a Article) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}

//以下的关联查询时要与文章标签表进行关联

//关联查询的实现
type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticletItle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

//文章列表的查询(依据标签id)
func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error) {
	//新建切片储存文章的信息
	fields := []string{
		"ar.id AS article_id",
		"ar.title AS article_title",
		"ar.desc AS article_desc",
		"ar.cover_image_url",
		"ar.content",
	}
	//将标签的内容新增到切片中
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)

	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	//两表关联语句
	rows, err := db.Select(fields).
		Table(ArticleTag{}.TableName()+" AS at").                          //查询指定数据库的字段
		Joins("LEFT JION`"+Tag{}.TableName()+"`AS t ON at.tag_id = t.id"). //指定关联语句
		Joins("LEFT JION`"+Article{}.TableName()+"`AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleID, &r.ArticletItle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagID, &r.TagName); err != nil {
			return nil, err
		}
		articles = append(articles, r)
	}
	return articles, nil
}

//文章列表总数的查询方法依据标签id
func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JION `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JION `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at. `tag_id`=? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
