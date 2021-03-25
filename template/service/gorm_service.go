package service

import "fmt"

var (
	GORMServiceTemplate = fmt.Sprintf(`package svc_{{.Package}}
{{- $Package := .Package }}
import(
	"{{.Mod}}/app/model/{{$Package}}"
	"github.com/mitchellh/mapstructure"
)

{{$StructName :=.StructName}}

// Add add one record
func Add(req *model_{{$Package}}.AddForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
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


// AddBatch add {{$StructName}}  batch record
func AddBatch(req model_{{$Package}}.AddBatchForm)(ret []* model_{{$Package}}.{{$StructName}} , err error) {
	var(
		datas []* model_{{$Package}}.{{$StructName}}
	)
	if err = mapstructure.Decode(req,&datas);err!=nil{
		return
	}
	// if needed todo add you business logic code
	if err =model_{{$Package}}.AddBatch(datas);err!=nil{
		return	
	}
	// 
	ret = datas
	return   
}

{{$PrimaryKey := ""}}

{{range .Fields -}}
      {{if .IsPrimary -}}
		 {{$PrimaryKey = .FieldName -}}
      {{end -}}
{{- end -}}

// EditOne edit {{$StructName}} one record
func EditOne(req *model_{{$Package}}.EditForm)(err error) {
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
	if err = data.SetQuery{{$PrimaryKey}}(req.{{$PrimaryKey}}).Update();err!=nil{return}

	return
}

// GetList get list {{$StructName}} data
func GetList(req *model_{{$Package}}.QueryForm)(ret []*model_{{$Package}}.{{$StructName}}, err error) {
	var(
		datas []*model_{{$Package}}.{{$StructName}}
	)
	if err = req.Valid();err!=nil{
		return
	}
	// if needed todo add you business logic code code
	if datas,err = model_{{$Package}}.GetList(req);err!=nil{return}
	
	// 
	ret = datas
	return
}


// Get get {{$StructName}} one record
func Get(req *model_{{$Package}}.OpOneForm)(ret *model_{{$Package}}.{{$StructName}}, err error) {
	var(
		d  *model_{{$Package}}.{{$StructName}}
		)
	{{range .Fields -}}
      {{if .IsUnique -}}
		if req.{{.FieldName}}!=nil{
			d = model_{{$Package}}.New{{$StructName}}()
			if err = d.SetQuery{{.FieldName}}(*req.{{.FieldName}}).GetBy{{.FieldName}}();err!=nil{
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

// Delete delete {{$StructName}} one record
func Delete(req *model_{{$Package}}.OpOneForm)( err error) {
	var(
		d  *model_{{$Package}}.{{$StructName}}
		)
	{{range .Fields -}}
      {{if .IsUnique -}}
		if req.{{.FieldName}}!=nil{
			d = model_{{$Package}}.New{{$StructName}}()
			return d.SetQuery{{.FieldName}}(*req.{{.FieldName}}).DeleteBy{{.FieldName}}()
		} 
      {{end -}}
	{{- end -}}
	return
}

// DeleteBatch delete {{$StructName}} batch record
func DeleteBatch(ids []string)( err error) {
	// if needed todo add you business logic code
	return   model_{{$Package}}.DeleteBatch(ids)
}

`)
)
