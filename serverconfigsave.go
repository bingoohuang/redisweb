package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net/http"
)

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
