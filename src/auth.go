package main

import (
	"github.com/mediocregopher/radix"
	"log"
	"net/http"
)

type SessionManager struct {
	Next http.Handler
	Mode string
	Pool *radix.Pool
}

func (s *SessionManager) AuthFailureRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/Login", 403)
}

func (s *SessionManager) Auth(w http.ResponseWriter, r *http.Request) {
	//if c, ece := r.Cookie("atx"); ece == nil {
	//
	//} else {
	//	s.AuthFailureRedirect(w, r)
	//}
	//s.Pool.Do(radix.Cmd())
}

func (s *SessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Mode == "auth" {
		s.Auth(w, r)
	} else if s.Mode == "login" {

	}
}

func withSession(next http.Handler, mode string) http.Handler {
	sm := SessionManager{
		Next: next,
		Mode: mode,
		Pool: getRedis(),
	}

	return &sm
}

var getRedis func() *radix.Pool = _getRedis()

func _getRedis() func() *radix.Pool {
	init := false
	var c *radix.Pool
	return func() *radix.Pool {
		if init {
			return c
		}
		var err error = nil
		c, err = radix.NewPool("tcp", redisURL, 10)

		if err != nil {
			log.Panic(err)
		}

		return c
	}
}
