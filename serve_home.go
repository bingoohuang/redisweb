package main

import (
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strconv"
	"strings"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := string(MustAsset("res/index.html"))
	html = strings.Replace(html, "<serverOptions/>", serverOptions(), 1)
	html = strings.Replace(html, "<databaseOptions/>", databaseOptions(), 1)
	html = go_utils.MinifyHtml(html, *devMode)

	mergeCss := go_utils.MergeCss(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".css"))
	css := go_utils.MinifyCss(mergeCss, *devMode)
	mergeScripts := go_utils.MergeJs(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".js"))
	js := go_utils.MinifyJs(mergeScripts, *devMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${ContextPath}", *contextPath, -1)
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
