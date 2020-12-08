package model

import "fmt"

var (
	GORMTemplate = fmt.Sprintf(`package model_{{.Package}}
{{$TFErr :=.TFErr}}
import (
{{- if $TFErr}} "errors" {{end}}
	{{- range $val := .Imports}}
		{{- if $val}}
			"{{$val}}"
		{{end}}
	{{end}}
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
// New{{.StructName}} new
func New{{.StructName}}()*{{.StructName}}{
	return new({{.StructName}})
}
	// Add add one record
	func (t *{{.StructName}}) Add(db *gorm.DB)(err error) {
		if err = db.Create(t).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrCreate{{.StructName}}{{end}}
			return
		}
		return
	}

	// Delete delete record
	func (t *{{.StructName}}) Delete(db *gorm.DB)(err error) {
		if err =  db.Delete(t).Error;err!=nil{
		
			{{- if $TFErr}} err = ErrDelete{{.StructName}} {{end}}
			return
		}
		return
	}

	// Updates update record
	func (t *{{.StructName}}) Updates(db *gorm.DB, m map[string]interface{})(err error) {
		if err = db.Model(t).Where("id = ?",t.ID).Updates(m).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err)
			err = ErrUpdate{{.StructName}} {{end}}
			return
		}
		return
	}

	// Get{{.StructName}}All get all record
	func Get{{.StructName}}All(db *gorm.DB)(ret []*{{.StructName}},err error){
		if err = db.Find(&ret).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrGet{{.StructName}} {{end}}
			return
		}
		return
	}

	// Get{{.StructName}}Count get count
	func Get{{.StructName}}Count(db *gorm.DB)(ret int64){
		db.Model(&{{.StructName}}{}).Count(&ret)
		return
	}

	// Delete{{.StructName}}Batch delete {{.StructName}} batch
	func Delete{{.StructName}}Batch(db *gorm.DB, ids []string)(err error){
		if err = db.Model(&{{.StructName}}{}).Delete("id in ?",ids).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrDelete{{.StructName}} {{end}}
			return
		}
		return 
	}
	
	// Add{{.StructName}}Batch add {{.StructName}} batch
	func Add{{.StructName}}Batch(db *gorm.DB, datas []*{{.StructName}})(err error){
		if err =  db.Model(&{{.StructName}}{}).Create(datas).Error;err!=nil{
			{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
			err = ErrCreate{{.StructName}} {{end}}
			return
		}
		return
	}

	{{$StructName := .StructName}}
	//  Query{{$StructName}}Form query form
	type Query{{$StructName}}Form struct{
	{{- range .Fields}}{{- if not .IsUnique}}		{{.FieldName}} *model_common.FieldData %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s  // if required, add binding:"required" to tag by self{{- end}}
{{end}}
		Order []string %sjson:"order" form:"order"%s
		PageNum int %sjson:"pageNum" form:"pageNum" binding:"required"%s
		PageSize int %sjson:"pageSize" form:"pageSize" binding:"required" %s
		}
	
	// Get{{$StructName}}List get {{$StructName}} list some field value or some condition
	func Get{{$StructName}}List(db *gorm.DB, q *Query{{$StructName}}Form)(ret []*{{$StructName}},err error){
		// order
		if len(q.Order)>0{
			for _,v:=range q.Order {
				db = db.Order(v)
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
	{{- range .Fields}}
		{{- if not .IsUnique}}
		if q.{{.FieldName}}!=nil{
			db = db.Where("{{.DBName}}" +q.{{.FieldName}}.Symbol +"?",q.{{.FieldName}}.Value)
		}
		{{- end}}
	{{- end}}
		if err = db.Find(&ret).Error;err!=nil{
			return
		}
		return
	}

	{{- range .Fields}}
		{{- if .IsUnique}}
			// QueryBy{{.FieldName}} query cond by {{.FieldName}}
		func (t *{{$StructName}}) SetQueryBy{{.FieldName}}({{.DBName}} {{.Type}})*{{$StructName}} {
			t.{{.FieldName}} = {{.DBName}}
			return  t
		}

		// GetBy{{.FieldName}} get one record by {{.FieldName}}
		func (t *{{$StructName}})GetBy{{.FieldName}}(db *gorm.DB)(err error){
			if err = db.First(t,"{{.DBName}} = ?",t.{{.FieldName}}).Error;err!=nil{
				{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
				err = ErrGet{{$StructName}} {{end}}
				return
			}
			return
		}

		// DeleteBy{{.FieldName}} delete record by {{.FieldName}}
		func (t *{{$StructName}}) DeleteBy{{.FieldName}}(db *gorm.DB)(err error) {
			if err= db.Delete(t,"{{.DBName}} = ?",t.{{.FieldName}}).Error;err!=nil{
				{{- if $TFErr}}model_common.ModelLog.Errorln(err) 
				err = ErrDelete{{$StructName}} {{end}}
				return
				}
			return
		}
		{{- end}}
	{{end}}
`, "`", "`", "`", "`", "`", "`", "`", "`")
	GORMInitDB = `
package model
{{$Mod :=.Mod}}
import(
	"fmt"

	"github.com/olongfen/contrib/log"
	"github.com/sirupsen/logrus"
	{{- range .Structs}}
	"{{$Mod}}/app/model/{{.LowerName}}"
	{{- end}}
	"{{$Mod}}/app/model/common"
	"{{$Mod}}/app/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres"	
)
	func init(){
	var (
		err error
		tables []interface{}
	)
	model_common.ModelLog = log.NewLogFile(log.ParamLog{Path: setting.Global.FilePath.LogDir + "/" + "models", Stdout: !setting.DevEnv, P: setting.Global.FilePath.LogPatent})
	dataSourceName := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", setting.Global.DB.Driver, setting.Global.DB.Username,
		setting.Global.DB.Password, setting.Global.DB.Host, setting.Global.DB.Port, setting.Global.DB.DatabaseName)
	if model_common.DB, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{Logger: logger.New(model_common.ModelLog, logger.Config{
		Colorful: true})}); err != nil {
		logrus.Fatal(err)
	}
	if setting.DevEnv {
		model_common.DB = model_common.DB.Debug()
	}

	{{- range  .Structs}}
		tables = append(tables,&model_{{.LowerName}}.{{.StructName}}{})
	{{end}}
	err = model_common.DB.AutoMigrate(tables ...)
	if err != nil {
		panic(err)
	}

	log.Infoln("database init success !")
}
`
)
