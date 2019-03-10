package main

import (
	"crypto/rand"
	"fmt"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

var (
	dirname  string
	cert     string
	key      string
	dev      bool
	redisURL string
)

func getDirname() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getSecret() []byte {
	// Todo: make real key for production
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		log.Panic(err)
	}

	return secret
}

//func test(w http.ResponseWriter, r *http.Request) {
//	file, _ := ioutil.ReadFile(path.Join(dirname, "index.html"))
//	_, err := fmt.Fprintf(w, string(file), csrf.Token(r))
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func test(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "hit me fam")
	if err != nil {
		log.Fatal(err)
	}
}

type TestHandler struct {
	serve func(w http.ResponseWriter, r *http.Request)
}

func (t *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.serve(w, r)
}

func getTestHandler(fn func(w http.ResponseWriter, r *http.Request)) *TestHandler {
	return &TestHandler{serve: fn}
}

//func testpost(w http.ResponseWriter, r *http.Request) {
//	err := csrf.FailureReason(r)
//	if err != nil {
//		log.Fatal(err)
//	}
//	re := struct {
//		res string
//	}{res: "hello"}
//	res, _ := json.Marshal(re)
//	fmt.Fprint(w, res)
//}

func loadEnv() {
	if v, ex := os.LookupEnv("GO_MODE"); ex && v == "production" {
		dev = false
	} else {
		dev = true
	}
	if v, ex := os.LookupEnv("CERT_PATH"); ex {
		cert = v
	} else {
		cert = path.Join(dirname, "../", "localhost", "cert.pem")
	}
	if v, ex := os.LookupEnv("KEY_PATH"); ex {
		key = v
	} else {
		key = path.Join(dirname, "../", "localhost", "key.pem")
	}
	if v, ex := os.LookupEnv("REDIS_URL"); ex {
		redisURL = v
	} else {
		redisURL = "127.0.0.1:6379"
	}
}

func getProtect() func(http.Handler) http.Handler {
	secret := getSecret()
	return csrf.Protect(
		secret,
		csrf.Secure(!dev),
		csrf.RequestHeader("_csrf"),
	)
}

func listen(c chan error, mux http.Handler, config *Config) {
	addr := ":" + strconv.Itoa(int((*config).Port))
	log.SetOutput(os.Stdout)
	log.Println("starting server at " + config.Host)
	var err error

	if config.CSRFProtection {
		protect := getProtect()
		err = http.ListenAndServeTLS(addr, cert, key, protect(mux))
	} else {
		err = http.ListenAndServeTLS(addr, cert, key, mux)
	}
	if err != nil {
		c <- err
	}
}

func main() {
	dirname = getDirname()
	loadEnv()
	config := getConfig()

	pmux := createProxyMux(config)
	tmux := NewPipeableMux()

	pmux.Handle("/signup", withSession(getTestHandler(test), 1))
	pmux.Handle("/verify", withSession(getTestHandler(test), 0))

	mux := ChainPipeable([]*PipableMux{pmux, tmux})

	c := make(chan error)
	go listen(c, mux, config)
	log.Fatal(<-c)

}
