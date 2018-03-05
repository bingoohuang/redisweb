package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/skratchdot/open-golang/open"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/gorilla/mux"
)

type RedisServer struct {
	ServerName string
	Addr       string
	Password   string
	DefaultDb  int
}

var (
	contextPath string
	port        string

	devMode    bool // to disable css/js minify
	argServers string
	servers    []RedisServer

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
	port = strconv.Itoa(*portArg)
	devMode = *devModeArg
	argServers = *serversArg
	servers = parseServers(argServers)
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

	redisServerConf := loadRedisServerConf()
	for key, val := range redisServerConf.Servers {
		result = append(result, RedisServer{
			ServerName: key,
			Addr:       val.Addr,
			Password:   val.Password,
			DefaultDb:  val.DefaultDb,
		})
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
			DefaultDb:  0,
		}
	} else if len == 2 {
		dbIndex, _ := strconv.Atoi(serverItems[1])
		return RedisServer{
			ServerName: serverName,
			Addr:       serverItems[0],
			Password:   "",
			DefaultDb:  dbIndex,
		}
	} else if len == 3 {
		dbIndex, _ := strconv.Atoi(serverItems[2])
		return RedisServer{
			ServerName: serverName,
			Addr:       serverItems[1],
			Password:   serverItems[0],
			DefaultDb:  dbIndex,
		}
	} else {
		panic("invalid servers argument")
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc(contextPath+"/", gzipWrapper(serveHome))
	r.HandleFunc(contextPath+"/favicon.png", serveImage("favicon.png"))
	r.HandleFunc(contextPath+"/spritesheet.png", serveImage("spritesheet.png"))
	r.HandleFunc(contextPath+"/listKeys", gzipWrapper(serveListKeys))
	r.HandleFunc(contextPath+"/showContent", gzipWrapper(serveShowContent))
	r.HandleFunc(contextPath+"/changeContent", serveNewKey)
	r.HandleFunc(contextPath+"/deleteKey", serveDeleteKey)
	r.HandleFunc(contextPath+"/deleteMultiKeys", serveDeleteMultiKeys)
	r.HandleFunc(contextPath+"/exportKeys", gzipWrapper(serveExportKeys))
	r.HandleFunc(contextPath+"/newKey", serveNewKey)
	r.HandleFunc(contextPath+"/redisInfo", gzipWrapper(serveRedisInfo))
	r.HandleFunc(contextPath+"/redisCli", serveRedisCli)
	r.HandleFunc(contextPath+"/redisImport", serveRedisImport)
	r.HandleFunc(contextPath+"/convenientConfig", serveConvenientConfigRead)
	r.HandleFunc(contextPath+"/convenientConfigAdd", serveConvenientConfigAdd)
	r.HandleFunc(contextPath+"/deleteConvenientConfigItem", serveDeleteConvenientConfigItem)
	r.HandleFunc(contextPath+"/loadRedisServerConfig", serveLoadRedisServerConfig)
	r.HandleFunc(contextPath+"/saveRedisServerConfig", serveSaveRedisServerConfig)
	r.HandleFunc(contextPath+"/changeRedisServer", serveChangeRedisServer)

	http.Handle(contextPath+"/", r)

	fmt.Println("start to listen at ", port)
	go openExplorer(port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func openExplorer(port string) {
	time.Sleep(100 * time.Millisecond)

	switch runtime.GOOS {
	case "windows":
		fallthrough
	case "darwin":
		open.Run("http://127.0.0.1:" + port)
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
	serverName := strings.TrimSpace(req.FormValue("server"))
	database := strings.TrimSpace(req.FormValue("database"))
	server := findServer(serverName)
	server.DefaultDb, _ = strconv.Atoi(database)
	return server
}

func findServer(serverName string) RedisServer {
	for _, server := range servers {
		if server.ServerName == serverName {
			fmt.Println("Found server ", serverName, server)
			return server
		}
	}

	fmt.Println("NotFound server ", serverName, ", user first ", servers[0])
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

const redisServerConfigFile = "redisServerConfig.toml"

func serveLoadRedisServerConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if _, err := os.Stat(redisServerConfigFile); os.IsNotExist(err) {
		json.NewEncoder(w).Encode(struct {
			RedisServerConfig string
		}{
			RedisServerConfig: `[servers]
    # [servers.demo1]
    # Address = "127.0.0.1:6379"
    # Password = ""
    # DefaultDb = 0

    # [servers.demo2]
    # Address = "127.0.0.1:7379"
    # Password = ""
    # DefaultDb = 0`,
		})
		return
	}

	redisServerConfig, _ := ioutil.ReadFile(redisServerConfigFile)
	json.NewEncoder(w).Encode(struct {
		RedisServerConfig string
	}{
		RedisServerConfig: string(redisServerConfig),
	})
}

type RedisServerConf struct {
	Servers map[string]RedisServer
}

func loadRedisServerConf() RedisServerConf {
	var redisServerConf RedisServerConf
	_, err := toml.DecodeFile(redisServerConfigFile, &redisServerConf)
	if err != nil {
		fmt.Println("DecodeFile error:", err)
	}

	return redisServerConf
}

func serveChangeRedisServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	redisServer := req.FormValue("redisServer")

	var foundServer *RedisServer = nil
	for _, server := range servers {
		if redisServer == server.ServerName {
			foundServer = &server
			break
		}
	}

	if foundServer != nil {
		dbs := configGetDatabases(*foundServer)
		json.NewEncoder(w).Encode(struct {
			OK        string
			DefaultDb int
			Dbs       int
		}{
			OK:        "OK",
			DefaultDb: foundServer.DefaultDb,
			Dbs:       dbs,
		})
	} else {
		json.NewEncoder(w).Encode(struct {
			OK string
		}{
			OK: "Server Unknown",
		})
	}
}

func serveSaveRedisServerConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	redisServerConfig := req.FormValue("redisServerConfig")

	var redisServerConf RedisServerConf
	_, err := toml.Decode(redisServerConfig, &redisServerConf)
	if err != nil {
		json.NewEncoder(w).Encode(struct {
			OK string
		}{
			OK: err.Error(),
		})
		return
	}

	err = ioutil.WriteFile(redisServerConfigFile, []byte(redisServerConfig), 0644)
	if err != nil {
		json.NewEncoder(w).Encode(struct {
			OK string
		}{
			OK: err.Error(),
		})
		return
	}

	loadRedisServerConf()
	servers = parseServers(argServers)

	serverNames := make([]string, 0)
	for _, server := range servers {
		serverNames = append(serverNames, server.ServerName)
	}

	dbs := 0
	defaultDb := 0

	if len(servers) > 0 {
		defaultDb = servers[0].DefaultDb
		dbs = configGetDatabases(servers[0])
	}

	json.NewEncoder(w).Encode(struct {
		OK        string
		Servers   []string
		DefaultDb int
		Dbs       int
	}{
		OK:        "OK",
		Servers:   serverNames,
		DefaultDb: defaultDb,
		Dbs:       dbs,
	})
}
