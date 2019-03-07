package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	Proxy *httputil.ReverseProxy
	URL   *url.URL
	Host  string
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = p.URL.Host
	r.URL.Scheme = p.URL.Scheme
	r.URL.Path = getRelativePath(r.URL.Path)
	r.Header.Set("X-Forward-Host", r.Header.Get("Host"))
	p.Proxy.ServeHTTP(w, r)
}

func getRelativePath(path string) string {
	count := 0
	var start int
	for i, c := range path {
		if c == '/' {
			count++
		}
		if count > 1 {
			start = i
			break
		}
	}
	return string(path[start:])
}

func createProxies(config *Config, mux *http.ServeMux) {
	for _, route := range config.ProxyRoutes {
		h := createProxy(&route, config.Host)
		mux.Handle(route.Name, h)
	}
}

func createProxy(cr *ProxyRoute, host string) http.Handler {
	t := cr.Target

	u, err := url.Parse(t)
	if err != nil {
		log.Fatal(err)
	}

	p := httputil.NewSingleHostReverseProxy(u)
	handler := ProxyHandler{Proxy: p, URL: u, Host: host}

	return &handler
}
