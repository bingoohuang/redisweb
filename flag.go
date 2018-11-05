package main

import (
	"flag"
	"github.com/bingoohuang/go-utils"
	"strconv"
)

var (
	contextPath    string
	port           string
	maxContentSize int64

	devMode bool // to disable css/js minify
	servers []RedisServer

	maxKeys              int
	convenientConfigFile string
	argServers           string

	authParam go_utils.MustAuthParam
)

func init() {
	flag.StringVar(&contextPath, "contextPath", "", "context path")
	var portArg int
	flag.IntVar(&portArg, "port", 8269, "Port to serve.")
	flag.Int64Var(&maxContentSize, "maxContentSize", 10000, "max content size to display.")
	flag.BoolVar(&devMode, "devMode", false, "devMode(disable js/css minify)")

	flag.StringVar(&argServers, "servers", "default=localhost:6379", "servers list, eg: Server1=localhost:6379,Server2=password2/localhost:6388/0")
	flag.IntVar(&maxKeys, "maxKeys", 1000, "Max keys to be listed(0 means all keys).")
	flag.StringVar(&convenientConfigFile, "convenientConfigFile", "convenient-config.ini", "convenient-config.ini file path")
	go_utils.PrepareMustAuthFlag(&authParam)

	flag.Parse()

	port = strconv.Itoa(portArg)
	servers = parseServers(argServers)
}
