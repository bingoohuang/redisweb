package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bingoohuang/gou/htt"
	"github.com/markbates/pkger"
)

func MustAsset(name string) []byte {
	f, err := pkger.Open(name)
	if err != nil {
		return nil
	}

	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}

	return content
}

func AssetNames() []string {
	var names []string

	f, err := pkger.Open("/res")
	if err != nil {
		return names
	}

	fis, _ := f.Readdir(0)

	for _, fi := range fis {
		if !fi.IsDir() {
			names = append(names, "/res/"+fi.Name())
		}
	}

	return names
}

func serveJsStatic(w http.ResponseWriter, r *http.Request) {
	dir := http.FileServer(pkger.Dir("/res"))
	ServeStaticFolder(w, r, dir, "/static")
}

func serveCssStatic(w http.ResponseWriter, r *http.Request) {
	dir := http.FileServer(pkger.Dir("/res"))
	ServeStaticFolder(w, r, dir, "/static")
}

func ServeStaticFolder(w http.ResponseWriter, r *http.Request, dir http.Handler, prefix string) {
	switch prefix {
	case "", "/": // ignore
	default:
		dir = http.StripPrefix(prefix, dir)
	}

	dir.ServeHTTP(w, r)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := string(MustAsset("/res/index.html"))
	html = strings.Replace(html, "<serverOptions/>", serverOptions(), 1)
	html = strings.Replace(html, "<databaseOptions/>", databaseOptions(), 1)
	html = htt.MinifyHTML(html, appConfig.DevMode)

	assetNames := AssetNames()
	mergeCss := htt.MergeCSS(MustAsset, htt.FilterAssetNames(assetNames, ".css"))
	css := htt.MinifyCSS(mergeCss, appConfig.DevMode)
	mergeScripts := htt.MergeJs(MustAsset, htt.FilterAssetNames(assetNames, ".js"))
	js := htt.MinifyJs(mergeScripts, appConfig.DevMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${ContextPath}", appConfig.ContextPath, -1)
	w.Write([]byte(html))
}

func databaseOptions() string {
	options := ""

	server0 := servers[0]
	databases := configGetDatabases(server0)
	for i := 0; i < databases; i++ {
		databaseIndex := strconv.Itoa(i)
		if server0.DefaultDb == i {
			options += `<option value="` + databaseIndex + `" selected>` + databaseIndex + `</option>`
		} else {
			options += `<option value="` + databaseIndex + `">` + databaseIndex + `</option>`
		}
	}

	return options
}

func serverOptions() string {
	options := ""

	for _, server := range servers {
		options += `<option value="` + server.ServerName + `">` + server.ServerName + `</option>`
	}

	return options
}
