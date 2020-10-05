package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type ListKeysResult struct {
	Keys   []KeysResult
	Cursor uint64
}

func serveListKeys(w http.ResponseWriter, req *http.Request) {
	HeadContentTypeJson(w)
	matchPattern := strings.TrimSpace(req.FormValue("pattern"))
	cursorStr := strings.TrimSpace(req.FormValue("cursor"))
	var cursor uint64
	var err error
	if cursorStr == "" {
		cursor = 0
	} else {
		cursor, err = strconv.ParseUint(cursorStr, 10, 64)
		if err != nil {
			http.Error(w, "bad cursor parameter", 411)
			return
		}
	}

	server := findRedisServer(req)

	keys, ncursor, _ := listKeys(server, cursor, matchPattern, appConfig.MaxKeys)
	sort.Slice(keys[:], func(i, j int) bool {
		return keys[i].Key < keys[j].Key
	})
	_ = json.NewEncoder(w).Encode(ListKeysResult{
		Keys:   keys,
		Cursor: ncursor,
	})
}
