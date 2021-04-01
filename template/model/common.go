package model

import "fmt"

var (
	CommonTemplate = fmt.Sprintf(`package {{.Package}}

import(
"fmt"
"github.com/olongfen/contrib/log"
"gorm.io/gorm"
)

type FieldData struct {
	Value interface{} %sjson:"value" form:"value"%s
	Symbol string %sjson:"symbol" form:"symbol"%s // symbol should send: "<", "<=", ">", ">=", "="
}

func (f *FieldData) Valid() (err error) {
	switch f.Symbol {
	case "<", "<=", ">", ">=", "=":
	default:
		err = fmt.Errorf("value: %v symbol %s invalid", f.Value, f.Symbol)
		return
	}
	return
}

var(
	ModelLog *log.Logger
	DB *gorm.DB
	Tables []interface{}
)

func GetDB(dbs ...*gorm.DB)(res *gorm.DB){
	if len(dbs)>0{
		return dbs[0]
	}
	return DB
}
`, "`", "`", "`", "`")
)
