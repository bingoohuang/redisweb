package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type RedisServer struct {
	ServerName string
	Addr       string
	Password   string
	DB         int
}

var (
	contextPath string
	port        int

	devMode bool // to disable css/js minify
	servers []RedisServer

	maxKeys              int
	convenientConfigFile string
)

func init() {
	contextPathArg := flag.String("contextPath", "", "context path")
	portArg := flag.Int("port", 8269, "Port to serve.")
	devModeArg := flag.Bool("devMode", false, "devMode(disable js/css minify)")
	serversArg := flag.String("servers", "default=localhost:6379", "servers list, eg: Server1=localhost:6379,Server2=password2/localhost:6388/0")
	maxKeysArg := flag.Int("maxKeys", 1000, "Max keys to be listed(0 means all keys).")
	convenientConfigFileArg := flag.String("convenientConfigFile", "convenient-config.ini", "convenient-config.ini file path")

	flag.Parse()

	contextPath = *contextPathArg
	port = *portArg
	devMode = *devModeArg
	servers = parseServers(*serversArg)
	maxKeys = *maxKeysArg
	convenientConfigFile = *convenientConfigFileArg
}

func parseServers(serversConfig string) []RedisServer {
	serverItems := splitTrim(serversConfig, ",")

	var result = make([]RedisServer, 0)
	for i, item := range serverItems {
		parts := splitTrim(item, "=")
		len := len(parts)
		if len == 1 {
			serverName := "Server" + strconv.Itoa(i+1)
			result = append(result, parseServerItem(serverName, parts[0]))
		} else if len == 2 {
			serverName := parts[0]
			result = append(result, parseServerItem(serverName, parts[1]))
		} else {
			panic("invalid servers argument")
		}
	}

	return result
}

func splitTrim(str, sep string) []string {
	subs := strings.Split(str, sep)
	ret := make([]string, 0)
	for i, v := range subs {
		v := strings.TrimSpace(v)
		if len(subs[i]) > 0 {
			ret = append(ret, v)
		}
	}

	return ret
}

func parseServerItem(serverName, serverConfig string) RedisServer {
	serverItems := splitTrim(serverConfig, "/")
	len := len(serverItems)
	if len == 1 {
		return RedisServer{
			ServerName: serverName,
			Addr:       serverItems[0],
			Password:   "",
			DB:         0,
		}
	} else if len == 2 {
		dbIndex, _ := strconv.Atoi(serverItems[1])
		return RedisServer{
			ServerName: serverName,
			Addr:       serverItems[0],
			Password:   "",
			DB:         dbIndex,
		}
	} else if len == 3 {
		dbIndex, _ := strconv.Atoi(serverItems[2])
		return RedisServer{
			ServerName: serverName,
			Addr:       serverItems[1],
			Password:   serverItems[0],
			DB:         dbIndex,
		}
	} else {
		panic("invalid servers argument")
	}
}

func main() {
	http.HandleFunc(contextPath+"/", gzipWrapper(serveHome))
	http.HandleFunc(contextPath+"/favicon.png", serveImage("favicon.png"))
	http.HandleFunc(contextPath+"/spritesheet.png", serveImage("spritesheet.png"))
	http.HandleFunc(contextPath+"/listKeys", gzipWrapper(serveListKeys))
	http.HandleFunc(contextPath+"/showContent", gzipWrapper(serveShowContent))
	http.HandleFunc(contextPath+"/changeContent", serveNewKey)
	http.HandleFunc(contextPath+"/deleteKey", serveDeleteKey)
	http.HandleFunc(contextPath+"/deleteMultiKeys", serveDeleteMultiKeys)
	http.HandleFunc(contextPath+"/exportKeys", gzipWrapper(serveExportKeys))
	http.HandleFunc(contextPath+"/newKey", serveNewKey)
	http.HandleFunc(contextPath+"/redisInfo", gzipWrapper(serveRedisInfo))
	http.HandleFunc(contextPath+"/redisCli", serveRedisCli)
	http.HandleFunc(contextPath+"/redisImport", serveRedisImport)
	http.HandleFunc(contextPath+"/convenientConfig", serveConvenientConfig)

	sport := strconv.Itoa(port)
	fmt.Println("start to listen at ", sport)
	if err := http.ListenAndServe(":"+sport, nil); err != nil {
		log.Fatal(err)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func serveImage(image string) func(w http.ResponseWriter, r *http.Request) {
	path := "res/" + image
	data := MustAsset(path)

	return func(w http.ResponseWriter, r *http.Request) {
		fi, _ := AssetInfo(path)
		buffer := bytes.NewReader(data)
		w.Header().Set("Content-Type", detectContentType(fi.Name()))
		w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, buffer)
	}
}

func detectContentType(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); t == "" {
		t = "application/octet-stream"
	}
	return
}

func serveListKeys(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	matchPattern := strings.TrimSpace(req.FormValue("pattern"))
	server := findRedisServer(req)

	keys, _ := listKeys(server, matchPattern, maxKeys)
	sort.Slice(keys[:], func(i, j int) bool {
		return keys[i].Key < keys[j].Key
	})
	json.NewEncoder(w).Encode(keys)
}

func findRedisServer(req *http.Request) RedisServer {
	serverName := strings.TrimSpace(req.FormValue("serverName"))
	database := strings.TrimSpace(req.FormValue("database"))
	server := findServer(serverName)
	server.DB, _ = strconv.Atoi(database)
	return server
}

func findServer(serverName string) RedisServer {
	for _, server := range servers {
		if server.ServerName == serverName {
			return server
		}
	}

	return servers[0]
}

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

	//log.Println("keyType:", keyType, ",key:", key, ",ttl:", ttl, ",format:", format, ",value:", value)

	ok := newKey(server, keyType, key, ttl, value)
	w.Write([]byte(ok))
}

func serveDeleteKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	key := strings.TrimSpace(req.FormValue("key"))
	server := findRedisServer(req)

	ok := deleteMultiKeys(server, key)
	w.Write([]byte(ok))
}

func serveExportKeys(w http.ResponseWriter, req *http.Request) {

	server := findRedisServer(req)
	exportKeys := strings.TrimSpace(req.FormValue("exportKeys"))
	exportType := strings.TrimSpace(req.FormValue("exportType"))
	download := strings.TrimSpace(req.FormValue("download"))
	if download == "true" {
		if exportType == "JSON" {
			w.Header().Set("Content-Disposition", "attachment; filename=export.json")
		} else {
			w.Header().Set("Content-Disposition", "attachment; filename=export.txt")
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
	}

	result := exportRedisKeys(server, exportKeys, exportType)
	switch result := result.(type) {
	case map[string]interface{}:
		jsonResult, _ := json.Marshal(result)
		str := jsonPrettyPrint(string(jsonResult))

		if download == "true" {
			io.Copy(w, strings.NewReader(str))
		} else {
			json.NewEncoder(w).Encode(str)
		}
	case []string:
		if download == "true" {
			joined := strings.Join(result, "\r\n")
			io.Copy(w, strings.NewReader(joined))
		} else {
			json.NewEncoder(w).Encode(result)
		}
	}
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

func serveRedisImport(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	server := findRedisServer(req)
	commands := strings.TrimSpace(req.FormValue("commands"))
	commandItems := splitTrim(commands, "\n")

	for index, commandItem := range commandItems {
		result := repl(server, commandItem)
		w.Write([]byte(strconv.Itoa(index+1) + ": "))
		w.Write([]byte(result))
		w.Write([]byte("\r\n"))
	}
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
