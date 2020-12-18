package generate

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/olongfen/gengo/template/controller"
	gin_main "github.com/olongfen/gengo/template/main"
	"github.com/olongfen/gengo/template/service"
	"github.com/olongfen/gengo/template/setting"
	"github.com/olongfen/gengo/utils"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
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
	initDB      = "initDB"
	common      = "common"
	response    = "response"
	initRouter  = "initRouter"
	middleware  = "middleware"
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
	controlBuf map[string]*bytes.Buffer
	routerBuf  map[string]*bytes.Buffer
	mainBuf    *bytes.Buffer
	outputDir  string
	config     parse.Config
	parser     *parse.Parser
	initDB     *InitDB
}

func (g *Generator) init() {
	if len(g.config.ORM) == 0 {
		g.config.ORM = "gorm"
	}
	if len(g.config.WEB) == 0 {
		g.config.WEB = "gin"
	}
	// init setting
	g.settingBuf[settingName] = &bytes.Buffer{}
	g.settingBuf[devName] = &bytes.Buffer{}
	g.settingBuf[testName] = &bytes.Buffer{}
	g.settingBuf[prodName] = &bytes.Buffer{}
	g.settingBuf[envName] = &bytes.Buffer{}
	//
	g.modelBuf[common] = &bytes.Buffer{}
	g.modelBuf[initDB] = &bytes.Buffer{}
	g.serviceBuf[common] = &bytes.Buffer{}
	g.controlBuf[common] = &bytes.Buffer{}
	g.controlBuf[middleware] = &bytes.Buffer{}
	g.controlBuf[response] = &bytes.Buffer{}
	g.routerBuf[initRouter] = &bytes.Buffer{}
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
		controlBuf: map[string]*bytes.Buffer{},
		routerBuf:  map[string]*bytes.Buffer{},
		mainBuf:    &bytes.Buffer{},
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
		g.controlBuf[v.StructName] = new(bytes.Buffer)
		g.routerBuf[v.StructName] = new(bytes.Buffer)
	}

	return g, nil
}

// genModel
func (g *Generator) genModel() (err error) {
	var (
		temp *template.Template
	)
	// 生成model公共代码
	if temp, err = template.New(common).Parse(model.CommonTemplate); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	c.Package = "model_common"
	if err = temp.Execute(g.modelBuf[common], c); err != nil {
		return
	}

	// 生成初始化数据库代码
	if temp, err = template.New(initDB).Parse(model.GORMInitDB); err != nil {
		return
	}
	if err = temp.Execute(g.modelBuf[initDB], g.initDB); err != nil {
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
				return
			}
		default:
			return fmt.Errorf("this type is not currently supported >>>> %s", g.config.ORM)
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

// genService
func (g *Generator) genService() (err error) {
	var (
		temp *template.Template
	)
	// 生成model公共代码
	if temp, err = template.New(common).Parse(service.CommonTemplate); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	c.Package = "srv_common"
	if err = temp.Execute(g.serviceBuf[common], c); err != nil {
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
				return
			}
		default:
			return fmt.Errorf("this type is not currently supported >>>> %s", g.config.ORM)
		}
	}
	return
}

// genControl
func (g *Generator) genControl() (err error) {
	var (
		temp *template.Template
	)
	// common
	if temp, err = template.New(common).Parse(controller.CommonTemplate); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	c.Package = "ctrl_common"
	if err = temp.Execute(g.controlBuf[common], c); err != nil {
		return
	}
	// 初始化代码
	switch g.config.WEB {
	case "gin":
		// response
		if temp, err = template.New(response).Parse(controller.ResponseTemplate); err != nil {
			return
		}
		if err = temp.Execute(g.controlBuf[response], nil); err != nil {
			return
		}

		// middleware
		if temp, err = template.New(middleware).Parse(controller.MiddlewareTemplate); err != nil {
			return
		}
		if err = temp.Execute(g.controlBuf[middleware], nil); err != nil {
			return
		}

		// init router
		if temp, err = template.New(initRouter).Parse(controller.InitRouterTemplate); err != nil {
			return
		}
		rc := struct {
			Structs []*parse.StructData
			Mod     string
		}{}
		rc.Mod = g.config.Mod
		rc.Structs = g.parser.Structs
		if err = temp.Execute(g.routerBuf[initRouter], rc); err != nil {
			return
		}
	default:
		return fmt.Errorf("this type is not currently supported >>>> %s", g.config.WEB)
	}

	//
	for _, v := range g.parser.Structs {
		var (
			t *template.Template
		)
		if _, ok := g.controlBuf[v.StructName]; !ok {
			continue
		}

		switch g.config.WEB {
		case "gin":
			if t, err = template.New(v.StructName).Parse(controller.GinTemplate); err != nil {
				return
			}
			if err = t.Execute(g.controlBuf[v.StructName], v); err != nil {
				log.Fatalln(err)
			}

			//
			if t, err = template.New(v.StructName).Parse(controller.StructRouterTemplate); err != nil {
				return
			}
			if err = t.Execute(g.routerBuf[v.StructName], v); err != nil {
				log.Fatalln(err)
			}
		default:
			return fmt.Errorf("this type is not currently supported >>>> %s", g.config.WEB)
		}
	}
	return
}

