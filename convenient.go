package main

import (
	"github.com/go-ini/ini"
	"strings"
	"net/http"
	"encoding/json"
)

type ConvenientItem struct {
	Section    string
	Name       string
	Template   string
	Operations []string
}

type ConvenientConfig struct {
	Ready bool
	Error string
	Items []ConvenientItem
}

func parseConvenientConfig() ConvenientConfig {
	items := make([]ConvenientItem, 0)

	cfg, err := ini.Load(convenientConfigFile)
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
		key, _ := section.GetKey("name")
		template, _ := section.GetKey("template")
		operations, _ := section.GetKey("operations")
		items = append(items, ConvenientItem{
			Section:    sectionName,
			Name:       key.MustString("Unnamed"),
			Template:   template.String(),
			Operations: strings.Split(operations.MustString("save,delete"), ","),
		})
	}

	return ConvenientConfig{
		Ready: true,
		Error: "",
		Items: items,
	}
}

func serveConvenientConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	config := parseConvenientConfig()
	json.NewEncoder(w).Encode(config)
}
