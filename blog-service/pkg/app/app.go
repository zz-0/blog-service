package app

import (
	"blog-service/pkg/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

//响应处理

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"pagesize"`
	TotalRows int `json:"total_rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	//返回响应的上下文
	return &Response{
		Ctx: ctx,
	}
}

//返回响应
func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

//错误响应
func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}
