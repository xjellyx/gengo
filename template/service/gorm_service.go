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
		{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // if required, add binding:"required" to tag by self
	  {{end -}}
	{{- end -}}
}

func (a *Edit{{$StructName}}ReqForm) Valid() (err error) {
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

	ret = data
	return
}

// Edit{{$StructName}} edit
func Edit{{$StructName}}(id int,req *Edit{{$StructName}}ReqForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		body []byte
		data =new(model_{{$Package}}.{{$StructName}})
		m = make(map[string]interface{},0)
	)
	data.ID = uint(req.ID)
	// todo add you business logic
	{{- range .Fields -}}
	  {{if not .IsBaseModel -}}
		m["{{.DBName}}"] = req.{{.FieldName}}
	  {{end -}}
	{{- end -}}
	if err = data.Updates(model_common.DB,m);err!=nil{return}
	ret = data
	return
}

`, "`", "`", "`", "`", "`", "`")
)
