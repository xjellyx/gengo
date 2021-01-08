package controller

import (
	"fmt"
)

var GinTemplate = fmt.Sprintf(`package api_{{.Package}}
{{- $Package := .Package }}
import(
	"github.com/gin-gonic/gin"
	"{{.Mod}}/app/service/{{$Package}}"
	"{{.Mod}}/app/model/{{$Package}}"
	"{{.Mod}}/app/controller/response"
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
		req = new(srv_{{$Package}}.Edit{{$StructName}}ReqForm)
		err error
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(nil).Response()
		}
	}()
	if err = c.ShouldBind(&req);err!=nil{
		return
	}
	if err = srv_{{$Package}}.Edit{{.StructName}}One(req);err!=nil{
		return
	}
}

// GetOne get {{$StructName}} one record
// @tags {{$StructName}}
// @Summary get {{$StructName}} one record
// @Description get {{$StructName}} one record
// @Accept json
// @Produce json
// @Param {} query srv_{{.LowerName}}.Operating{{$StructName}}OneReqForm true "{{$StructName}} form, just pass a parameter"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Package}}/get  [get]
func (ct *Ctrl{{$StructName}}) GetOne(c *gin.Context) {
	var(
		data interface{}
		req =new(srv_{{.LowerName}}.Operating{{$StructName}}OneReqForm)
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
	if err =  c.ShouldBindQuery(req);err!=nil{
		return
	}
	if data,err = srv_{{$Package}}.Get{{$StructName}}One(req);err!=nil{
		return
	}
}

// GetList get {{$StructName}} list record
// @tags {{$StructName}}
// @Summary get {{$StructName}} list record
// @Description get {{$StructName}} list record
// @Accept json
// @Produce json
// @Param {} query model_{{$Package}}.Query{{$StructName}}Form true "获取{{$StructName}}列表form"
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
// @Param {} body srv_{{.LowerName}}.Operating{{$StructName}}OneReqForm true "{{$StructName}} form, just pass a parameter"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router  /api/v1/{{$Package}}/delete [delete]
func (ct *Ctrl{{$StructName}}) DeleteOne(c *gin.Context) {
	var(
		data interface{}
		err error
		req = new(srv_{{.LowerName}}.Operating{{$StructName}}OneReqForm)
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	if err = c.ShouldBind(req);err!=nil{return}
	if err = srv_{{$Package}}.Delete{{$StructName}}One(req);err!=nil{return}
}

// DeleteList delete {{$StructName}} list record
// @tags {{$StructName}}
// @Summary delete {{$StructName}} list record
// @Description delete {{$StructName}} list record
// @Accept json
// @Produce json
// @Param ids body []string true "{{$StructName}} ID list"
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
	if err = c.ShouldBind(&ids);err!=nil{return}
	if err = srv_{{$Package}}.Delete{{$StructName}}Batch(ids);err!=nil{return}
}

`)

var (
	ResponseTemplate = fmt.Sprintf(`
package response

import (
	"github.com/gin-gonic/gin"
	"sync"
)

const (
	CodeSuccess = 0
	CodeFail    = -1
)

type Gin struct {
	c      *gin.Context
	resp   *Response
	status int
}

type Response struct {
	Meta    Meta        %sjson:"meta"%s
	Code    int         %sjson:"code"%s
	Message string      %sjson:"message"%s
	Data    interface{} %sjson:"data"%s
}

type Meta map[string]interface{}

var (
	l = &sync.RWMutex{}
)

func (m Meta) Set(key string, val interface{}) {
	l.Lock()
	m[key] = val
	l.Unlock()
}

// NewGinResponse
func NewGinResponse(c *gin.Context) *Gin {
	return &Gin{
		c,
		&Response{
			Meta: Meta{},
		},
		200,
	}
}

func (g *Gin) Fail(code int, message string) *Gin {
	g.resp.Code = code
	g.resp.Message = message
	return g
}

func (g *Gin) SetStatus(status int) *Gin {
	g.status = status
	return g
}

func (g *Gin) SetMeta(key string, val interface{}) *Gin {
	g.resp.Meta.Set(key, val)
	return g
}

func (g *Gin) Success(data interface{}) *Gin {
	g.resp.Code = CodeSuccess
	g.resp.Message = "success"
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

var (
	InitRouterTemplate = fmt.Sprintf(`
package router
{{$Mod :=.Mod}}
import (
	"io"
	"net/http"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/olongfen/contrib/log"

	{{- range .Structs}}
	_"{{$Mod}}/app/controller/router/{{.LowerName}}"
	{{- end}}
	"{{.Mod}}/app/controller/middleware"
	"{{.Mod}}/app/setting"
	"{{.Mod}}/app/controller/common"
)

// 初始化路由
var (
	Engine   = gin.Default()
)
// init 初始化路由模块
func init() {
	if !setting.DevEnv {
		gin.SetMode(gin.ReleaseMode)
		Engine.Use(gin.Logger())
		// 创建记录日志的文件
		f, _ := rotatelogs.New(
			setting.Global.FilePath.LogDir + "/router" + setting.Global.FilePath.LogPatent+".log",
		)
		gin.DefaultWriter = io.MultiWriter(f)
	}


	// 添加中间件
	Engine.Use(middleware.CORS())
	Engine.Use(middleware.GinLogFormatter())
	Engine.Use(gin.Recovery())
	// 没有路由请求
	Engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": " 404 " + http.StatusText(http.StatusNotFound),
		})
	})
	// TODO 路由
	{
		var api = Engine.Group("api/v1")
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.Use(middleware.Common())

		// 测试连接
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ping": "pong >>>>>>> update"})
		})
		for _,v:=range ctrl_common.RouterGroupFunctions{
			v(api)
		}

	}
	log.Infoln("router init success !")
}
`)
)

var (
	MiddlewareTemplate = fmt.Sprintf(`package middleware
