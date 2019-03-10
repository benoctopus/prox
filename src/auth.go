package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/mediocregopher/radix"
	"log"
	"net/http"
	"time"
)

type SessionManager struct {
	Next http.Handler
	Mode int
	Pool *radix.Pool
}

func (s *SessionManager) AuthFailureRedirect(w http.ResponseWriter, r *http.Request) {
	var code int
	if r.Method == "GET" {
		code = 303
		http.Redirect(w, r, "/login", code)
		return
	}
	code = 403
	w.WriteHeader(code)
	return
}

func (s *SessionManager) makeToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return token, nil
}

func (s *SessionManager) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		return
	}
	// TODO: get body, validate, store

	token, terr := s.makeToken()

	// TODO: remove
	id, ierr := s.makeToken()

	if terr != nil {
		log.Fatal(terr)
		w.WriteHeader(501)
		return
	}

	// TODO: remove
	if ierr != nil {
		log.Fatal(ierr)
		w.WriteHeader(501)
		return
	}

	cookie := http.Cookie{
		Name:     "atx",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(1 * 24 * time.Hour),
	}

	serr := s.Pool.Do(radix.Cmd(nil, "Set", token, id ))

	if serr != nil {
		log.Fatal(serr)
		w.WriteHeader(501)
		return
	}

	http.SetCookie(w, &cookie)

	s.Next.ServeHTTP(w, r)
}

//func (s *SessionManager) Login(w http.ResponseWriter, r *http.Request) {
//	// Todo: make Post only Route, Validate Data, Check
//	r.ParseForm()
//	var token string
//	s.Pool.Do(radix.Cmd(nil, "set", token, Id))
//}

func (s *SessionManager) Auth(w http.ResponseWriter, r *http.Request) {
	if c, ece := r.Cookie("atx"); ece == nil {
		var id string
		rerr := s.Pool.Do(radix.Cmd(&id, "get", c.Value))
		r.Header.Set("User", id)

		if rerr != nil {
			log.Fatal(rerr)
			w.WriteHeader(501)
			return
		}
	} else {
		s.AuthFailureRedirect(w, r)
		return
	}

	s.Next.ServeHTTP(w, r)
}

func (s *SessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Mode == 0 {
		s.Auth(w, r)
	} else if s.Mode == 1 {
		s.Signup(w, r)
	}
}

func withSession(next http.Handler, mode int) http.Handler {
	sm := SessionManager{
		Next: next,
		Mode: mode,
		Pool: getRedis(),
	}

	return &sm
}

var getRedis = _getRedis()

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