// genMain
func (g *Generator) genMain() (err error) {
	var (
		temp *template.Template
	)
	switch g.config.WEB {
	case "gin":
		// common
		if temp, err = template.New(common).Parse(gin_main.GINMainTemplate); err != nil {
			return
		}
		c := parse.Config{}
		c = g.config
		c.Package = "main"
		if err = temp.Execute(g.mainBuf, c); err != nil {
			return
		}
	default:
		return fmt.Errorf("this type is not currently supported >>>> %s", g.config.WEB)

	}
	return
}

// Generate executes the template and store it in an internal buffer.
func (g *Generator) Generate() *Generator {
	if err := g.genModel(); err != nil {
		log.Fatalln(err)
	}
	if err := g.genSetting(); err != nil {
		log.Fatalln(err)
	}
	if err := g.genService(); err != nil {
		log.Fatalln(err)
	}
	if err := g.genControl(); err != nil {
		log.Fatalln(err)
	}
	if err := g.genMain(); err != nil {
		log.Fatalln(err)
	}
	return g
}

// formatModel
func (g *Generator) formatModel() {
	for k, _ := range g.modelBuf {
		formatedOutput, err := format.Source(g.modelBuf[k].Bytes())
		if err != nil {
			continue
			log.Fatalln(err)
		}
		g.modelBuf[k] = bytes.NewBuffer(formatedOutput)
	}
}

