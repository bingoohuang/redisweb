package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/pkger"
)

func downloadContent(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimSpace(req.FormValue("key"))
	server := findRedisServer(req)
	fileName := strings.TrimSpace(req.FormValue("fileName"))

	// tell the browser the returned content should be downloaded
	w.Header().Add("Content-Disposition", "Attachment; filename="+fileName)
	content := displayContent(server, key, false, true)
	http.ServeContent(w, req, fileName, time.Now(), bytes.NewReader([]byte(content.Content.(string))))
}

func serveShowContent(w http.ResponseWriter, req *http.Request) {
	HeadContentTypeJson(w)
	key := strings.TrimSpace(req.FormValue("key"))
	server := findRedisServer(req)

	if k, err := strconv.Unquote(key); err == nil {
		key = k
	}
	content := displayContent(server, key, true, false)
	_ = json.NewEncoder(w).Encode(content)
}

func serveNewKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	keyType := strings.TrimSpace(req.FormValue("type"))
	key := strings.TrimSpace(req.FormValue("key"))
	ttl := strings.TrimSpace(req.FormValue("ttl"))
	value := strings.TrimSpace(req.FormValue("value"))

	server := findRedisServer(req)

	// log.Println("keyType:", keyType, ",key:", key, ",ttl:", ttl, ",format:", format, ",value:", value)
	if k, err := strconv.Unquote(key); err == nil {
		key = k
	}

	ok := newKey(server, keyType, key, ttl, value)
	_, _ = w.Write([]byte(ok))
}

func serveRedisCli(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	server := findRedisServer(req)
	cmd := strings.TrimSpace(req.FormValue("cmd"))

	result := repl(server, cmd)
	_, _ = w.Write([]byte(result))
}

func serveRedisInfo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	server := findRedisServer(req)

	ok := redisInfo(server)
	_, _ = w.Write([]byte(ok))
}

func serveImage(image string) func(w http.ResponseWriter, r *http.Request) {
	data := MustAsset("/res/" + image)
	fi := AssetInfo("/res/" + image)
	return ServeImage(data, fi)
}

func ServeImage(imageBytes []byte, fi os.FileInfo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer := bytes.NewReader(imageBytes)
		w.Header().Set("Content-Type", DetectContentType(fi.Name()))
		w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, buffer)
	}
}

func DetectContentType(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); t == "" {
		t = "application/octet-stream"
	}
	return
}

func AssetInfo(name string) os.FileInfo {
	f, err := pkger.Open(name)
	if err != nil {
		return nil
	}

	defer f.Close()

	stat, _ := f.Stat()

	return stat
}
