package service

import "fmt"

var (
	GORMServiceTemplate = fmt.Sprintf(`package srv_{{.Package}}
{{- $Package := .Package }}
import(
	"strconv"
	"{{.Mod}}/app/model/{{$Package}}"
	"{{.Mod}}/app/model/common"
	"github.com/mitchellh/mapstructure"
)
	{{$StructName :=.StructName}}
// Add{{$StructName}}ReqForm
type Add{{$StructName}}ReqForm struct {
	{{- range .Fields -}}
	  {{if not .IsBaseModel -}} 
		{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // if required, add binding:"required" to tag by self
	  {{end -}}
	{{end -}}
}

func (a *Add{{$StructName}}ReqForm) Valid() (err error) {
	return
}

// Edit{{$StructName}}ReqForm
type Edit{{$StructName}}ReqForm struct {
	ID int64 %sjson:"id" form:"id" binding:"required"%s
	{{range .Fields -}}
	  {{if not .IsBaseModel -}}
		{{.FieldName}} *{{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // if required, add binding:"required" to tag by self
	  {{end -}}
	{{- end -}}
}

func (a *Edit{{$StructName}}ReqForm) Valid() (err error) {
	return
}

func (a *Edit{{$StructName}}ReqForm)ToMAP()(ret map[string]interface{}){
	ret= make(map[string]interface{},0)
	{{range .Fields}}{{if not .IsBaseModel}} if a.{{.FieldName}}!=nil{ ret["{{.DBName}}"] = *a.{{.FieldName}}};{{end}}{{end}}
	return 
}

// Add{{$StructName}}One add
func Add{{$StructName}}One(req *Add{{$StructName}}ReqForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data = new(model_{{$Package}}.{{$StructName}})
	)
	if err = mapstructure.Decode(req,data);err!=nil{
		return
	}
	// if needed todo add you business logic code

	if err = data.Add(model_common.DB);err!=nil{
		return
	}

	// 
	ret = data
	return
}

type {{$StructName}}BatchForm []*Add{{$StructName}}ReqForm

// Add{{$StructName}}Batch add {{$StructName}} 
func Add{{$StructName}}Batch(req {{$StructName}}BatchForm)(ret []* model_{{$Package}}.{{$StructName}} , err error) {
	var(
		datas []* model_{{$Package}}.{{$StructName}}
	)
	if err = mapstructure.Decode(req,&datas);err!=nil{
		return
	}
	// if needed todo add you business logic code
	if err =model_{{$Package}}.Add{{$StructName}}Batch(model_common.DB,datas);err!=nil{
		return	
	}
	// 
	ret = datas
	return   
}

// Edit{{$StructName}}One edit
func Edit{{$StructName}}One(req *Edit{{$StructName}}ReqForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data =model_{{$Package}}.New{{$StructName}}()
	)
	// if needed todo add you business logic code code
	
	if err = data.SetQueryByID(uint(req.ID)).Updates(model_common.DB,req.ToMAP());err!=nil{return}
	
	// 
	ret = data
	return
}

// Get{{$StructName}}Page get page {{$StructName}} data
func Get{{$StructName}}Page(req *model_{{$Package}}.Query{{$StructName}}Form)(ret []*model_{{$Package}}.{{$StructName}}, err error) {
	var(
		datas []*model_{{$Package}}.{{$StructName}}
	)
	// if needed todo add you business logic code code
	
	if datas,err = model_{{$Package}}.Get{{$StructName}}List(model_common.DB,req);err!=nil{return}
	
	// 
	ret = datas
	return
}

// Get{{$StructName}}One get {{$StructName}} 
func Get{{$StructName}}One(in string)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	var(
		id,_ = strconv.Atoi(in)
		d = model_{{$Package}}.New{{$StructName}}().SetQueryByID(uint(id))
	)
	if err = d.GetByID(model_common.DB);err!=nil{return}

	ret = d
	return   
}

// Delete{{$StructName}}One delete {{$StructName}} 
func Delete{{$StructName}}One(in string)( err error) {

	var(
	id,_ = strconv.Atoi(in)
		d = model_{{$Package}}.New{{$StructName}}().SetQueryByID(uint(id))
	)
	// if needed todo add you business logic code
	return   d.DeleteByID(model_common.DB)
}

// Delete{{$StructName}}Batch delete {{$StructName}} 
func Delete{{$StructName}}Batch(ids []string)( err error) {
	// if needed todo add you business logic code
	return   model_{{$Package}}.Delete{{$StructName}}Batch(model_common.DB,ids)
}

`, "`", "`", "`", "`", "`", "`")
)
