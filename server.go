package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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
