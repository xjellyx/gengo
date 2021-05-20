package controller

import (
	"fmt"
)

var GinTemplate = fmt.Sprintf(`
{{$Sep :=.Separate}}
{{- if $Sep}}package api_{{.Package}}{{- else}}package apis{{- end}}
{{- $Package := .Package }}
{{- $Router :=.HumpName}}
import(
	"fmt"
	"github.com/gin-gonic/gin"
{{if $Sep}}"{{.Mod}}/app/services/{{.PackageName}}"{{- else}}"{{.Mod}}/app/services"{{end}}
{{if $Sep}}"{{.Mod}}/app/models/{{.PackageName}}"{{- else}}"{{.Mod}}/app/models"{{end}}
	"{{.Mod}}/app/controller/response"
)
{{$Primary:= ""}}
{{$PrimaryType := ""}}
{{$PrimaryHumpName := ""}}
{{- range .Fields}}{{if .IsPrimary}}
{{$Primary = .DBName}}
{{$PrimaryType = .Type}}
{{$PrimaryHumpName = .HumpName}}{{- end}}{{- end}}
{{$StructName := .StructName}}

// Ctl{{$StructName}} ctrl
type Ctl{{$StructName}} struct {}

// Add add {{$StructName}} one record
// @tags {{$StructName}}
// @Summary add {{$StructName}} one record
// @Description add {{$StructName}} one record
// @Accept json
// @Produce json
// @Param {} body {{ if $Sep}}model_{{$Package}}{{else}}models{{end}}.Add{{$StructName}}Form true "添加{{$StructName}}表单" 
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Router}} [post]
func (ct *Ctl{{$StructName}}) Add(c *gin.Context) {
	var(
		req = new({{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Add{{$StructName}}Form)
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
	if data, err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Add{{$StructName}}(req); err != nil {
		return
	}
	
}

// AddBatch add {{$StructName}} batch record
// @tags {{$StructName}}
// @Summary add {{$StructName}} list record
// @Description add {{$StructName}} list record
// @Accept json
// @Produce json
// @Param  {} body {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Add{{$StructName}}BatchForm true "添加{{$StructName}}表单列表" 
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Router}}/batch [post]
func (ct *Ctl{{$StructName}}) AddBatch(c *gin.Context) {
	var(
	data interface{}
	code = response.CodeFail
	req {{ if $Sep}}model_{{$Package}}{{else}}models{{- end}}.Add{{$StructName}}BatchForm
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
	
	if data,err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Add{{$StructName}}Batch(req);err!=nil{
		return
	}
}

// Update update {{$StructName}} one record
// @tags {{$StructName}}
// @Summary edit {{$StructName}} one record
// @Description edit {{$StructName}} one record
// @Accept json
// @Produce json
// @Param {{$PrimaryHumpName}} path string true "{{$PrimaryHumpName}}"
// @Param {} body {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Up{{$StructName}}Form true "update {{$StructName}} form"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Router}}/:{{$PrimaryHumpName}} [put]
func (ct *Ctl{{$StructName}}) Update(c *gin.Context) {
	var(
		req = new({{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Up{{$StructName}}Form)
		key string
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
	if key=c.Param("{{$PrimaryHumpName}}");len(key)==0{
		err=fmt.Errorf("%s must be send","{{$PrimaryHumpName}}")
		return
	}
	if err = c.ShouldBind(&req);err!=nil{
		return
	}
	if err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Up{{$StructName}}(key,req);err!=nil{
		return
	}
}

// Get get {{$StructName}} one record
// @tags {{$StructName}}
// @Summary get {{$StructName}} one record
// @Description get {{$StructName}} one record
// @Accept json
// @Produce json
// @Param {{$PrimaryHumpName}} path string true "{{$PrimaryHumpName}}"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Router}}/:{{$Primary}}  [get]
func (ct *Ctl{{$StructName}}) Get(c *gin.Context) {
	var(
		data interface{}
		key  string
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
	if key=c.Param("{{$PrimaryHumpName}}");len(key)==0{
		err=fmt.Errorf("%s must be send","{{$PrimaryHumpName}}")
		return
	}
	if data,err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Get{{$StructName}}("{{$PrimaryHumpName}}",key);err!=nil{
		return
	}
}

{{$Int :=  "int" }}
{{$Int8  :="int8" }}
{{$Int16 :="int16" }}
{{$Int32 :="int32" }}
{{$Int64 :="int64" }}
{{$Float64 :="float64" }}
{{$Float32 :="float32" }}
{{$Time :="time.Time" }}
// GetList get {{$StructName}} list record
// @tags {{$StructName}}
// @Summary get {{$StructName}} list record
// @Description get {{$StructName}} list record
// @Accept json
// @Produce json
{{- range .Fields}}{{- if not .IsUnique}}		
{{if eq .Type $Time -}}
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Int -}}
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Int8 -}}
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Int16 -}} 
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Int32 -}} 
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Int64 -}} 
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Float32 -}} 
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else if eq .Type $Float64 -}} 
// @Param	{{.HumpName}}Map query string  false  "example: {{.HumpName}}Map[>]=some value&{{.HumpName}}Map[<]=some value; key must be >,>=,<,<=,!=,=,gt,ge,lt,le,ne,eq"
{{else -}}
// @Param	{{.HumpName}} query {{.Type}}  false "{{.FieldName}}"
{{- end -}}	
{{- end -}}
{{end -}}
// @Param orderMap query string false "example: orderMap[column]=desc"
// @Param pageSize query int false "page size"
// @Param pageNum query int false "page num"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router /api/v1/{{$Router}}  [get]
func (ct *Ctl{{$StructName}}) GetList(c *gin.Context) {
	var(
		data interface{}
		req = new({{ if $Sep}}model_{{$Package}}.{{- else}}models.{{- end}}Query{{$StructName}}Form )
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
{{- range .Fields}}{{- if not .IsUnique}}		
{{if eq .Type $Time -}}
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Int -}}
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Int8 -}}
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Int16 -}} 
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Int32 -}} 
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Int64 -}} 
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Float32 -}} 
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{else if eq .Type $Float64 -}} 
	req.{{.FieldName}}Map=c.QueryMap("{{.HumpName}}Map")
{{- end -}}	
{{- end -}}
{{end -}}
	if data,err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Get{{$StructName}}List(req);err!=nil{return}
}

// Delete delete {{$StructName}} one record
// @tags {{$StructName}}
// @Summary delete {{$StructName}} one record
// @Description delete {{$StructName}} one record
// @Accept json
// @Produce json
// @Param {{$PrimaryHumpName}} path string true "{{$PrimaryHumpName}}"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router  /api/v1/{{$Router}}/:{{$PrimaryHumpName}} [delete]
func (ct *Ctl{{$StructName}}) Delete(c *gin.Context) {
	var(
		data interface{}
		err error
		key string
		code = response.CodeFail	
)
	defer func(){
		if err!=nil{
			response.NewGinResponse(c).Fail(code,err.Error()).Response()
		}else{
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	if key=c.Param("{{$PrimaryHumpName}}");len(key)==0{
		err=fmt.Errorf("%s must be send","{{$PrimaryHumpName}}")
		return
	}
	if err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Del{{$StructName}}("{{$PrimaryHumpName}}",key);err!=nil{return}
}

// DeleteBatch delete {{$StructName}} list record
// @tags {{$StructName}}
// @Summary delete {{$StructName}} list record
// @Description delete {{$StructName}} list record
// @Accept json
// @Produce json
// @Param {{$PrimaryHumpName}}s body []string true "{{$StructName}} {{$PrimaryHumpName}} list"
// @Success 200  {object} response.Response
// @Failure 500  {object} response.Response
// @router  /api/v1/{{$Router}}/batch [delete]
func (ct *Ctl{{$StructName}}) DeleteBatch(c *gin.Context) {
	var(
		data interface{}
		{{$PrimaryHumpName}}s []string
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
	if err = c.ShouldBind(&{{$PrimaryHumpName}}s);err!=nil{return}
	if err = {{- if $Sep}}svc_{{$Package}}{{- else}}services{{- end}}.Del{{$StructName}}Batch(ids);err!=nil{return}
}

`, "%s", "%s", "%s")

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

