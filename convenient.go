package main

import (
	"encoding/json"
	"github.com/go-ini/ini"
	"net/http"
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

		items = append(items, ConvenientItem{
			Section:    sectionName,
			Name:       name,
			Template:   template.String(),
			Operations: strings.Split(operations, ","),
			Ttl:        ttl,
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
