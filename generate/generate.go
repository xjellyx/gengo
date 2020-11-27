package generate

import (
	"bytes"
	"errors"
	"github.com/olongfen/gengo/template/service"
	"github.com/olongfen/gengo/template/setting"
	"github.com/olongfen/gengo/utils"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/olongfen/contrib/log"
	"github.com/olongfen/gengo/parse"
	"github.com/olongfen/gengo/template/model"
)

const (
	settingName = "setting"
	envName     = "settingEnv"
	devName     = "dev"
	testName    = "test"
	prodName    = "prod"
)

type Struct struct {
	LowerName  string
	StructName string
}

type InitDB struct {
	Structs []Struct
	Mod     string
}

// The Generator is the one responsible for generating the code, adding the imports, formating, and writing it to the file.
type Generator struct {
	modelBuf   map[string]*bytes.Buffer
	serviceBuf map[string]*bytes.Buffer
	settingBuf map[string]*bytes.Buffer
	outputDir  string
	config     parse.Config
	parser     *parse.Parser
	initDB     *InitDB
}

func (g *Generator) init() {
	if len(g.config.ORM) == 0 {
		g.config.ORM = "gorm"
	}
	// init setting
	g.settingBuf[settingName] = &bytes.Buffer{}
	g.settingBuf[devName] = &bytes.Buffer{}
	g.settingBuf[testName] = &bytes.Buffer{}
	g.settingBuf[prodName] = &bytes.Buffer{}
	g.settingBuf[envName] = &bytes.Buffer{}
}

// NewGenerator function creates an instance of the generator given the name of the output file as an argument.
func NewGenerator(output string, p *parse.Parser, c parse.Config) (ret *Generator, err error) {

	if err = p.ParserFile(); err != nil {
		return
	}

	if err = p.ParserStruct(); err != nil {
		return
	}
	g := &Generator{
		modelBuf:   map[string]*bytes.Buffer{},
		serviceBuf: map[string]*bytes.Buffer{},
		settingBuf: map[string]*bytes.Buffer{},
		outputDir:  output,
		parser:     p,
		initDB:     new(InitDB),
	}
	g.config = c
	g.init()
	g.initDB.Mod = c.Mod
	for _, v := range p.Structs {
		g.initDB.Structs = append(g.initDB.Structs, Struct{
			LowerName:  strings.ToLower(v.StructName),
			StructName: v.StructName,
		})
		v.Config = c
		v.Config.Package = strings.ToLower(v.StructName)
		g.modelBuf[v.StructName] = new(bytes.Buffer)
		g.serviceBuf[v.StructName] = new(bytes.Buffer)
	}
	g.modelBuf["common"] = &bytes.Buffer{}
	g.modelBuf["initDB"] = &bytes.Buffer{}

	g.serviceBuf["common"] = &bytes.Buffer{}

	return g, nil
}

func (g *Generator) checkConfig() (err error) {
	if len(g.config.Imports) == 0 {
		err = errors.New("import package dose'n set")
		return
	}
	return
}

// genModel
func (g *Generator) genModel() (err error) {
	var (
		temp *template.Template
	)
	// 生成model公共代码
	if temp, err = template.New("common").Parse(model.CommonTemplate); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	c.Package = "model_common"
	if err = temp.Execute(g.modelBuf["common"], c); err != nil {
		return
	}

	// 生成初始化数据库代码
	if temp, err = template.New("initDB").Parse(model.GORMInitDB); err != nil {
		return
	}
	if err = temp.Execute(g.modelBuf["initDB"], g.initDB); err != nil {
		return
	}
	for _, v := range g.parser.Structs {
		var (
			t *template.Template
		)
		if _, ok := g.modelBuf[v.StructName]; !ok {
			continue
		}
		switch g.config.ORM {
		case "gorm":
			if t, err = template.New(v.StructName).Parse(model.GORMTemplate); err != nil {
				return
			}
			if err = t.Execute(g.modelBuf[v.StructName], v); err != nil {
				log.Fatalln(err)
			}
		}
	}
	return
}

