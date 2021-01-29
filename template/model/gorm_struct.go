package model

import "fmt"

var (
	GORMTemplate = fmt.Sprintf(`package model_{{.Package}}
{{$TFErr :=.TFErr}}
import (
{{- if $TFErr}} "errors" {{end}}
	"{{.Mod}}/app/model/common"
	"gorm.io/gorm"
)

// Error
{{if $TFErr}} var(
	ErrCreate{{.StructName}} = errors.New("create {{.StructName}} failed")
	ErrDelete{{.StructName}} = errors.New("delete {{.StructName}} failed")
	ErrGet{{.StructName}} = errors.New("get {{.StructName}} failed")
	ErrUpdate{{.StructName}} = errors.New("update {{.StructName}} failed")
)
{{end}}
// {{.StructName}}
type {{.StructName}} {{.StructDetail}}

func init(){
	model_common.Tables = append(model_common.Tables,&{{.StructName}}{})
}

// New{{.StructName}} new
func New{{.StructName}}()*{{.StructName}}{
	return new({{.StructName}})
}

// TableName 
func TableName()string{
	return "{{.LowerName}}s"
}
	// Add add one record
	func (obj *{{.StructName}}) Add(dbs ...*gorm.DB)(err error) {
		if err = model_common.GetDB(dbs...).Create(obj).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrCreate{{.StructName}}{{end}}
			return
		}
		return
	}

	// Delete delete record
	func (obj *{{.StructName}}) Delete(dbs ...*gorm.DB)(err error) {
		if err =  model_common.GetDB(dbs...).Delete(obj).Error;err!=nil{
		
			{{- if $TFErr}} err = ErrDelete{{.StructName}} {{end}}
			return
		}
		return
	}

	// Update update record
	func (obj *{{.StructName}}) Update(dbs ...*gorm.DB)(err error) {
		if err = model_common.GetDB(dbs...).Updates(obj).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err)
			err = ErrUpdate{{.StructName}} {{end}}
			return
		}
		return
	}

	// Get{{.StructName}}All get all record
	func Get{{.StructName}}All(dbs ...*gorm.DB)(ret []*{{.StructName}},err error){
		if err = model_common.GetDB(dbs...).Find(&ret).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrGet{{.StructName}} {{end}}
			return
		}
		return
	}

	// Get{{.StructName}}Count get count
	func Get{{.StructName}}Count(dbs ...*gorm.DB)(ret int64){
		model_common.GetDB(dbs...).Model(&{{.StructName}}{}).Count(&ret)
		return
	}

	// Delete{{.StructName}}Batch delete {{.StructName}} batch
	func Delete{{.StructName}}Batch( ids []string, dbs ...*gorm.DB)(err error){
		if err = model_common.GetDB(dbs...).Model(&{{.StructName}}{}).Delete("id in ?",ids).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrDelete{{.StructName}} {{end}}
			return
		}
		return 
	}
	
	// Add{{.StructName}}Batch add {{.StructName}} batch
	func Add{{.StructName}}Batch( datas []*{{.StructName}},dbs ...*gorm.DB)(err error){
		if err =  model_common.GetDB(dbs...).Model(&{{.StructName}}{}).Create(datas).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrCreate{{.StructName}} {{end}}
			return
		}
		return
	}

	{{$StructName := .StructName}}
	{{$Int :=  "int" }}
	{{$Int8  :="int8" }}
	{{$Int16 :="int16" }}
	{{$Int32 :="int32" }}
	{{$Int64 :="int64" }}
	{{$Float64 :="float64" }}
	{{$Float32 :="float32" }}
	{{$Time :="time.Time" }}
	//  Query{{$StructName}}Form query form ;  if some field is required, add binding:"required" to tag by self
	type Query{{$StructName}}Form struct{
{{- range .Fields}}{{- if not .IsUnique}}		
{{- if eq .Type $Time -}}
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Int -}}
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Int8 -}}
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Int16 -}}cond 
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Int32 -}}cond 
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Int64 -}}cond 
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Float32 -}}cond 
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else if eq .Type $Float64 -}}cond 
		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- else -}}
		{{.FieldName}} *{{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- end}}	
{{- end}}
{{end}}
		Order []string %sjson:"order" form:"order"%s
		PageNum int %sjson:"pageNum" form:"pageNum"%s // get all without uploading
		PageSize int %sjson:"pageSize" form:"pageSize"%s // get all without uploading
		}
	
	// Get{{$StructName}}List get {{$StructName}} list some field value or some condition
	func Get{{$StructName}}List( q *Query{{$StructName}}Form,dbs ...*gorm.DB)(ret []*{{$StructName}},err error){
		var(
			db = model_common.GetDB(dbs...)
		)
		// order
		if len(q.Order)>0{
			for _,v:=range q.Order {
				db =db.Order(v)
			}
		}
		// pageSize
		if q.PageSize!=0{
			db = db.Limit(q.PageSize)
		}
		// pageNum
		if q.PageNum!=0{
			q.PageNum = (q.PageNum - 1) * q.PageSize
			db = db.Offset(q.PageNum)
		}
	{{range .Fields}}
		{{- if not .IsUnique}}
		{{- if eq .Type $Time}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Int}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Int8}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Int16}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Int32}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Int64}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Float32}}
	if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else if eq .Type $Float64}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
{{- else -}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}} = ?",*q.{{.FieldName}})
		}
{{- end}}
{{end}}
{{- end}}
		if err =db.Find(&ret).Error;err!=nil{
			return
		}
		return
	}

	{{- range .Fields}}
		{{- if .IsUnique}}
			// QueryBy{{.FieldName}} query cond by {{.FieldName}}
		func (obj *{{$StructName}}) SetQueryBy{{.FieldName}}({{.HumpName}} {{.Type}})*{{$StructName}} {
			obj. {{.FieldName}} = {{.HumpName}}
			return  obj
		}

		// GetBy{{.FieldName}} get one record by {{.FieldName}}
		func (obj *{{$StructName}})GetBy{{.FieldName}}(dbs ...*gorm.DB)(err error){
			if err = model_common.GetDB(dbs...).First(obj,{{ if not .IsBaseModel }}"{{.DBName}} = ?",obj. {{.FieldName}} {{end}}).Error;err!=nil{
				{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
				err = ErrGet{{$StructName}} {{end}}
				return
			}
			return
		}

		// DeleteBy{{.FieldName}} delete record by {{.FieldName}}
		func (obj *{{$StructName}}) DeleteBy{{.FieldName}}(dbs ...*gorm.DB)(err error) {
			if err= model_common.GetDB(dbs...).Delete(obj,{{ if not .IsBaseModel }}"{{.DBName}} = ?",obj. {{.FieldName}}{{end}}).Error;err!=nil{
				{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
				err = ErrDelete{{$StructName}} {{end}}
				return
				}
			return
		}
		{{- end}}
	{{end}}
`, "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`")
	GORMInitDB = `
package model
{{$Mod :=.Mod}}
import(
	"fmt"

	"github.com/olongfen/contrib/log"
	"github.com/sirupsen/logrus"
	{{- range .Structs}}
	_"{{$Mod}}/app/model/{{.LowerName}}"
	{{- end}}
	"{{$Mod}}/app/model/common"
	"{{$Mod}}/app/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
)
func init() {
	var (
		err            error
		dataSourceName string
		dialector      gorm.Dialector
	)
	model_common.ModelLog = log.NewLogFile(log.ParamLog{Path: setting.Global.FilePath.LogDir + "/" + "models", Stdout: setting.DevEnv, P: setting.Global.FilePath.LogPatent})
	switch setting.Global.DB.Driver {
	case "postgres":
		dataSourceName = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", setting.Global.DB.Username,
			setting.Global.DB.Password, setting.Global.DB.Host, setting.Global.DB.Port, setting.Global.DB.DatabaseName)
		//dataSourceName = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", setting.Global.DB.Driver, setting.Global.DB.Username,
		//	setting.Global.DB.Password, setting.Global.DB.Host, setting.Global.DB.Port, setting.Global.DB.DatabaseName)
		dialector = postgres.Open(dataSourceName)
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", setting.Global.DB.Username,
			setting.Global.DB.Password, setting.Global.DB.Host, setting.Global.DB.Port, setting.Global.DB.DatabaseName)
		dialector = mysql.Open(dataSourceName)
	case "sqlite":
		dialector = sqlite.Open(setting.Global.DB.Source)
	case "sqlserver":
		dataSourceName = fmt.Sprintf("%s://%s:%s@%s:%s?database=%s", setting.Global.DB.Driver, setting.Global.DB.Username,
			setting.Global.DB.Password, setting.Global.DB.Host, setting.Global.DB.Port, setting.Global.DB.DatabaseName)
		dialector = sqlserver.Open(dataSourceName)
	case "clickhouse":
		dataSourceName = fmt.Sprintf("tcp://%s:%sdatabase=%s&username=%s&password=%s&read_timeout=10&write_timeout=30", setting.Global.DB.Host, setting.Global.DB.Port,
			setting.Global.DB.DatabaseName, setting.Global.DB.Username, setting.Global.DB.Password)
		dialector = clickhouse.Open(dataSourceName)
	default:
		log.Fatalln("dose not support this sql driver >>> ", setting.Global.DB.Driver)
	}

	if model_common.DB, err = gorm.Open(dialector, &gorm.Config{Logger: logger.New(model_common.ModelLog, logger.Config{
		Colorful: true})}); err != nil {
		logrus.Fatal(err)
	}
	if setting.DevEnv {
		model_common.DB = model_common.DB.Debug()
		err = model_common.DB.AutoMigrate(model_common.Tables...)
		if err != nil {
			panic(err)
		}
	}

	log.Infoln("database init success !")
}
`
)
