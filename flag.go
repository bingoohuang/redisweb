package main

import (
	"flag"
	"log"
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

	if _, err := toml.DecodeFile(configFile, &appConfig); err != nil {
		log.Panic("config file decode error", err.Error())
	}

	if appConfig.ContextPath != "" && strings.Index(appConfig.ContextPath, "/") < 0 {
		appConfig.ContextPath = "/" + appConfig.ContextPath
	}

	port = strconv.Itoa(appConfig.ListenPort)
	servers = parseServers(appConfig.Servers)
}
