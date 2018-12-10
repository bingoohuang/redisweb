package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net/http"
	"os"
)

const redisServerConfigFile = "redisServerConfig.toml"

func serveLoadRedisServerConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if _, err := os.Stat(redisServerConfigFile); os.IsNotExist(err) {
		_ = json.NewEncoder(w).Encode(struct {
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
	_ = json.NewEncoder(w).Encode(struct {
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
	if _, err := os.Stat(redisServerConfigFile); os.IsNotExist(err) {
		return redisServerConf
	}

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
		_ = json.NewEncoder(w).Encode(struct {
			OK        string
			DefaultDb int
			Dbs       int
		}{
			OK:        "OK",
			DefaultDb: foundServer.DefaultDb,
			Dbs:       dbs,
		})
	} else {
		_ = json.NewEncoder(w).Encode(struct {
			OK string
		}{
			OK: "Server Unknown",
		})
	}
}
