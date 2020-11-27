package service

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
	ModelLog = log.NewLogFile(log.ParamLog{Path:   "./log/model",Stdout: true})
	DB *gorm.DB
)
`, "`", "`", "`", "`")
)
