package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strings"
)

func serveShowContent(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	key := strings.TrimSpace(req.FormValue("key"))
	server := findRedisServer(req)

	content := displayContent(server, key)
	json.NewEncoder(w).Encode(content)
}

func serveNewKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	keyType := strings.TrimSpace(req.FormValue("type"))
	key := strings.TrimSpace(req.FormValue("key"))
	ttl := strings.TrimSpace(req.FormValue("ttl"))
	value := strings.TrimSpace(req.FormValue("value"))

	server := findRedisServer(req)

	// log.Println("keyType:", keyType, ",key:", key, ",ttl:", ttl, ",format:", format, ",value:", value)

	ok := newKey(server, keyType, key, ttl, value)
	w.Write([]byte(ok))
}

func serveRedisCli(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	server := findRedisServer(req)
	cmd := strings.TrimSpace(req.FormValue("cmd"))

	result := repl(server, cmd)
	w.Write([]byte(result))
}

func serveRedisInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	server := findRedisServer(req)

	ok := redisInfo(server)
	w.Write([]byte(ok))
}

func serveImage(image string) func(w http.ResponseWriter, r *http.Request) {
	data := MustAsset("res/" + image)
	fi, _ := AssetInfo("res/" + image)
	return go_utils.ServeImage(data, fi)
}
