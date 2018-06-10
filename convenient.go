package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"github.com/go-ini/ini"
	"net/http"
	"os"
	"strings"
)

type ConvenientItem struct {
	Section    string
	Name       string
	Template   string
	Operations []string
	Ttl        string
}

type ConvenientConfig struct {
	Ready bool
	Error string
	Items []ConvenientItem
}

func convenientConfigNew(name, template, operations, ttl string) (string, string) {
	err := createIniFileIfNotExists()
	if err != nil {
		return "", err.Error()
	}

	cfg, err := ini.Load(*convenientConfigFile)
	if err != nil {
		return "", err.Error()
	}

	sectionName := go_utils.RandString(10)
	section, err := cfg.NewSection(sectionName)
	if err != nil {
		return sectionName, err.Error()
	}

	section.NewKey("name", name)
	section.NewKey("template", template)
	section.NewKey("operations", operations)
	section.NewKey("ttl", ttl)

	err = cfg.SaveTo(*convenientConfigFile)
	if err != nil {
		return sectionName, err.Error()
	}

	return sectionName, "OK"
}

func deleteConvenientConfigItem(sectionName string) error {
	err := createIniFileIfNotExists()
	if err != nil {
		return err
	}

	cfg, err := ini.Load(*convenientConfigFile)
	if err != nil {
		return err
	}

	cfg.DeleteSection(sectionName)

	err = cfg.SaveTo(*convenientConfigFile)

	return err
}

func parseConvenientConfig() ConvenientConfig {
	items := make([]ConvenientItem, 0)
	err := createIniFileIfNotExists()
	if err != nil {
		return ConvenientConfig{
			Ready: false,
			Error: err.Error(),
			Items: items,
		}
	}

	cfg, err := ini.Load(*convenientConfigFile)
	if err != nil {
		return ConvenientConfig{
			Ready: false,
			Error: err.Error(),
			Items: items,
		}
	}

	sectionNames := cfg.SectionStrings()
	for _, sectionName := range sectionNames {
		if sectionName == "DEFAULT" {
			continue
		}

		section := cfg.Section(sectionName)
		nameKey, _ := section.GetKey("name")
		name := "Unnamed"
		if nameKey != nil {
			name = nameKey.MustString(name)
		}

		template, _ := section.GetKey("template")
		operationsKey, _ := section.GetKey("operations")
		operations := "save, delete"
		if operationsKey != nil {
			operations = operationsKey.MustString(operations)
		}

		ttlKey, _ := section.GetKey("ttl")
		ttl := "-1s"
		if ttlKey != nil {
			ttl = ttlKey.MustString("-1s")
		}

		items = append([]ConvenientItem{{
			Section:    sectionName,
			Name:       name,
			Template:   template.String(),
			Operations: strings.Split(operations, ","),
			Ttl:        ttl,
		}}, items...)
	}

	return ConvenientConfig{
		Ready: true,
		Error: "",
		Items: items,
	}
}
func createIniFileIfNotExists() error {
	file, err := os.OpenFile(*convenientConfigFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func serveConvenientConfigRead(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	config := parseConvenientConfig()
	json.NewEncoder(w).Encode(config)
}

func serveDeleteConvenientConfigItem(w http.ResponseWriter, req *http.Request) {
	sectionName := strings.TrimSpace(req.FormValue("sectionName"))
	deleteConvenientConfigItem(sectionName)
}

func serveConvenientConfigAdd(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	name := strings.TrimSpace(req.FormValue("name"))
	template := strings.TrimSpace(req.FormValue("template"))
	operations := strings.TrimSpace(req.FormValue("operations"))
	ttl := strings.TrimSpace(req.FormValue("ttl"))

	sectionName, result := convenientConfigNew(name, template, operations, ttl)
	json.NewEncoder(w).Encode(struct {
		Section string
		Message string
	}{sectionName, result})
}
