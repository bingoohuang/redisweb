package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/go-utils"
	"log"
	"strconv"
	"strings"
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

var configFile string
var appConfig AppConfig
var authParam go_utils.MustAuthParam
var servers []RedisServer
var port string

func init() {
	flag.StringVar(&configFile, "configFile", "appConfig.toml", "config file path")

	flag.Parse()
	if _, err := toml.DecodeFile(configFile, &appConfig); err != nil {
		log.Panic("config file decode error", err.Error())
	}

	if appConfig.ContextPath != "" && strings.Index(appConfig.ContextPath, "/") < 0 {
		appConfig.ContextPath = "/" + appConfig.ContextPath
	}

	port = strconv.Itoa(appConfig.ListenPort)
	servers = parseServers(appConfig.Servers)

	authParam = go_utils.MustAuthParam{
		EncryptKey:  appConfig.EncryptKey,
		CookieName:  appConfig.CookieName,
		RedirectUri: appConfig.RedirectUri,
		LocalUrl:    appConfig.LocalUrl,
		ForceLogin:  appConfig.ForceLogin,
	}
}
