package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
)

type ProxyRoute struct {
	Name     string
	Target   string
	Restrict bool
}

type Config struct {
	ProxyRoutes    []ProxyRoute
	Port           uint
	Host           string
	CSRFProtection bool
}

func getConfig() *Config {
	raw, err := ioutil.ReadFile(path.Join(dirname, "config.json"))

	if err != nil {
		log.Fatal("FAILED TO READ CONFIGURATION FILE")
		log.Panic(err)
	}

	var j Config
	err = json.Unmarshal(raw, &j)

	if err != nil {
		log.Fatal("FAILED TO DESERIALIZE CONFIGURATION STRING")
		log.Panic(err)
	}

	return &j
}
