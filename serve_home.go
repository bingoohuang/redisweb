package main

import (
	"bytes"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"net/http"
	"strconv"
	"strings"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != contextPath+"/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := string(MustAsset("res/index.html"))
	html = strings.Replace(html, "<serverOptions/>", serverOptions(), 1)
	html = strings.Replace(html, "<databaseOptions/>", databaseOptions(), 1)
	html = minifyHtml(html, devMode)

	css, js := minifyCssJs(mergeCss(), mergeScripts(), devMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)

	w.Write([]byte(html))
}

func minifyHtml(htmlStr string, devMode bool) string {
	if devMode {
		return htmlStr
	}

	mini := minify.New()
	mini.AddFunc("text/html", html.Minify)
	minified, _ := mini.String("text/html", htmlStr)
	return minified
}

func minifyCssJs(mergedCss, mergedJs string, devMode bool) (string, string) {
	if devMode {
		return mergedCss, mergedJs
	}

	mini := minify.New()
	mini.AddFunc("text/css", css.Minify)
	mini.AddFunc("text/javascript", js.Minify)

	minifiedCss, _ := mini.String("text/css", mergedCss)
	minifiedJs, _ := mini.String("text/javascript", mergedJs)

	return minifiedCss, minifiedJs
}

func databaseOptions() string {
	options := ""

	databases := configGetDatabases(servers[0])
	for i := 0; i < databases; i++ {
		databaseIndex := strconv.Itoa(i)
		options += `<option value="` + databaseIndex + `">` + databaseIndex + `</option>`
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

func mergeCss() string {
	return mergeStatic("stylesheet.css", "codemirror-5.29.0.min.css", "index.css")
}

func mergeScripts() string {
	return mergeStatic("jquery-3.2.1.min.js", "codemirror-5.29.0.min.js", "matchbrackets-5.29.0.min.js",
		"javascript-5.29.0.min.js", "autosize-4.0.0.min.js", "js.cookie.js", "index.js", "resizebar.js")
}

func mergeStatic(statics ...string) string {
	var scripts bytes.Buffer
	for _, static := range statics {
		scripts.Write(MustAsset("res/" + static))
		scripts.Write([]byte("\n"))
	}

	return scripts.String()
}