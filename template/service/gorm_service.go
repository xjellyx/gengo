package service

import "fmt"

var (
	GORMServiceTemplate = fmt.Sprintf(`package srv_{{.Package}}
{{- $Package := .Package }}
import(
	"{{.Mod}}/app/model/{{$Package}}"
	"github.com/mitchellh/mapstructure"
)
{{$PrimaryKey := ""}}
{{$PrimaryKeyType := ""}}
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

{{$PrimaryKey:=""}}
// Edit{{$StructName}}ReqForm
type Edit{{$StructName}}ReqForm struct {
	{{range .Fields -}}
      {{if .IsPrimary -}}
		{{$PrimaryKey = .FieldName -}} 
		{{$PrimaryKeyType = .Type -}}
		{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}" binding:"required"%s 
      {{end -}}
	  {{if not .IsBaseModel -}}
		{{.FieldName}} {{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // if required, add binding:"required" to tag by self
	  {{end -}}
	{{- end -}}
}

func (a *Edit{{$StructName}}ReqForm) Valid() (err error) {
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

	if err = data.Add();err!=nil{
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
	if err =model_{{$Package}}.Add{{$StructName}}Batch(datas);err!=nil{
		return	
	}
	// 
	ret = datas
	return   
}

// Edit{{$StructName}}One edit
func Edit{{$StructName}}One(req *Edit{{$StructName}}ReqForm)(err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data =model_{{$Package}}.New{{$StructName}}()
	)
	// if needed todo add you business logic code code
	if err = mapstructure.Decode(req, data); err != nil {
		return
	}
	if err = data.SetQueryBy{{$PrimaryKey}}(req.{{$PrimaryKey}}).Update();err!=nil{return}

	return
}

// Get{{$StructName}}Page get page {{$StructName}} data
func Get{{$StructName}}Page(req *model_{{$Package}}.Query{{$StructName}}Form)(ret []*model_{{$Package}}.{{$StructName}}, err error) {
	var(
		datas []*model_{{$Package}}.{{$StructName}}
	)
	// if needed todo add you business logic code code
	
	if datas,err = model_{{$Package}}.Get{{$StructName}}List(req);err!=nil{return}
	
	// 
	ret = datas
	return
}

// Operating{{$StructName}}OneReqForm
type Operating{{$StructName}}OneReqForm struct {
	{{range .Fields -}}
      {{if .IsUnique -}}
		{{.FieldName}} *{{.Type}} %sjson:"{{.HumpName}}" form:"{{.HumpName}}"%s // this form just pass a parameter 
      {{end -}}
	{{- end -}}
}
// Get{{$StructName}}One get {{$StructName}} 
func Get{{$StructName}}One(req *Operating{{$StructName}}OneReqForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	var(
		d  *model_{{$Package}}.{{$StructName}}
		)
	{{range .Fields -}}
      {{if .IsUnique -}}
		if req.{{.FieldName}}!=nil{
			d = model_{{$Package}}.New{{$StructName}}()
			if err = d.SetQueryBy{{.FieldName}}(*req.{{.FieldName}}).GetBy{{.FieldName}}();err!=nil{
				return
			}
			goto RETURN
		} 
      {{end -}}
	{{- end -}}

RETURN:
	ret = d
	return   
}

// Delete{{$StructName}}One delete {{$StructName}} 
func Delete{{$StructName}}One(req *Operating{{$StructName}}OneReqForm)( err error) {
	var(
		d  *model_{{$Package}}.{{$StructName}}
		)
	{{range .Fields -}}
      {{if .IsUnique -}}
		if req.{{.FieldName}}!=nil{
			d = model_{{$Package}}.New{{$StructName}}()
			return d.SetQueryBy{{.FieldName}}(*req.{{.FieldName}}).DeleteBy{{.FieldName}}()
		} 
      {{end -}}
	{{- end -}}
	return
}

// Delete{{$StructName}}Batch delete {{$StructName}} 
func Delete{{$StructName}}Batch(ids []string)( err error) {
	// if needed todo add you business logic code
	return   model_{{$Package}}.Delete{{$StructName}}Batch(ids)
}

`, "`", "`", "`", "`", "`", "`", "`", "`")
)
