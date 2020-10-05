package main

import (
	"strconv"
	"strings"
)

type RedisServer struct {
	ServerName string
	Addr       string
	Password   string
	DefaultDb  int
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
