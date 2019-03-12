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
	TLSCertPath    string
	TLSKeyPath     string
	HTTPSPort      uint
	HTTPSHost      string
	HTTPPort       uint
	HTTPHost       string
	CSRFProtection bool
	RedisURL       string
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

	if j.TLSCertPath == "" {
		j.TLSCertPath = path.Join(dirname, "../", "localhost", "cert.pem")
	} else if !path.IsAbs(j.TLSCertPath) {
		j.TLSCertPath = path.Join(dirname, j.TLSCertPath)
	}

	if j.TLSKeyPath == "" {
		j.TLSKeyPath = path.Join(dirname, "../", "localhost", "key.pem")
	} else if !path.IsAbs(j.TLSKeyPath) {
		j.TLSKeyPath = path.Join(dirname, j.TLSKeyPath)
	}

	if j.RedisURL == "" {
		j.RedisURL = "127.0.0.1:6379"
	}

	getRedis = _getRedis(&j)

	return &j
}
