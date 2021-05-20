package gin_main

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	TemplateGinMain bytes.Buffer
)

func init() {
	var (
		b   []byte
		err error
	)
	_, fullFilename, _, _ := runtime.Caller(0)
	dir, _ := path.Split(fullFilename)
	if b, err = os.ReadFile(dir + "gin_gorm.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateGinMain.Write(b)
}
