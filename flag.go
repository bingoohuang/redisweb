package main

import (
	"flag"
	"github.com/bingoohuang/go-utils"
	"strconv"
)

var (
	contextPath *string
	port        string

	devMode    *bool // to disable css/js minify
	argServers *string
	servers    []RedisServer

	maxKeys              *int
	convenientConfigFile *string

	authParam go_utils.MustAuthParam
)

func init() {
	contextPath = flag.String("contextPath", "", "context path")
	portArg := flag.Int("port", 8269, "Port to serve.")
	devMode = flag.Bool("devMode", false, "devMode(disable js/css minify)")
	argServers = flag.String("servers", "default=localhost:6379", "servers list, eg: Server1=localhost:6379,Server2=password2/localhost:6388/0")
	maxKeys = flag.Int("maxKeys", 1000, "Max keys to be listed(0 means all keys).")
	convenientConfigFile = flag.String("convenientConfigFile", "convenient-config.ini", "convenient-config.ini file path")
	go_utils.PrepareMustAuthFlag(&authParam)

	flag.Parse()

	port = strconv.Itoa(*portArg)
	servers = parseServers(*argServers)
}
