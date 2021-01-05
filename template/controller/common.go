package controller

import "fmt"

var (
	CommonTemplate = fmt.Sprintf(`package {{.Package}}
{{$Mod :=.Mod}}
import(
"github.com/gin-gonic/gin"
"github.com/olongfen/contrib/log"
"{{$Mod}}/app/setting"
)



var(
	ControlLog = log.NewLogFile(log.ParamLog{Path: setting.Global.FilePath.LogDir + "/" + "controller", Stdout: setting.DevEnv, P: setting.Global.FilePath.LogPatent})
	RouterGroupFunctions []func(group *gin.RouterGroup)
)
`)
)
