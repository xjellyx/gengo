package main

import (
	"flag"
	"os"
	"strings"

	"github.com/olongfen/gengo/generate"
	"github.com/olongfen/gengo/parse"
)

var (
	tfErr bool
	input,
	output,
	mod,
	web,
	orm,
	imports string
)

func parseFlags() {

	flag.StringVar(&output, "output", "", "[Required] The name of schema output to generate output for, comma separated")
	flag.StringVar(&input, "input", "", "[Required] The name of the input go file path")
	flag.StringVar(&mod, "mod", "", "[Required] The name of project go module")
	flag.StringVar(&web, "web", "", "[Option] The name of project web frame")
	flag.StringVar(&orm, "orm", "", "[Option] The name of project orm frame")
	flag.StringVar(&imports, "imports", "", "[Required] The name of the import  to import package")
	flag.BoolVar(&tfErr, "tfErr", false, "[Option] The name of transform db err")
	flag.Parse()

	if len(input) == 0 || len(mod) == 0 || len(output) == 0 {
		flag.Usage()
		os.Exit(1)
	}

}

func main() {
	parseFlags()
	var (
		err error
		gen *generate.Generator
	)

	if gen, err = generate.NewGenerator(output, parse.NewParser(input), parse.Config{
		Imports: func() (ret []string) {
			s := strings.Split(imports, ",")
			for _, v := range s {
				ret = append(ret, v)
			}
			return ret
		}(),
		Mod:   mod,
		TFErr: tfErr,
		WEB:   web,
		ORM:   orm,
	}); err != nil {
		panic(err)
	}
	gen.Generate().Format().Flush().GenDocs()

}