import (
	"fmt"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

const (
	PayPasswordHeader = "X-Pay-Password"
	SignatureHeader   = "X-Signature"
	Token             = "X-Token"
)

var (
	allowHeaders = strings.Join([]string{
		"accept",
		"origin",
		"Authorization",
		"Content-Type",
		"Content-Length",
		"Content-Length",
		"Accept-Encoding",
		"Cache-Control",
		"X-CSRF-Token",
		"X-Requested-With",
		Token,
		SignatureHeader,    // 接受签名的 Header
		PayPasswordHeader,  // 接收交易密码的 Header
		"X-Wechat-Binding", // 激活微信帐号
	}, ",")
	allowMethods = strings.Join([]string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}, ",")
)

// CORS
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

// Common
func Common() gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.Writer.Header()
		// alone dns prefect
		header.Set("X-DNS-Prefetch-Control", "on")
		// IE No Open
		header.Set("X-Download-Options", "noopen")
		// not cache
		header.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		header.Set("Expires", "max-age=0")
		// Content Security Policy
		header.Set("Content-Security-Policy", "default-src 'self'")
		// xss protect
		// it will caught some problems is old IE
		header.Set("X-XSS-Protection", "1; mode=block")
		// Referrer Policy
		header.Set("Referrer-Header", "no-referrer")
		// cros frame, allow same origin
		header.Set("X-Frame-Options", "SAMEORIGIN")
		// HSTS
		header.Set("Strict-Transport-Security", "max-age=5184000;includeSubDomains")
		// no sniff
		header.Set("X-Content-Type-Options", "nosniff")
	}
}
`) + `
func GinLogFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("address: %s, time: %s, method: %s, path: %s, errMessage: %s, proto: %s, code: %d, latency: %s, body: %v %v",
				params.ClientIP, params.TimeStamp.Format("2006-01-02 15:04:05"), params.Method, params.Path,
				params.ErrorMessage, params.Request.Proto, params.StatusCode, params.Latency,
				params.Request.Body, "\n")
		})
}`
)

var (
	StructRouterTemplate = fmt.Sprintf(`package router

import  ("{{.Mod}}/app/controller/api/{{.Package}}"
"{{.Mod}}/app/controller/common"
"github.com/gin-gonic/gin"
)

func init{{.StructName}}(r *gin.RouterGroup) {
	c := &api_{{.Package}}.Ctrl{{.StructName}}{}
	{{.Package}} := r.Group("{{.Package}}")
	{{.Package}}.GET("get", c.GetOne)
	{{.Package}}.GET("list", c.GetList)
	{{.Package}}.POST("add", c.AddOne)
	{{.Package}}.POST("addList", c.AddList)
	{{.Package}}.PUT("edit", c.Edit)
	{{.Package}}.DELETE("delete", c.DeleteOne)
	{{.Package}}.DELETE("deleteList", c.DeleteList)
}

func init() {
ctrl_common.RouterGroupFunctions = append(ctrl_common.RouterGroupFunctions,init{{.StructName}})
}


`)
)
