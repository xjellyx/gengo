package model

import "fmt"

var (
	CommonTemplate = fmt.Sprintf(`package {{.Package}}

import(
"github.com/olongfen/contrib/log"
"gorm.io/gorm"
)

type FieldData struct {
	Value interface{} %sjson:"value" form:"value"%s
	Symbol string %sjson:"symbol" form:"symbol"%s
}


var(
	ModelLog *log.Logger
	DB *gorm.DB
)
`, "`", "`", "`", "`")
)
