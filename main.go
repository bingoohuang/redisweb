package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"

	"github.com/bingoohuang/gou/ran"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	r := mux.NewRouter()

	handleFunc(r, "/", serveHome, true)
	handleFunc(r, "/static/css/{key}", serveCssStatic, true)
	handleFunc(r, "/static/js/{key}", serveJsStatic, true)
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
	OpenExplorer(port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	select {} // 阻塞
}

func OpenExplorerWithContext(contextPath, port string) {
	go func() {
		time.Sleep(100 * time.Millisecond)

		switch runtime.GOOS {
		case "windows":
			fallthrough
		case "darwin":
			open.Run("http://127.0.0.1:" + port + contextPath + "/?" + ran.String(10))
		}
	}()
}

func OpenExplorer(port string) {
	OpenExplorerWithContext("", port)
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	wrap := DumpRequest(f)

	if requiredGzip {
		wrap = GzipHandlerFunc(wrap)
	}

	r.HandleFunc(appConfig.ContextPath+path, wrap)
}

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := GzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func DumpRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))
		fn(w, r)
	}
}
