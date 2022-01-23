package routers

import (
	"blog-service/global"
	"blog-service/internal/middleware"
	"blog-service/internal/routers/api"
	v1 "blog-service/internal/routers/api/v1"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-programming-tour-book/blog-service/docs"
	"github.com/go-programming-tour-book/blog-service/pkg/limiter"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

//路由初始化
func NewRouter() *gin.Engine {
	r := gin.New()

	//使用自己定义的中间件，不使用gin自带的
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	//各种中间件的注册
	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTime(60 * time.Second))
	r.Use(middleware.Translations())
	//r.Use()
	r.Use(middleware.Translations())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	article := v1.NewArticle() //获取文章结构体
	tag := v1.NewTag()         //获取标签结构体
	upload := api.NewUpload()  //获取静态文件地址

	r.POST("/upload/file", upload.UploadFile)                         //上传图片到数据库
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath)) //静态文件的存放
	r.GET("/auth", api.GetAuth)

	//业务接口的基础编码
	apivl := r.Group("/api/v1")
	apivl.Use() //middleware.JWT()
	{
		apivl.POST("/tags", tag.Create)
		apivl.DELETE("/tags/:id", tag.Delete)
		apivl.PUT("/tags/:id", tag.Update)
		apivl.PATCH("/tags/:id/state", tag.Update)
		apivl.GET("/tags", tag.List)

		apivl.POST("/articles", article.Create)
		apivl.DELETE("/articles/:id", article.Delete)
		apivl.PUT("/articles/:id", article.Update)
		apivl.PATCH("/article/:id/state", article.Update)
		apivl.GET("/article", article.List)
	}
	return r
}