// genSetting
func (g *Generator) genSetting() (err error) {
	var (
		temp *template.Template
	)

	if temp, err = template.New(settingName).Parse(setting.SettingTemplate); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[settingName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(envName).Parse(setting.EnvTemplate); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[envName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(devName).Parse(setting.ConfigTemplate); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[devName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(testName).Parse(setting.ConfigTemplate); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[testName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(prodName).Parse(setting.ConfigTemplate); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[prodName], nil); err != nil {
		return
	}
	return
}

// genSetting
func (g *Generator) genService() (err error) {
	var (
		temp *template.Template
	)
	// 生成model公共代码
	if temp, err = template.New("common").Parse(service.CommonTemplate); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	c.Package = "srv_common"
	if err = temp.Execute(g.serviceBuf["common"], c); err != nil {
		return
	}

	for _, v := range g.parser.Structs {
		var (
			t *template.Template
		)
		if _, ok := g.serviceBuf[v.StructName]; !ok {
			continue
		}

		switch g.config.ORM {
		case "gorm":
			if t, err = template.New(v.StructName).Parse(service.GORMServiceTemplate); err != nil {
				return
			}
			if err = t.Execute(g.serviceBuf[v.StructName], v); err != nil {
				log.Fatalln(err)
			}
		}
	}
	return
}

// Generate executes the template and store it in an internal buffer.
func (g *Generator) Generate() *Generator {
	if err := g.checkConfig(); err != nil {
		panic(err)
	}
	if err := g.genModel(); err != nil {
		log.Fatalln(err)
	}
	if err := g.genSetting(); err != nil {
		log.Fatalln(err)
	}
	if err := g.genService(); err != nil {
		log.Fatalln(err)
	}
	return g
}

func (g *Generator) formatModel() {
	for k, _ := range g.modelBuf {
		formatedOutput, err := format.Source(g.modelBuf[k].Bytes())
		if err != nil {
			log.Warnln(string(g.modelBuf[k].Bytes()))
			panic(err)
		}
		g.modelBuf[k] = bytes.NewBuffer(formatedOutput)
	}
}

func (g *Generator) formatSetting() {
	formatedOutput, err := format.Source(g.settingBuf[settingName].Bytes())
	if err != nil {
		panic(err)
	}
	g.settingBuf[settingName] = bytes.NewBuffer(formatedOutput)
}

func (g *Generator) formatService() {
	for k, _ := range g.serviceBuf {
		formatedOutput, err := format.Source(g.serviceBuf[k].Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		g.serviceBuf[k] = bytes.NewBuffer(formatedOutput)
	}
}

// Format function formates the output of the generation.
func (g *Generator) Format() *Generator {
	g.formatSetting()
	g.formatModel()
	g.formatService()
	return g
}

// Flush function writes the output to the output file.
func (g *Generator) Flush() (err error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for k, _ := range g.modelBuf {
			s := strings.ToLower(k)
			dir := g.outputDir + "/model/" + s
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					log.Fatalln(err)
				}
			}
			filename := dir + "/gen_" + s + ".go"
			if utils.Exists(filename) {
				continue
			}

			if err = ioutil.WriteFile(filename, g.modelBuf[k].Bytes(), 0777); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for k, _ := range g.settingBuf {
			switch k {
			case settingName:
				dir := g.outputDir + "/setting"
				filename := dir + "/gen_setting.go"
				if err = os.MkdirAll(dir, 0777); err != nil {
					if !os.IsExist(err) {
						log.Fatalln(err)
					}
				}
				if utils.Exists(filename) {
					continue
				}
				if err = ioutil.WriteFile(filename, g.settingBuf[k].Bytes(), 0777); err != nil {
					log.Fatalln(err)
				}
			case envName:
				dir := g.outputDir + "/conf"
				filename := dir + "/.env"
				if err = os.MkdirAll(dir, 0777); err != nil {
					if !os.IsExist(err) {
						log.Fatalln(err)
					}
				}
				if utils.Exists(filename) {
					continue
				}
				if err = ioutil.WriteFile(filename, g.settingBuf[k].Bytes(), 0777); err != nil {
					log.Fatalln(err)
				}
			case devName, testName, prodName:
				dir := g.outputDir + "/conf/"
				filename := dir + k + "-global-config" + ".yaml"
				if err = os.MkdirAll(dir, 0777); err != nil {
					if !os.IsExist(err) {
						log.Fatalln(err)
					}
				}

				if utils.Exists(filename) {
					continue
				}
				if err = ioutil.WriteFile(filename, g.settingBuf[k].Bytes(), 0777); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for k, _ := range g.serviceBuf {
			s := strings.ToLower(k)
			dir := g.outputDir + "/service/" + s
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					log.Fatalln(err)
				}
			}
			filename := dir + "/gen_" + s + ".go"
			if utils.Exists(filename) {
				continue
			}

			if err = ioutil.WriteFile(filename, g.serviceBuf[k].Bytes(), 0777); err != nil {
				log.Fatalln(err)
			}
		}
	}()
	wg.Wait()
	return nil
}
