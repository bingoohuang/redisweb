package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	handleFunc(r, "/", serveHome, true)
	handleFunc(r, "/favicon.png", serveImage("favicon.png"), false)
	handleFunc(r, "/spritesheet.png", serveImage("spritesheet.png"), false)
	handleFunc(r, "/listKeys", serveListKeys, true)
	handleFunc(r, "/showContent", serveShowContent, true)
	handleFunc(r, "/downloadContent", downloadContent, true)
	handleFunc(r, "/changeContent", serveNewKey, false)
	handleFunc(r, "/deleteKey", serveDeleteKey, false)
	handleFunc(r, "/deleteMultiKeys", serveDeleteMultiKeys, false)
	handleFunc(r, "/exportKeys", serveExportKeys, true)
	handleFunc(r, "/newKey", serveNewKey, false)
	handleFunc(r, "/redisInfo", serveRedisInfo, true)
	handleFunc(r, "/redisCli", serveRedisCli, false)
	handleFunc(r, "/redisImport", serveRedisImport, false)
	handleFunc(r, "/convenientConfig", serveConvenientConfigRead, false)
	handleFunc(r, "/convenientConfigAdd", serveConvenientConfigAdd, false)
	handleFunc(r, "/deleteConvenientConfigItem", serveDeleteConvenientConfigItem, false)
	handleFunc(r, "/loadRedisServerConfig", serveLoadRedisServerConfig, false)
	handleFunc(r, "/saveRedisServerConfig", serveSaveRedisServerConfig, false)
	handleFunc(r, "/changeRedisServer", serveChangeRedisServer, false)

	http.Handle(appConfig.ContextPath+"/", r)

	fmt.Println("start to listen at ", port)
	go_utils.OpenExplorer(port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	select {} // 阻塞
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	wrap := go_utils.DumpRequest(f)

	if requiredGzip {
		wrap = go_utils.GzipHandlerFunc(wrap)
	}

	wrap = go_utils.MustAuth(wrap, authParam)

	r.HandleFunc(appConfig.ContextPath+path, wrap)
}
