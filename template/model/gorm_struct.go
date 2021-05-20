package model

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	TemplateGORM       bytes.Buffer
	TemplateGORMInitDB bytes.Buffer
	TemplateGORMForm   bytes.Buffer
	TemplateCommon     bytes.Buffer
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

	if b, err = os.ReadFile(dir + "init_db_gorm.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateGORMInitDB.Write(b)

	if b, err = os.ReadFile(dir + "gorm.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateGORM.Write(b)

	if b, err = os.ReadFile(dir + "form.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateGORMForm.Write(b)
}
