package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/BurntSushi/toml"
)

func serveSaveRedisServerConfig(w http.ResponseWriter, req *http.Request) {
	HeadContentTypeJson(w)
	redisServerConfig := req.FormValue("redisServerConfig")

	var redisServerConf RedisServerConf
	_, err := toml.Decode(redisServerConfig, &redisServerConf)
	if err != nil {
		_ = json.NewEncoder(w).Encode(struct {
			OK string
		}{
			OK: err.Error(),
		})
		return
	}

	err = ioutil.WriteFile(redisServerConfigFile, []byte(redisServerConfig), 0644)
	if err != nil {
		_ = json.NewEncoder(w).Encode(struct {
			OK string
		}{
			OK: err.Error(),
		})
		return
	}

	loadRedisServerConf()
	servers = parseServers(appConfig.Servers)

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

	_ = json.NewEncoder(w).Encode(struct {
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
