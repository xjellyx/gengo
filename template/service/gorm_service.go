package service

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	TemplateCommon bytes.Buffer
	TemplateGorm   bytes.Buffer
)

func init() {
	var (
		b   []byte
		err error
	)
	_, fullFilename, _, _ := runtime.Caller(0)
	dir, _ := path.Split(fullFilename)
	if b, err = os.ReadFile(dir + "common.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateCommon.Write(b)

	if b, err = os.ReadFile(dir + "gorm.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateGorm.Write(b)

}
