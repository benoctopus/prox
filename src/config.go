package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ProxyRoute struct {
	Name     string
	Target   string
	Restrict bool
}

type Config struct {
	ProxyRoutes    []ProxyRoute
	HTTPSPort      uint
	HTTPSHost      string
	HTTPPort       uint
	HTTPHost       string
	CSRFProtection bool
}

func getConfig() *Config {
	var p string

	if v, ex := os.LookupEnv("CONFIG_PATH"); ex {
		if path.IsAbs(v) {
			p = v
		} else {
			p = path.Join(dirname, v)
		}
	} else {
		p = path.Join(dirname, "config.json")
	}

	raw, err := ioutil.ReadFile(p)

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