// formatSetting
func (g *Generator) formatSetting() {
	formatedOutput, err := format.Source(g.settingBuf[settingName].Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	g.settingBuf[settingName] = bytes.NewBuffer(formatedOutput)
}

// formatService
func (g *Generator) formatService() {
	for k, _ := range g.serviceBuf {
		formatedOutput, err := format.Source(g.serviceBuf[k].Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		g.serviceBuf[k] = bytes.NewBuffer(formatedOutput)
	}
}

// formatController
func (g *Generator) formatController() {
	for k, _ := range g.controlBuf {
		formatedOutput, err := format.Source(g.controlBuf[k].Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		g.controlBuf[k] = bytes.NewBuffer(formatedOutput)
	}
	//
	for k, _ := range g.routerBuf {
		formatedOutput, err := format.Source(g.routerBuf[k].Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		g.routerBuf[k] = bytes.NewBuffer(formatedOutput)
	}
}

// formatMain
func (g *Generator) formatMain() {
	formatedOutput, err := format.Source(g.mainBuf.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	g.mainBuf = bytes.NewBuffer(formatedOutput)
}

// Format function formates the output of the generation.
func (g *Generator) Format() *Generator {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		g.formatSetting()
	}()
	go func() {
		defer wg.Done()
		g.formatModel()
	}()
	go func() {
		defer wg.Done()
		g.formatService()
	}()

	go func() {
		defer wg.Done()
		g.formatController()
	}()
	wg.Wait()
	g.formatMain()
	return g
}

// flushModel
func (g *Generator) flushModel() (err error) {
	for k, _ := range g.modelBuf {
		s := gorm.ToDBName(k)
		dir := g.outputDir + "/app/model/" + s
		if err = os.MkdirAll(dir, 0777); err != nil {
			if !os.IsExist(err) {
				return
			}
		}
		filename := dir + "/gen_" + s + ".go"
		if utils.Exists(filename) && k != initDB {
			continue
		}

		if err = ioutil.WriteFile(filename, g.modelBuf[k].Bytes(), 0777); err != nil {
			return
		}
	}
	return
}

func (g *Generator) flushService() (err error) {
	for k, _ := range g.serviceBuf {
		s := gorm.ToDBName(k)
		dir := g.outputDir + "/app/service/" + s
		if err = os.MkdirAll(dir, 0777); err != nil {
			if !os.IsExist(err) {
				return
			}
		}
		filename := dir + "/gen_" + s + ".go"
		if utils.Exists(filename) {
			continue
		}

		if err = ioutil.WriteFile(filename, g.serviceBuf[k].Bytes(), 0777); err != nil {
			return
		}
	}
	return
}

// flushSetting
func (g *Generator) flushSetting() (err error) {
	for k, _ := range g.settingBuf {
		switch k {
		case settingName:
			dir := g.outputDir + "/app/setting"
			filename := dir + "/gen_setting.go"
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			if utils.Exists(filename) {
				continue
			}
			if err = ioutil.WriteFile(filename, g.settingBuf[k].Bytes(), 0777); err != nil {
				return
			}
		case envName:
			dir := g.outputDir + "/conf"
			filename := dir + "/.env"
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			if utils.Exists(filename) {
				continue
			}
			if err = ioutil.WriteFile(filename, g.settingBuf[k].Bytes(), 0777); err != nil {
				return
			}
		case devName, testName, prodName:
			dir := g.outputDir + "/conf/"
			filename := dir + k + "-global-config" + ".yaml"
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}

			if utils.Exists(filename) {
				continue
			}
			if err = ioutil.WriteFile(filename, g.settingBuf[k].Bytes(), 0777); err != nil {
				return
			}
		}
	}
	return
}

// flushController
func (g *Generator) flushController() (err error) {
	for k, _ := range g.controlBuf {
		s := gorm.ToDBName(k)
		var (
			dir string
		)
		switch k {
		case common, response, middleware:
			dir = g.outputDir + "/app/controller/" + s
		default:
			dir = g.outputDir + "/app/controller/api/" + s
		}

		if err = os.MkdirAll(dir, 0777); err != nil {
			if !os.IsExist(err) {
				return
			}
		}
		filename := dir + "/gen_" + s + ".go"
		if utils.Exists(filename) {
			continue
		}

		if err = ioutil.WriteFile(filename, g.controlBuf[k].Bytes(), 0777); err != nil {
			return
		}
	}
	// 生成router代码
	for k, _ := range g.routerBuf {
		s := gorm.ToDBName(k)
		var (
			dir string
		)
		dir = g.outputDir + "/app/controller/router/" + s
		if err = os.MkdirAll(dir, 0777); err != nil {
			if !os.IsExist(err) {
				return
			}
		}
		filename := dir + "/gen_" + s + ".go"
		if utils.Exists(filename) && k != initRouter {
			continue
		}

		if err = ioutil.WriteFile(filename, g.routerBuf[k].Bytes(), 0777); err != nil {
			return
		}
	}
	return
}

// flushMain
func (g *Generator) flushMain() (err error) {
	var (
		dir string
	)
	dir = g.outputDir
	if err = os.MkdirAll(dir, 0777); err != nil {
		if !os.IsExist(err) {
			return
		}
	}
	filename := dir + "/main" + ".go"
	if utils.Exists(filename) {
		return
	}

	if err = ioutil.WriteFile(filename, g.mainBuf.Bytes(), 0777); err != nil {
		return
	}

	return
}

// Flush function writes the output to the output file.
func (g *Generator) Flush() *Generator {
	var (
		err error
		wg  = sync.WaitGroup{}
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = g.flushModel(); err != nil {
			log.Fatalln(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = g.flushService(); err != nil {
			log.Fatalln(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = g.flushSetting(); err != nil {
			log.Fatalln(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = g.flushController(); err != nil {
			log.Fatalln(err)
		}
	}()
	wg.Wait()
	if err = g.flushMain(); err != nil {
		log.Fatalln(err)
	}
	return g
}

// GenDocs　gen swagger doc
func (g *Generator) GenDocs() {
	var (
		err error
	)
	if err = exec.Command("swag", "init", "-d", g.outputDir, "-o", g.outputDir+"/docs").Run(); err != nil {
		if strings.Contains(err.Error(), "executable file not found in $PATH") {
			if err = exec.Command("go", "get", "-u", "github.com/swaggo/swag/cmd/swag").Run(); err != nil {
				log.Fatalln(err)
			} else {
				if err = exec.Command("swag", "init", "-d", g.outputDir, "-o", g.outputDir+"/docs").Run(); err != nil {
					log.Fatalln(err)
				}
			}
		} else {
			log.Warnln(err)
		}
	}
	return
}
