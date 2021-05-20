package service

import "fmt"

var (
	GORMServiceTemplate = fmt.Sprintf(`
{{$Sep :=.Separate}}
{{- if $Sep}}package svc_{{.Package}}{{- else}}package services{{- end}}

{{- $Package := .Package }}
import(
	{{- if $Sep}}"{{.Mod}}/app/models/{{.PackageName}}"{{- else}}"{{.Mod}}/app/models"{{- end}}
	"fmt"
	"github.com/mitchellh/mapstructure"
)

{{$StructName :=.StructName}}
{{$Primary:= ""}}
{{$PrimaryType := ""}}
{{$PrimaryHumpName := ""}}
{{- range .Fields}}{{if .IsPrimary}}
{{$Primary = .DBName}}
{{$PrimaryType = .Type}}
{{$PrimaryHumpName = .HumpName}}{{- end}}{{- end}}
// Add{{$StructName}} add one record
func Add{{$StructName}}(req *{{ if $Sep}}model_{{$Package}}{{else}}models{{end}}.Add{{$StructName}}Form)(res *{{if $Sep}}model_{{$Package}}{{else}}models{{end}}.{{$StructName}}, err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data = new({{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.{{$StructName}})
	)
	if err = mapstructure.Decode(req,data);err!=nil{
		return
	}
	// if needed todo add you business logic code

	if err = {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Add{{$StructName}}(data);err!=nil{
		return
	}

	// 
	res = data
	return
}


// Add{{$StructName}}Batch add {{$StructName}}  batch record
func Add{{$StructName}}Batch(req {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Add{{$StructName}}BatchForm)(res []*{{if $Sep}}model_{{$Package}}{{else}}models{{end}}.{{$StructName}} , err error) {
	var(
		data []* {{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.{{$StructName}}
	)
	if err = mapstructure.Decode(req,&data);err!=nil{
		return
	}
	// if needed todo add you business logic code
	if err ={{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Add{{$StructName}}Batch(data);err!=nil{
		return	
	}
	// 
	res = data
	return   
}

{{$PrimaryKey := ""}}

{{range .Fields -}}
      {{if .IsPrimary -}}
		 {{$PrimaryKey = .FieldName -}}
      {{end -}}
{{- end -}}

// Up{{$StructName}} edit {{$StructName}} one record
func Up{{$StructName}}({{$PrimaryHumpName}} interface{},req *{{ if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Up{{$StructName}}Form)(err error) {
	if err = req.Valid();err!=nil{
		return
	}
	var(
		data ={{if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.New{{$StructName}}()
	)
	// if needed todo add you business logic code code
	if err = mapstructure.Decode(req, data); err != nil {
		return
	}
	if err = {{ if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Up{{$StructName}}({{$PrimaryHumpName}},data);err!=nil{return}

	return
}

// Get{{$StructName}}List get {{$StructName}} list  data
func Get{{$StructName}}List(req *{{ if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Query{{$StructName}}Form)(res interface{}, err error) {
	var(
		data []*{{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.{{$StructName}} // default get all column,if choice some column, define struct response form by yourself
	)
	// if needed todo add you business logic code code
	if err = {{- if $Sep}}model_{{$Package}}{{- else}}models{{- end}}.Get{{$StructName}}List(req,&data);err!=nil{return}
	
	// 
	res = data
	return
}


// Get{{$StructName}} get {{$StructName}} one record
func Get{{$StructName}}(field string,value interface{})(res *{{if $Sep}}model_{{$Package}}{{else}}models{{end}}.{{$StructName}}, err error) {
	var(
		d  *{{if $Sep}}model_{{$Package}}{{else}}models{{end}}.{{$StructName}}
	)
	switch field{
	{{range .Fields -}}
      {{if .IsUnique -}}
		case "{{.DBName}}":
			if d,err = {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Get{{$StructName}}{{.FieldName}}(value);err!=nil{
				return
			} 
      {{end -}}
	{{- end -}}
	default:
		err = fmt.Errorf("field: %s not support in this way",field)
	}
	res = d
	return   
}

// Del{{$StructName}} delete {{$StructName}} one record
func Del{{$StructName}}(field string,value interface{})( err error) {
	switch field{
	{{range .Fields -}}
      {{if .IsUnique -}}
		case "{{.DBName}}":
			if err = {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Del{{$StructName}}{{.FieldName}}(value);err!=nil{
				return
			} 
      {{end -}}
	{{- end -}}
	default:
		err = fmt.Errorf("field: %s not support in this way",field)
	}
	return
}

// Del{{$StructName}}Batch delete {{$StructName}} batch record
func Del{{$StructName}}Batch(ids []string)( err error) {
	// if needed todo add you business logic code
	return   {{if $Sep}}model_{{$Package}}{{else}}models{{end}}.Del{{$StructName}}Batch(ids)
}

`, "%s", "%s")
)