type ErrMsgData struct {
	ErrCode int
	ErrMsg  string
	Details interface{}
}

type Response struct {
	Error ErrMsgData    %sjson:"error"%s
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

func (g *Gin) NewMeta(m Meta) *Gin {
	g.resp.Meta = m
	return g
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

func (g *Gin) Fail(code int, message string, errMsg ...ErrMsgData) *Gin {
	g.resp.Code = code
	if len(errMsg) > 0 {
		g.resp.Error=errMsg[0]
	}
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

`, "`", "`", "`", "`", "`", "`", "`", "`", "`", "`")
)

var (
	InitRouterTemplate = fmt.Sprintf(`
{{$Sep :=.Separate}}
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
	{{if $Sep}}"{{$Mod}}/app/controller/common"{{else}}"{{$Mod}}/app/controller/apis"{{end}}
	"{{.Mod}}/app/controller/middleware"
	"{{.Mod}}/app/setting"
)

// 初始化路由
var (
	Engine   = gin.Default()
)
// Init 初始化路由模块
func Init() {
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
		for _,v:=range {{if $Sep}}ctl_common{{else}}apis{{end}}.RouterGroupFunctions{
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

// CORS cors
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

// Common common head
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

import  (
{{$Sep :=.Separate}}
{{- if $Sep}}
"{{.Mod}}/app/controller/apis/{{.Package}}"
"{{.Mod}}/app/controller/common"
{{- else}}"{{.Mod}}/app/controller/apis"{{- end}}

"github.com/gin-gonic/gin"
)
{{$Primary:= ""}}
{{$PrimaryType := ""}}
{{$PrimaryHumpName := ""}}
{{- range .Fields}}{{if .IsPrimary}}
{{$Primary = .DBName}}
{{$PrimaryType = .Type}}
{{$PrimaryHumpName = .HumpName}}{{- end}}{{- end}}
func init{{.StructName}}(r *gin.RouterGroup) {
	c := {{- if $Sep}}&api_{{.Package}}.{{- else}}&apis.{{- end}}Ctl{{.StructName}}{}
	{{.HumpName}} := r.Group("{{.HumpName}}")
	{{.HumpName}}.GET(":{{$PrimaryHumpName}}", c.Get)
	{{.HumpName}}.GET("", c.GetList)
	{{.HumpName}}.POST("", c.Add)
	{{.HumpName}}.POST("batch", c.AddBatch)
	{{.HumpName}}.PUT(":{{$PrimaryHumpName}}", c.Update)
	// {{.HumpName}}.DELETE(":{{$PrimaryHumpName}}", c.Delete) // default close
	// {{.HumpName}}.DELETE("batch", c.DeleteBatch) // default close
}

func init() {
 {{- if $Sep}}ctl_common{{else}}apis{{end}}.RouterGroupFunctions = append( {{- if $Sep}}ctl_common{{else}}apis{{end}}.RouterGroupFunctions,init{{.StructName}})
}


`)
)
