package service

import "fmt"

var (
	CommonTemplate = fmt.Sprintf(`package {{.Package}}
{{$Mod :=.Mod}}
import(
"github.com/olongfen/contrib/log"
"{{$Mod}}/setting"
"gorm.io/gorm"
)

type FieldData struct {
	Value interface{} %sjson:"value" form:"value"%s
	Symbol string %sjson:"symbol" form:"symbol"%s
}


var(
	ServiceLog = log.NewLogFile(log.ParamLog{Path: setting.Global.FilePath.LogDir + "/" + "service", Stdout: setting.DevEnv, P: setting.Global.FilePath.LogPatent})
	DB *gorm.DB
)
`, "`", "`", "`", "`")
)
