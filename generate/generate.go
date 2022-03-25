package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/olongfen/gengo/template/controller"
	"github.com/olongfen/gengo/template/entity"
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
	entities    = "entities"
)

type Struct struct {
	LowerName   string
	StructName  string
	PackageName string
}

type InitDB struct {
	Structs  []Struct
	Mod      string
	Separate bool
}

// The Generator is the one responsible for generating the code, adding the imports, formating, and writing it to the file.
type Generator struct {
	modelBuf   map[string]*bytes.Buffer
	formBuf    map[string]*bytes.Buffer
	serviceBuf map[string]*bytes.Buffer
	settingBuf map[string]*bytes.Buffer
	controlBuf map[string]*bytes.Buffer
	routerBuf  map[string]*bytes.Buffer
	entity     map[string]*bytes.Buffer
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
	g.entity[entities] = &bytes.Buffer{}
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
		formBuf:    map[string]*bytes.Buffer{},
		serviceBuf: map[string]*bytes.Buffer{},
		settingBuf: map[string]*bytes.Buffer{},
		controlBuf: map[string]*bytes.Buffer{},
		routerBuf:  map[string]*bytes.Buffer{},
		mainBuf:    &bytes.Buffer{},
		outputDir:  output,
		parser:     p,
		initDB:     new(InitDB),
		entity:     map[string]*bytes.Buffer{},
	}
	g.config = c
	g.init()
	g.initDB.Mod = c.Mod
	g.initDB.Separate = c.Separate
	for _, v := range p.Structs {
		g.initDB.Structs = append(g.initDB.Structs, Struct{
			LowerName:   v.LowerName,
			StructName:  v.StructName,
			PackageName: v.PackageName,
		})
		v.Config = c
		v.Config.Package = v.LowerName
		g.modelBuf[v.StructName] = new(bytes.Buffer)
		g.entity[v.StructName] = new(bytes.Buffer)
		g.serviceBuf[v.StructName] = new(bytes.Buffer)
		g.controlBuf[v.StructName] = new(bytes.Buffer)
		g.routerBuf[v.StructName] = new(bytes.Buffer)
	}

	return g, nil
}

// GenEntity gen
func (g *Generator) GenEntity() {
	var (
		err  error
		file *os.File
		temp *template.Template
	)
	if temp, err = template.New("entity").Parse(entity.Template.String()); err != nil {
		log.Fatalln(err)
	}
	for _, v := range g.parser.Structs {
		if _, ok := g.entity[v.StructName]; !ok {
			continue
		}

		if err = temp.Execute(g.entity[v.StructName], v); err != nil {
			log.Fatalln(err)
		}

		if file, err = os.OpenFile(g.config.Input, os.O_WRONLY|os.O_APPEND, os.ModeAppend); err != nil {
			log.Fatalln(err)
		}
		writer := bufio.NewWriter(file)
		pak := []byte("package tmpl")
		all := append(pak, g.entity[v.StructName].Bytes()...)

		data := g.parser.CheckExistFunc(g.config.Input, all)
		if len(data) == len(all) {
			data = data[len(pak):]
		}
		// 已经生成或者存在Set Get方法不会再添加生成
		if _, err = writer.Write(data); err != nil {
			log.Fatalln(err)
		}
		if err = writer.Flush(); err != nil {
			log.Fatalln(err)
		}

	}
	return
}

