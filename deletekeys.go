package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func serveDeleteKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	key := strings.TrimSpace(req.FormValue("key"))
	server := findRedisServer(req)

	ok := deleteMultiKeys(server, key)
	w.Write([]byte(ok))
}

func serveDeleteMultiKeys(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	keys := strings.TrimSpace(req.FormValue("keys"))
	server := findRedisServer(req)

	var multiKeys []string
	json.Unmarshal([]byte(keys), &multiKeys)

	ok := deleteMultiKeys(server, multiKeys...)
	w.Write([]byte(ok))
}
