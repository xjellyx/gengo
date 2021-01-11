package setting

import "fmt"

var (
	SettingTemplate = fmt.Sprintf(`package setting

import(
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"github.com/olongfen/contrib/log"
	"os"
	"strings"
)

// Settings global config
type Settings struct {
	FilePath FilePath
	DB       Database
	Serve    Serve
	Project  Project
}

// Project project config
type Project struct {
}

// Serve serve config
type Serve struct {
	ServerAddr string
	ServerPort string
	IsTLS      bool
	TLS        TLS
}

type TLS struct {
	Cert string
	Key  string
}

// Database database config
type Database struct {
	Host         string
	Port         string
	Driver       string
	DatabaseName string
	Username     string
	Password     string
	MaxIdleConn  int
	MaxOpenConn  int
 	Source string
}

// FilePath file save path
type FilePath struct {
	LogDir    string // log save dir
	LogPatent string // log patent
}

var (
	Global = new(Settings)
	DevEnv       = false
)

func init() {
	var (
		err        error
		configFile string
	)
	if err = gotenv.Load("./conf/.env"); err != nil {
		log.Fatal(err)
	}
	env := os.Getenv("ENVIRONMENT")
	switch {
	case strings.Contains(env, "prod"):
		configFile = "./conf/prod-global-config.yaml"
	case strings.Contains(env, "test"):
		configFile = "./conf/test-global-config.yaml"
	default: // default is dev
		DevEnv = true
		configFile = "./conf/dev-global-config.yaml"

	}
	viper.SetConfigFile(configFile)
	_ = viper.ReadInConfig()
	if err = viper.Unmarshal(Global); err != nil {
		log.Fatal(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file: ",e.Name," Op: ",  e.Op)
		if err = viper.Unmarshal(Global); err != nil {
			log.Fatal(err)
		}
	})
	log.Infoln("setting init success !")
}

`)
	EnvTemplate    = `ENVIRONMENT=dev`
	ConfigTemplate = `db:
  databasename: postgres
  driver: postgres
  host: 127.0.0.1
  maxidleconn: 200
  maxopenconn: 2000
  password: 123456
  port: "5432"
  username: postgres
filepath:
  logdir: ./log
  logpatent: _%Y-%m-%d.log_project
serve:
  serveraddr: 0.0.0.0
  serverport: "8080"
project:

`
)