// genModel
func (g *Generator) genModel() (err error) {
	var (
		temp *template.Template
	)
	// 生成model公共代码
	if temp, err = template.New(common).Parse(model.TemplateCommon.String()); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	if !c.Separate {
		c.Package = "models"
	} else {
		c.Package = "model_common"
	}
	if err = temp.Execute(g.modelBuf[common], c); err != nil {
		return
	}

	// 生成初始化数据库代码
	if temp, err = template.New(initDB).Parse(model.TemplateGORMInitDB.String()); err != nil {
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
		g.formBuf[v.StructName] = new(bytes.Buffer)
		if t, err = template.New(v.StructName).Parse(model.TemplateGORMForm.String()); err != nil {
			return
		}
		if err = t.Execute(g.formBuf[v.StructName], v); err != nil {
			return
		}
		switch g.config.ORM {
		case "gorm":
			if t, err = template.New(v.StructName).Parse(model.TemplateGORM.String()); err != nil {
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

	if temp, err = template.New(settingName).Parse(setting.TemplateSetting.String()); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[settingName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(envName).Parse(setting.TemplateEnv); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[envName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(devName).Parse(setting.TemplateConfig.String()); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[devName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(testName).Parse(setting.TemplateConfig.String()); err != nil {
		return
	}
	if err = temp.Execute(g.settingBuf[testName], nil); err != nil {
		return
	}

	//
	if temp, err = template.New(prodName).Parse(setting.TemplateConfig.String()); err != nil {
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
	if temp, err = template.New(common).Parse(service.TemplateCommon.String()); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	if !c.Separate {
		c.Package = "services"
	} else {
		c.Package = "svc_common"
	}
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
			if t, err = template.New(v.StructName).Parse(service.TemplateGorm.String()); err != nil {
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
	if temp, err = template.New(common).Parse(controller.TemplateCommon.String()); err != nil {
		return
	}
	c := parse.Config{}
	c = g.config
	if !c.Separate {
		c.Package = "apis"
	} else {
		c.Package = "ctl_common"
	}

	if err = temp.Execute(g.controlBuf[common], c); err != nil {
		return
	}
	// 初始化代码
	switch g.config.WEB {
	case "gin":
		// response
		if temp, err = template.New(response).Parse(controller.TemplateResponse.String()); err != nil {
			return
		}
		if err = temp.Execute(g.controlBuf[response], nil); err != nil {
			return
		}

		// middleware
		if temp, err = template.New(middleware).Parse(controller.TemplateMiddleware.String()); err != nil {
			return
		}
		if err = temp.Execute(g.controlBuf[middleware], nil); err != nil {
			return
		}

		// init router
		if temp, err = template.New(initRouter).Parse(controller.TemplateInitRouter.String()); err != nil {
			return
		}
		rc := struct {
			Structs []*parse.StructData
			parse.Config
		}{}
		rc.Config = g.config
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
			if t, err = template.New(v.StructName).Parse(controller.TemplateGin.String()); err != nil {
				return
			}
			if err = t.Execute(g.controlBuf[v.StructName], v); err != nil {
				log.Fatalln(err)
			}

			//
			if t, err = template.New(v.StructName).Parse(controller.TemplateStructRouter.String()); err != nil {
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
		if temp, err = template.New(common).Parse(gin_main.TemplateGinMain.String()); err != nil {
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

func (g *Generator) formatForm() {
	for k, _ := range g.formBuf {
		formatedOutput, err := format.Source(g.formBuf[k].Bytes())
		if err != nil {
			continue
			log.Fatalln(err)
		}
		g.formBuf[k] = bytes.NewBuffer(formatedOutput)
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
			log.Warnln(string(g.serviceBuf[k].Bytes()))
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
			log.Warnln(string(g.controlBuf[k].Bytes()))
			log.Fatalln(err)
		}
		g.controlBuf[k] = bytes.NewBuffer(formatedOutput)
	}
	//
	for k, _ := range g.routerBuf {
		formatedOutput, err := format.Source(g.routerBuf[k].Bytes())
		if err != nil {
			log.Warnln(string(g.routerBuf[k].Bytes()))
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
		g.formatForm()
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
		var (
			filename string
			s        = gorm.ToDBName(k)
		)
		if g.config.Separate {
			dir := g.outputDir + "/app/models/"
			dir += utils.SQLColumn2PkgStyle(s)
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = dir + "/" + s + ".go"
		} else {
			if err = os.MkdirAll(g.outputDir+"/app/models/", 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = g.outputDir + "/app/models/" + s + ".go"
		}

		if utils.Exists(filename) {
			continue
		}
		if err = ioutil.WriteFile(filename, g.modelBuf[k].Bytes(), 0777); err != nil {
			return
		}
	}
	return
}

// flushModel
func (g *Generator) flushForm() (err error) {
	for k, _ := range g.formBuf {
		var (
			filename string
			s        = gorm.ToDBName(k)
		)
		if g.config.Separate {
			dir := g.outputDir + "/app/form/"
			dir += utils.SQLColumn2PkgStyle(s)
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = dir + "/" + s + ".go"
		} else {
			if err = os.MkdirAll(g.outputDir+"/app/form/", 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = g.outputDir + "/app/form/" + s + ".go"
		}

		if utils.Exists(filename) {
			continue
		}
		if err = ioutil.WriteFile(filename, g.formBuf[k].Bytes(), 0777); err != nil {
			return
		}
	}
	return
}

func (g *Generator) flushService() (err error) {
	for k, _ := range g.serviceBuf {
		var (
			filename string
			s        = gorm.ToDBName(k)
		)
		if g.config.Separate {

			dir := g.outputDir + "/app/services/" + utils.SQLColumn2PkgStyle(s)
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = dir + "/" + s + ".go"
		} else {
			dir := g.outputDir + "/app/services/"
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = g.outputDir + "/app/services/" + s + ".go"
		}
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
			filename := dir + "/setting.go"
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
		var (
			s        = gorm.ToDBName(k)
			filename string
		)
		var (
			dir    string
			prefix string
		)
		switch k {
		case response, middleware:
			dir = g.outputDir + "/app/controller/" + utils.SQLColumn2PkgStyle(s)
			prefix = "/"
			if err = os.MkdirAll(dir, 0777); err != nil {
				if !os.IsExist(err) {
					return
				}
			}
			filename = dir + prefix + s + ".go"
		default:
			if k == common {
				dir = g.outputDir + "/app/controller/" + s + "/"
			} else {
				dir = g.outputDir + "/app/controller/apis/" + s
				prefix = "/"
			}
			if g.config.Separate {
				if err = os.MkdirAll(dir, 0777); err != nil {
					if !os.IsExist(err) {
						return
					}
				}
				filename = dir + prefix + s + ".go"
			} else {
				if err = os.MkdirAll(g.outputDir+"/app/controller/apis/", 0777); err != nil {
					if !os.IsExist(err) {
						return
					}
				}
				filename = g.outputDir + "/app/controller/apis/" + s + ".go"
			}
		}

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
			filename string
			dir      string
		)
		dir = g.outputDir + "/app/controller/router/"
		if err = os.MkdirAll(dir, 0777); err != nil {
			if !os.IsExist(err) {
				return
			}
		}
		filename = dir + s + ".go"
		if utils.Exists(filename) {
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
		if err = g.flushForm(); err != nil {
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
