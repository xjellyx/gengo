package parse

import (
	"fmt"
	"github.com/jinzhu/inflection"
	"github.com/olongfen/contrib/log"
	"github.com/olongfen/gengo/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/jinzhu/gorm"
)

// Field struct field detail
type Field struct {
	DBName      string // database name
	FieldName   string // field name
	HumpName    string // hump name
	PackageName string // package name
	Type        string // field type
	IsBaseModel bool   // base model field
	IsUnique    bool   // is unique true
	IsPrimary   bool   // is primary kry
}

// StructData struct data
type StructData struct {
	Config
	StructDetail string // struct detail
	StructName   string // struct name
	PluralName   string
	LowerName    string
	HumpName     string
	PackageName  string
	Fields       []*Field // struct field
}

type Config struct {
	Package      string
	TFErr        bool
	Mod          string
	ORM          string
	WEB          string
	Separate     bool
	RemoveSource bool
	GenPkg       string
}

// Parser parse struct
type Parser struct {
	Filepath      string
	Structs       []*StructData
	Files         map[string]*ast.File // key: filename value: ast.File
	CacheFileByte map[string][]byte
	fs            *token.FileSet
}

// NewParser new
func NewParser(f string) *Parser {
	return &Parser{
		Files:         map[string]*ast.File{},
		fs:            token.NewFileSet(),
		CacheFileByte: make(map[string][]byte),
		Filepath:      f,
	}
}

// ParserFile 导入文件并获取go文件
func (p *Parser) ParserFile() (err error) {
	var (
		f *ast.File
	)
	// 解析文件数据
	if f, err = parser.ParseFile(p.fs, p.Filepath, nil, 0); err != nil {
		return
	}
	// 缓存ｇｏ文件数据
	p.Files[p.Filepath] = f

	return
}

// ParserStruct 解析结构体数据
func (p *Parser) ParserStruct() (err error) {
	for k, f := range p.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				return true
			}
			for _, v := range decl.Specs {
				ts, _ok := v.(*ast.TypeSpec)
				if !_ok {
					continue
				}
				var (
					data      = new(StructData)
					haveModel = false
				)
				data.StructName = ts.Name.Name
				data.LowerName = gorm.ToDBName(data.StructName)
				data.HumpName = utils.SQLColumnToHumpStyle(data.LowerName)
				data.PluralName = inflection.Plural(data.LowerName)
				data.PackageName = utils.SQLColumn2PkgStyle(data.LowerName)
				var structType *ast.StructType
				if structType, ok = ts.Type.(*ast.StructType); !ok {
					continue
				}
				// 只读取含有结构体的.go文件,读取文件缓存起来，已经读取的略过
				if _, ok = p.CacheFileByte[k]; !ok {
					d, _ := ioutil.ReadFile(k)
					p.CacheFileByte[k] = d

				}
				data.StructDetail = string(p.CacheFileByte[k][structType.Pos()-1 : structType.End()-1])
				for _, fd := range structType.Fields.List {
					var (
						fieldData = new(Field)
					)
					// 字段
					if t, ok1 := fd.Type.(*ast.Ident); ok1 {
						fieldData.Type = t.String()
						var r rune
						for _, d := range fieldData.Type {
							r = d
							break
						}
						if unicode.IsUpper(r) {
							continue
						}
						fieldData.FieldName = fd.Names[0].String()
						fieldData.DBName = gorm.ToDBName(fieldData.FieldName)
						fieldData.HumpName = utils.SQLColumnToHumpStyle(fieldData.DBName)

						if fd.Tag != nil && (strings.Contains(fd.Tag.Value, "primary") ||
							strings.Contains(fd.Tag.Value, "unique")) {
							if !strings.Contains(fd.Tag.Value, "unique:") &&
								!strings.Contains(fd.Tag.Value, "uniqueIndex:") {
								fieldData.IsUnique = true
							}
							//
							if strings.Contains(fd.Tag.Value, "primary") {
								fieldData.IsPrimary = true
								haveModel = true
							}
							if strings.Contains(fd.Tag.Value, "primary") || strings.HasSuffix(fieldData.DBName,
								"id") {
								fieldData.IsBaseModel = true
							}
						}
						data.Fields = append(data.Fields, fieldData)
						continue
					}

					// 基本model字段,自动添加,为后面搜索使用
					if _v, ok2 := fd.Type.(*ast.SelectorExpr); ok2 {
						if _v.Sel.Name == "Model" {
							haveModel = true
							idField := new(Field)
							idField.FieldName = "ID"
							idField.Type = "uint"
							idField.IsUnique = true
							idField.IsBaseModel = true
							idField.IsPrimary = true
							idField.HumpName = "id"
							idField.DBName = gorm.ToDBName("ID")

							createdAtField := new(Field)
							createdAtField.FieldName = "CreatedAt"
							createdAtField.Type = "time.Time"
							createdAtField.IsBaseModel = true
							createdAtField.HumpName = "createdAt"
							createdAtField.DBName = gorm.ToDBName("CreatedAt")

							updatedAtField := new(Field)
							updatedAtField.FieldName = "UpdatedAt"
							updatedAtField.Type = "time.Time"
							updatedAtField.HumpName = "updatedAt"
							updatedAtField.IsBaseModel = true
							updatedAtField.DBName = gorm.ToDBName("UpdatedAt")
							data.Fields = append(data.Fields, idField, createdAtField, updatedAtField)
						}
						if _v.Sel.Name == "Time" {
							fieldData.Type = "time.Time"
							fieldData.FieldName = fd.Names[0].String()
							fieldData.DBName = gorm.ToDBName(fieldData.FieldName)
							fieldData.HumpName = utils.SQLColumnToHumpStyle(fieldData.DBName)
							switch fieldData.FieldName {
							case "CreatedAt", "CreatedTime", "CreateTime", "CreateAt",
								"UpdatedAt", "UpdatedTime", "UpdateTime", "UpdateAt",
								"DeletedAt", "DeletedTime", "DeleteTime", "DeleteAt":
								fieldData.IsBaseModel = true

							}
							data.Fields = append(data.Fields, fieldData)
						}

					}
				}
				if !haveModel {
					err = fmt.Errorf("%s %s", "please set the primary key of the table ", data.StructName)
					log.Fatal(err)
				}
				p.Structs = append(p.Structs, data)
			}
			return true
		})
	}
	return
}
