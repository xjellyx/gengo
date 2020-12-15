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
)

var (
	initFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     outputDirFlag,
			Aliases:  []string{"o"},
			Usage:    "The name of schema output to generate output",
			Required: true,
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
	}
)

func initAction(c *cli.Context) error {
	var (
		err error
		gen *generate.Generator
	)
	if gen, err = generate.NewGenerator(c.String(outputDirFlag), parse.NewParser(c.String(inputDirFlag)), parse.Config{
		Mod:   c.String(modFlag),
		TFErr: c.Bool(transformErrorFlag),
		WEB:   c.String(webFlag),
		ORM:   c.String(ormFlag),
	}); err != nil {
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
