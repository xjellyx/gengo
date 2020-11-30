package service

import "fmt"

var (
	GORMServiceTemplate = fmt.Sprintf(`package srv_{{.Package}}
{{- $Package := .Package }}
import(
	"{{.Mod}}/model/{{$Package}}"
	"{{.Mod}}/model/common"
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
		{{.FieldName}} *{{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}" binding:"required"%s
	  {{end -}}
	{{- end -}}
}

func (a *Edit{{$StructName}}ReqForm) Valid() (err error) {
	return
}

func (a *Edit{{$StructName}}ReqForm)ToMAP()(ret map[string]interface{}){
	ret= make(map[string]interface{},0)
	{{range .Fields}}{{if not .IsBaseModel}} ret["{{.DBName}}"] = a.{{.FieldName}};{{end}}{{end}}
	return 
}

// Add{{$StructName}} add
func Add{{$StructName}}(req *Add{{$StructName}}ReqForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data = new(model_{{$Package}}.{{$StructName}})
	)
	if err = mapstructure.Decode(req,data);err!=nil{
		return
	}
	if err = data.Add(model_common.DB);err!=nil{
		return
	}

	// 
	ret = data
	return
}

// Edit{{$StructName}} edit
func Edit{{$StructName}}(id int,req *Edit{{$StructName}}ReqForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data =model_{{$Package}}.New{{$StructName}}()
	)
	// todo add you business logic
	
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
	// todo add you business logic
	
	if datas,err = model_{{$Package}}.Get{{$StructName}}List(model_common.DB,req);err!=nil{return}
	
	// 
	ret = datas
	return
}

`, "`", "`", "`", "`", "`", "`")
)
