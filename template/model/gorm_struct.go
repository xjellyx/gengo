package model

import "fmt"

var (
	GORMTemplate = fmt.Sprintf(`
{{$Sep := .Separate}}
{{- if $Sep}}package model_{{.Package}}{{- else}}package models{{- end}}
{{$TFErr :=.TFErr}}
{{$Int :=  "int" }}
{{$Int8  :="int8" }}
{{$Int16 :="int16" }}
{{$Int32 :="int32" }}
{{$Int64 :="int64" }}
{{$Float64 :="float64" }}
{{$Float32 :="float32" }}
{{$Time :="time.Time" }}
import (
{{- if $TFErr}} "errors" {{end}}
	{{if $Sep}}"{{.Mod}}/app/models/common"{{- end}}
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	v1 "github.com/jinzhu/gorm"
)
{{$StructName :=.StructName}}
// Error
{{if $TFErr}} var(
	ErrCreate{{$StructName}} = errors.New("create {{$StructName}} failed")
	ErrDelete{{$StructName}} = errors.New("delete {{$StructName}} failed")
	ErrGet{{$StructName}} = errors.New("get {{$StructName}} failed")
	ErrUpdate{{$StructName}} = errors.New("update {{$StructName}} failed")
)
{{end}}
// {{$StructName}} table
type {{$StructName}} {{.StructDetail}}

func init(){
  {{- if $Sep}}model_common.{{- end}}Tables = append({{- if $Sep}}model_common.{{- end}}Tables,&{{$StructName}}{})
}
{{$Primary:= ""}}
{{$PrimaryType := ""}}
{{$PrimaryHumpName := ""}}
{{- range .Fields}}{{if .IsPrimary}}
{{$Primary = .DBName}}
{{$PrimaryType = .Type}}
{{$PrimaryHumpName = .HumpName}}{{- end}}{{- end}}

// New{{$StructName}} new
func New{{$StructName}}()*{{$StructName}}{
	return new({{$StructName}})
}

// {{$StructName}}TableName TableName 
func {{$StructName}}TableName()string{
	return "{{.LowerName}}s"
}

// Add{{$StructName}} add one record
func Add{{$StructName}}(data *{{$StructName}},dbs ...*gorm.DB)(err error) {
		if err = {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Create(data).Error;err!=nil{
			{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err) 
			err = ErrCreate{{$StructName}}{{end}}
			return
		}
		return
	}

// Del{{$StructName}} delete record
func Del{{$StructName}}({{$PrimaryHumpName}} interface{},dbs ...*gorm.DB)(err error) {
		if err =  {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).Delete("{{$Primary}} = ?",{{$PrimaryHumpName}}).Error;err!=nil{
			{{- if $TFErr}} err = ErrDelete{{$StructName}} {{end}}
			return
		}
		return
	}

// Up{{$StructName}} update record
func Up{{$StructName}}({{$PrimaryHumpName}},m interface{},dbs ...*gorm.DB)(err error) {
		if err = {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).Where("{{$Primary}} = ?",{{$PrimaryHumpName}}).Updates(m).Error;err!=nil{
			{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err)
			err = ErrUpdate{{$StructName}} {{end}}
			return
		}
		return
	}

// GetAll{{$StructName}} get all record
func GetAll{{$StructName}}(dbs ...*gorm.DB)(res []*{{$StructName}},err error){
		if err = {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Find(&res).Error;err!=nil{
			{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err) 
			err = ErrGet{{$StructName}} {{end}}
			return
		}
		return
	}

// Count{{$StructName}} get count
func Count{{$StructName}}(dbs ...*gorm.DB)(res int64){
		{{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).Count(&res)
		return
	}

// Del{{$StructName}}Batch delete {{$StructName}} batch
func Del{{$StructName}}Batch( {{$PrimaryHumpName}}s []string, dbs ...*gorm.DB)(err error){
		if err = {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).Delete("id in ?",{{$PrimaryHumpName}}s).Error;err!=nil{
			{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err) 
			err = ErrDelete{{$StructName}} {{end}}
			return
		}
		return 
	}
	
// Add{{$StructName}}Batch add {{$StructName}} batch
func Add{{$StructName}}Batch( datas []*{{$StructName}},dbs ...*gorm.DB)(err error){
		if err =  {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).Create(datas).Error;err!=nil{
			{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err) 
			err = ErrCreate{{$StructName}} {{end}}
			return
		}
		return
	}

// Get{{$StructName}}List get {{$StructName}} list some field value or some condition
func Get{{$StructName}}List(q *Query{{$StructName}}Form ,res interface{},dbs ...*gorm.DB)(err error){
		var(
			db = {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{})
		)
{{- range .Fields}}{{- if not .IsUnique}}		
{{- if eq .Type $Time -}}
		for k,v:=range q.{{.FieldName}}Map{
			if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
		db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Int -}}
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Int8 -}}
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Int16 -}} 
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Int32 -}} 
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Int64 -}} 
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Float32 -}} 
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else if eq .Type $Float64 -}} 
			for k,v:=range q.{{.FieldName}}Map{
				if k,err={{- if $Sep}}model_common.{{- end}}ValidFieldSymbol("{{.DBName}}",k);err!=nil{return}
			db=db.Where("{{.DBName}} "+k+ " ?",v)
		}
{{- else -}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}} = ?",q.{{.FieldName}})	
		}
{{- end}}	
{{- end}}
{{end}}

for k, v := range q.OrderMap {
	var(desc bool)	
	if v=="desc"{
		desc =true
	}else{
		desc=false
	}
		db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: v1.ToColumnName(k)}, Desc: desc})
}
if q.PageSize!=0{
	db = db.Limit(q.PageSize)
}
if q.PageNum!=0{
	q.PageNum = (q.PageNum - 1) * q.PageSize
	db = db.Offset(q.PageNum)
}
if err =db.Find(res).Error;err!=nil{
			return
		}
		return
	}

{{- range .Fields}}
{{- if .IsUnique}}
// Get{{$StructName}}{{.FieldName}} get one record by {{.FieldName}}
func Get{{$StructName}}{{.FieldName}}({{.HumpName}} interface{},dbs ...*gorm.DB)(res *{{$StructName}},err error){
			res=new({{$StructName}})
			if err = {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).First(res,"{{.DBName}} = ?",{{.HumpName}}).Error;err!=nil{
				{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err) 
				err = ErrGet{{$StructName}} {{end}}
				return
			}
			return
		}

// Del{{$StructName}}{{.FieldName}} delete record by {{.FieldName}}
func  Del{{$StructName}}{{.FieldName}}({{.HumpName}} interface{},dbs ...*gorm.DB)(err error) {
			if err= {{- if $Sep}}model_common.{{- end}}GetDB(dbs...).Model(&{{$StructName}}{}).Delete("{{.DBName}} = ?",{{.HumpName}}).Error;err!=nil{
				{{- if $TFErr}}{{- if $Sep}}model_common.{{- end}}ModelLog.Errorln(err) 
				err = ErrDelete{{$StructName}} {{end}}
				return
				}
			return
		}
		{{- end}}
	{{end}}
`)
	GORMInitDB = `
package models
{{$Mod :=.Mod}}
{{$Sep:= .Separate}}
import(
	"fmt"

	"github.com/olongfen/contrib/log"
	"github.com/sirupsen/logrus"
	{{if $Sep }}	{{- range .Structs}}
	_"{{$Mod}}/app/models/{{.PackageName}}"
	{{- end}}
	"{{$Mod}}/app/models/common"{{end}}
	"{{$Mod}}/app/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
)
func Init() {
	var (
		err            error
		dataSourceName string
		dialector      gorm.Dialector
	)
	{{if $Sep}}model_common.{{- end}}ModelLog = log.NewLogFile(log.ParamLog{Path: setting.Global.FilePath.LogDir + "/" + "models", Stdout: setting.DevEnv, P: setting.Global.FilePath.LogPatent})
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

	if {{if $Sep}}model_common.{{- end}}DB, err = gorm.Open(dialector, &gorm.Config{Logger: logger.New({{- if $Sep}}model_common.{{- end}}ModelLog, logger.Config{
		Colorful: true})}); err != nil {
		logrus.Fatal(err)
	}
	if setting.DevEnv {
		{{- if $Sep}}model_common.{{- end}}DB = {{- if $Sep}}model_common.{{- end}}DB.Debug()
		err = {{- if $Sep}}model_common.{{- end}}DB.AutoMigrate({{- if $Sep}}model_common.{{- end}}Tables...)
		if err != nil {
			panic(err)
		}
	}

	log.Infoln("database init success !")
}
`
	GORMForm = fmt.Sprintf(`
{{$Sep := .Separate}}
{{- if $Sep}}package model_{{.Package}}{{- else}}package models{{- end}}
{{$StructName := .StructName}}
{{$Int :=  "int" }}
{{$Int8  :="int8" }}
{{$Int16 :="int16" }}
{{$Int32 :="int32" }}
{{$Int64 :="int64" }}
{{$Float64 :="float64" }}
{{$Float32 :="float32" }}
{{$Time :="time.Time" }}
// Query{{$StructName}}Form query {{$StructName}}  form ;  if some field is required, add binding:"required" to tag by self
type Query{{$StructName}}Form struct{
{{- range .Fields}}{{- if not .IsUnique}}		
{{- if eq .Type $Time -}}
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Int -}}
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Int8 -}}
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Int16 -}} 
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Int32 -}} 
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Int64 -}} 
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Float32 -}} 
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else if eq .Type $Float64 -}} 
		{{.FieldName}}Map map[string]string %sjson:"{{.HumpName}}Map" form:"{{.HumpName}}Map"%s  // example: {{.FieldName}}Map[>]=some value&{{.FieldName}}Map[<]=some value; key must be ">,>=,<,<=,="
{{- else -}}
		{{.FieldName}} *{{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // cond {{.FieldName}}
{{- end}}	
{{- end}}
{{end}}OrderMap map[string]string %sjson:"orderMap" form:"orderMap"%s  // example: orderMap[column]=desc
		PageNum int %sjson:"pageNum" form:"pageNum"%s // get all without uploading
		PageSize int %sjson:"pageSize" form:"pageSize"%s // get all without uploading
		}


// Add{{$StructName}}Form add {{$StructName}} form
type Add{{$StructName}}Form struct {
	{{range .Fields -}}
	  {{if not .IsBaseModel -}} 
		{{if .IsUnique -}}
			{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}" binding:"required"%s // {{.HumpName}}
		{{else -}}
			{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // {{.HumpName}}
        {{end -}}
	  {{end -}}
	{{- end -}}
}

// Valid add {{$StructName}}  form verify
func (a *Add{{$StructName}}Form) Valid() (err error) {
	return
}

type Add{{$StructName}}BatchForm []*Add{{$StructName}}Form

// Up{{$StructName}}Form  edit {{$StructName}} form 
type Up{{$StructName}}Form struct {
	{{range .Fields -}}
	  {{if not .IsBaseModel -}}
		{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // {{.HumpName}}
	  {{end -}}
	{{- end -}}
}

// Valid  edit {{$StructName}} form verify
func (a *Up{{$StructName}}Form) Valid() (err error) {
	return
}

`, "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`",
		"`", "`", "`", "`", "`", "`", "`")
)
