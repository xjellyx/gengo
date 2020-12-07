package controller

import (
	"fmt"
)

var GinTemplate = fmt.Sprintf(`package api_{{.Package}}
{{- $Package := .Package }}
import(
	"github.com/gin-gonic/gin"
	"{{.Mod}}/service/{{$Package}}"
	"{{.Mod}}/model/{{$Package}}"
	"{{.Mod}}/controller/response"
)

{{$StructName := .StructName}}

// Ctrl{{$StructName}} 
type Ctrl{{$StructName}} struct {}

// AddOne add {{$StructName}} one record
// @tags {{$StructName}}
// @Summary add {{$StructName}} one record
// @Description add {{$StructName}} one record
// @Accept json
// @Produce json
// @Param {} body srv_{{$Package}}.Add{{$StructName}}ReqForm true "添加{{$StructName}}表单" 
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Package}}/add [post]
func (ct *Ctrl{{$StructName}}) AddOne(c *gin.Context) {
	var(
		req = &srv_{{$Package}}.Add{{$StructName}}ReqForm{}
		data interface{}
		code = response.CodeFail
		err error
	)
	
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else {
			response.NewGinResponse(c).Success(data).Response()
		}
	}()

	if err = c.ShouldBind(req);err!=nil{
		return
	}
	if data, err = srv_{{$Package}}.Add{{$StructName}}One(req); err != nil {
		return
	}
	
}

// AddList add {{$StructName}} list record
// @tags {{$StructName}}
// @Summary add {{$StructName}} list record
// @Description add {{$StructName}} list record
// @Accept json
// @Produce json
// @Param  {} body srv_{{$Package}}.{{$StructName}}BatchForm true "添加{{$StructName}}表单列表" 
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Package}}/addList [post]
func (ct *Ctrl{{$StructName}}) AddList(c *gin.Context) {
	var(
	data interface{}
	code = response.CodeFail
	req srv_{{$Package}}.{{$StructName}}BatchForm
	err error)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()

	if err = c.ShouldBind(&req);err!=nil{
		return
	}
	
	if data,err = srv_{{$Package}}.Add{{$StructName}}Batch(req);err!=nil{
		return
	}
}

// Edit edit {{$StructName}} one record
// @tags {{$StructName}}
// @Summary edit {{$StructName}} one record
// @Description edit {{$StructName}} one record
// @Accept json
// @Produce json
// @Param  {} body srv_{{$Package}}.Edit{{$StructName}}ReqForm true "编辑{{$StructName}}表单" 
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Package}}/edit [put]
func (ct *Ctrl{{$StructName}}) Edit(c *gin.Context) {
	var(
		data interface{}
		req = new(srv_{{$Package}}.Edit{{$StructName}}ReqForm)
		err error
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	if err = c.ShouldBind(&req);err!=nil{
		return
	}
	if data,err = srv_{{$Package}}.Edit{{.StructName}}One(req);err!=nil{
		return
	}
}

// GetOne get {{$StructName}} one record
// @tags {{$StructName}}
// @Summary get {{$StructName}} one record
// @Description get {{$StructName}} one record
// @Accept json
// @Produce json
// @Param id query string true "{{$StructName}} ID"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Package}}/get  [get]
func (ct *Ctrl{{$StructName}}) GetOne(c *gin.Context) {
	var(
		data interface{}
		id string
		err error
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	id = c.Query("id")
	if data,err = srv_{{$Package}}.Get{{$StructName}}One(id);err!=nil{
		return
	}
}

// GetList get {{$StructName}} list record
// @tags {{$StructName}}
// @Summary get {{$StructName}} list record
// @Description get {{$StructName}} list record
// @Accept json
// @Produce json
// @Param {} body model_{{$Package}}.Query{{$StructName}}Form true "获取{{$StructName}}列表form"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Package}}/list  [get]
func (ct *Ctrl{{$StructName}}) GetList(c *gin.Context) {
	var(
		data interface{}
		req = new(model_{{$Package}}.Query{{$StructName}}Form)
		err error
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	if err = c.ShouldBindQuery(req);err!=nil{
		return
	}
	if data,err = srv_{{$Package}}.Get{{$StructName}}Page(req);err!=nil{return}
}

// DeleteOne delete {{$StructName}} one record
// @tags {{$StructName}}
// @Summary delete {{$StructName}} one record
// @Description delete {{$StructName}} one record
// @Accept json
// @Produce json
// @Param id param string true "{{$StructName}} ID"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router  /api/v1/{{$Package}}/delete [delete]
func (ct *Ctrl{{$StructName}}) DeleteOne(c *gin.Context) {
	var(
		data interface{}
		err error
		id string
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	id = c.Param("id")
	if err = srv_{{$Package}}.Delete{{$StructName}}One(id);err!=nil{return}
}

// DeleteList delete {{$StructName}} list record
// @tags {{$StructName}}
// @Summary delete {{$StructName}} list record
// @Description delete {{$StructName}} list record
// @Accept json
// @Produce json
// @Param ids param [int] true "{{$StructName}} ID list"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router  /api/v1/{{$Package}}/deleteList [delete]
func (ct *Ctrl{{$StructName}}) DeleteList(c *gin.Context) {
	var(
		data interface{}
		ids []string
		err error
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	ids = c.PostFormArray("ids")
	if err = srv_{{$Package}}.Delete{{$StructName}}Batch(ids);err!=nil{return}
}

`)

var (
	ResponseTemplate = fmt.Sprintf(`
package response

import(	"github.com/gin-gonic/gin")

const(
	CodeSuccess = 0
	CodeFail = 4000
)
type Gin struct {
	c      *gin.Context
	resp   *Response
	status int
}

type Response struct {
	Meta Meta        %sjson:"meta"%s
	Data interface{} %sjson:"data"%s
}

type Meta struct {
	Code    int    %sjson:"code"%s
	Message string %sjson:"message"%s
}

// NewGinResponse
func NewGinResponse(c *gin.Context) *Gin {
	return &Gin{
		c,
		&Response{},
		200,
	}
}

func (g *Gin) Fail(code int, message string) *Gin {
	g.resp.Meta.Code = code
	g.resp.Meta.Message = message
	return g
}

func (g *Gin) SetStatus(status int) *Gin {
	g.status = status
	return g
}

func (g *Gin) Success(data interface{}) *Gin {
	g.resp.Meta.Code = CodeSuccess
	g.resp.Meta.Message = "success"
	g.resp.Data = data
	return g
}

// Response setting gin.JSON
func (g *Gin) Response() {
	g.c.JSON(g.status, g.resp)
	g.c.Abort()
	return
}
`, "`", "`", "`", "`", "`", "`", "`", "`")
)
