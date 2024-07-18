package main

import (
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"
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

	http.Handle("/", r)

	fmt.Println("start to listen at ", port)
	OpenExplorerWithContext(appConfig.ContextPath, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func OpenExplorerWithContext(contextPath, port string) {
	go func() {
		time.Sleep(100 * time.Millisecond)

		switch runtime.GOOS {
		case "windows":
			fallthrough
		case "darwin":
			open.Run("http://127.0.0.1:" + port + contextPath + "?" + ran.String(10))
		}
	}()
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	//wrap := DumpRequest(f)
	wrap := f

	if requiredGzip {
		wrap = GzipHandlerFunc(wrap)
	}

	p := filepath.Join(appConfig.ContextPath, path)
	if p != "/" {
		p = strings.TrimSuffix(p, "/")
	}

	if appConfig.BasicAuth != "" {
		wrap = basicAuth(wrap, appConfig.BasicAuth)
	}

	r.HandleFunc(p, wrap)
}

// AsciiEqualFold is [strings.EqualFold], ASCII only. It reports whether s and t
// are equal, ASCII-case-insensitively.
func AsciiEqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lower(s[i]) != lower(t[i]) {
			return false
		}
	}
	return true
}

// lower returns the ASCII lowercase version of b.
func lower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (usernamePassword string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !AsciiEqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", false
	}
	return string(c), true
}

func basicAuth(next http.HandlerFunc, basicUserPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			realAuth, _ := parseBasicAuth(auth)
			if realAuth == basicUserPassword {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
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
