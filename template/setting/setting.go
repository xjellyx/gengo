package setting

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	TemplateSetting bytes.Buffer
	TemplateEnv     = `ENVIRONMENT=dev`
	TemplateConfig  bytes.Buffer
)

func init() {
	var (
		b   []byte
		err error
	)
	_, fullFilename, _, _ := runtime.Caller(0)
	dir, _ := path.Split(fullFilename)
	if b, err = os.ReadFile(dir + "setting.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateSetting.Write(b)
	if b, err = os.ReadFile(dir + "config.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateConfig.Write(b)
}
