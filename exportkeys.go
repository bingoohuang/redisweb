package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

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
