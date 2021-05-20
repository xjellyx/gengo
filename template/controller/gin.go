package controller

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	TemplateCommon       bytes.Buffer
	TemplateGin          bytes.Buffer
	TemplateResponse     bytes.Buffer
	TemplateInitRouter   bytes.Buffer
	TemplateMiddleware   bytes.Buffer
	TemplateStructRouter bytes.Buffer
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

	if b, err = os.ReadFile(dir + "gin_gorm.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateGin.Write(b)

	if b, err = os.ReadFile(dir + "gin_init_router.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateInitRouter.Write(b)

	if b, err = os.ReadFile(dir + "gin_middleware.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateMiddleware.Write(b)

	if b, err = os.ReadFile(dir + "gin_response.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateResponse.Write(b)

	if b, err = os.ReadFile(dir + "gin_struct_router.tmpl"); err != nil {
		log.Fatalln(err)
	}
	TemplateStructRouter.Write(b)
}
