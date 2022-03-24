package entity

import (
	"bytes"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	Template bytes.Buffer
)

func init() {
	var (
		b   []byte
		err error
	)
	_, fullFilename, _, _ := runtime.Caller(0)
	dir, _ := path.Split(fullFilename)
	if b, err = os.ReadFile(dir + "entity.tmpl"); err != nil {
		log.Fatalln(err)
	}
	Template.Write(b)

}
