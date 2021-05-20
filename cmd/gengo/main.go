package main

import (
	"github.com/olongfen/gengo/generate"
	"github.com/olongfen/gengo/parse"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const (
	transformErrorFlag = "transformError"
	inputDirFlag       = "input"
	outputDirFlag      = "output"
	modFlag            = "mod"
	webFlag            = "web"
	ormFlag            = "orm"
	separateFlag       = "separate"
	genPkgFlag         = "genPkg"
	removeSourceFlag   = "removeSource"
)

var (
	initFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     outputDirFlag,
			Aliases:  []string{"o"},
			Usage:    "The name of schema output to generate output",
			Required: false,
			Value:    ".",
		},
		&cli.StringFlag{
			Name:     inputDirFlag,
			Aliases:  []string{"i"},
			Usage:    "The name of the input go file path",
			Required: true,
		},
		&cli.StringFlag{
			Name:     modFlag,
			Aliases:  []string{"m"},
			Usage:    "The name of project go module",
			Required: true,
		},
		&cli.StringFlag{
			Name:     genPkgFlag,
			Aliases:  []string{"g"},
			Usage:    "The name of define model struct package name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     webFlag,
			Aliases:  []string{"w"},
			Usage:    "The name of project web frame",
			Required: false,
			Value:    "gin",
		},
		&cli.StringFlag{
			Name:     ormFlag,
			Aliases:  []string{"r"},
			Usage:    "The name of project orm frame",
			Required: false,
			Value:    "gorm",
		}, &cli.BoolFlag{
			Name:     transformErrorFlag,
			Aliases:  []string{"t"},
			Usage:    "The name of transform db err",
			Required: false,
			Value:    true,
		},
		&cli.BoolFlag{
			Name:     separateFlag,
			Aliases:  []string{"s"},
			Usage:    "The name of separate package",
			Required: false,
			Value:    false,
		},
		//&cli.BoolFlag{
		//	Name:     removeSourceFlag,
		//	Aliases:  []string{"rm"},
		//	Usage:    "The name of remove source",
		//	Required: false,
		//	Value:    false,
		//},
	}
)

func initAction(c *cli.Context) error {
	var (
		err error
		gen *generate.Generator
		cfg = parse.Config{
			Mod:      c.String(modFlag),
			TFErr:    c.Bool(transformErrorFlag),
			WEB:      c.String(webFlag),
			ORM:      c.String(ormFlag),
			Separate: c.Bool(separateFlag),
			GenPkg:   c.String(genPkgFlag),
			// RemoveSource: c.Bool(removeSourceFlag),
		}
	)
	//if len(c.String(outputDirFlag)) == 0 {
	//	d, _ := os.Getwd()
	//	fmt.Println(d, cfg.GenPkg)
	//	index := strings.Index(d, cfg.GenPkg)
	//	output = d[:index]
	//}

	if gen, err = generate.NewGenerator(c.String(outputDirFlag), parse.NewParser(c.String(inputDirFlag)), cfg); err != nil {
		return err
	}
	gen.Generate().Format().Flush().GenDocs()
	return nil
}

func main() {
	app := cli.NewApp()
	app.Action = initAction
	app.Flags = initFlags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
