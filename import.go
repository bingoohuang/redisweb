package main

import (
	"net/http"
	"strconv"
	"strings"
)

func serveRedisImport(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	server := findRedisServer(req)
	commands := strings.TrimSpace(req.FormValue("commands"))
	commandItems := splitTrim(commands, "\n")

	for index, commandItem := range commandItems {
		result := repl(server, commandItem)
		_, _ = w.Write([]byte(strconv.Itoa(index+1) + ": "))
		_, _ = w.Write([]byte(result))
		_, _ = w.Write([]byte("\r\n"))
	}
}
