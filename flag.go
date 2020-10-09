package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	ContextPath string
	ListenPort  int

	EncryptKey  string
	CookieName  string
	RedirectUri string
	LocalUrl    string
	ForceLogin  bool

	DevMode bool

	Servers              string
	MaxContentSize       int64
	MaxKeys              int
	ConvenientConfigFile string
}

var (
	configFile string
	appConfig  AppConfig
	servers    []RedisServer
	port       string
)

func init() {
	flag.StringVar(&configFile, "config", "redisweb.toml", "config file path")
	flag.Parse()

	createDefaultRedisWebConfigFile()

	if _, err := toml.DecodeFile(configFile, &appConfig); err != nil {
		log.Panic("config file decode error", err.Error())
	}

	createDefaultConvenientConfigFile()

	if appConfig.ContextPath != "" && strings.Index(appConfig.ContextPath, "/") < 0 {
		appConfig.ContextPath = "/" + appConfig.ContextPath
	}

	port = strconv.Itoa(appConfig.ListenPort)
	servers = parseServers(appConfig.Servers)
}

func createDefaultConvenientConfigFile() {
	if appConfig.ConvenientConfigFile == "" {
		appConfig.ConvenientConfigFile = "convenient.ini"
	}

	if _, err := os.Stat(appConfig.ConvenientConfigFile); os.IsNotExist(err) {
		ioutil.WriteFile(appConfig.ConvenientConfigFile, []byte(`
[convenient1]
name       = 登录验证码
template   = captcha:{mobile}:/login
operations = save,delete
ttl        = 15s

[convenient2]
name       = 课种缓存
template   = westcache:yoga:{tcode}:CourseTypeService.queryCourseTypes
operations = delete

[convenient3]
name       = 强制退出
template   = shiro:session-operation:{userId}
operations = save

[blNTRhBqDP]
name       = 日志跟踪手机号
template   = TraceLogMobile:{mobile}
operations = Delete,Save
ttl        = -1s
`), 0644)
	}
}

func createDefaultRedisWebConfigFile() {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		ioutil.WriteFile(configFile, []byte(`
#ContextPath = ""
ListenPort  = 8269
# max content size to display.
MaxContentSize = 10000
DevMode = false
# servers list, eg: Server1=localhost:6379,Server2=password2/localhost:6388/0
Servers = "default=localhost:6379"
# Max keys to be listed(0 means all keys).
MaxKeys = 1000
# convenient.ini file path
ConvenientConfigFile = "convenient.ini"
`), 0644)
	}
}
