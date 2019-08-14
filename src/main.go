package main

import (
	"crypto/rand"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	dirname  string
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

//func test(w http.ResponseWriter, r *http.Request) {
//	_, err := fmt.Fprint(w, "hit me fam")
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
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
	if v, ex := os.LookupEnv("GO_MODE"); ex && v == "development" {
		dev = true
	} else {
		dev = false
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
	addr := ":" + strconv.Itoa(int((*config).HTTPSPort))

	if config.CSRFProtection {
		mux = getProtect()(mux)
	}

	mux = handlers.CombinedLoggingHandler(os.Stdout, mux)

	log.SetOutput(os.Stdout)
	//log.Println("starting server at " + config.HTTPSHost)

	err := http.ListenAndServeTLS(addr, config.TLSCertPath, config.TLSKeyPath, mux)

	if err != nil {
		c <- err
	}
}

type HTTPSRedirector struct {
	config *Config
}

func (h *HTTPSRedirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if dev {
		//println(r.URL.Path)
		//log.Fatal(r.URL.Path)
		http.Redirect(w, r, "https://localhost:8443"+r.URL.Path, 303)
	} else {
		http.Redirect(w, r, h.config.HTTPSHost+r.URL.Path, 303)
	}
}

func redirectToHTTPS(c chan error, config *Config) {
	addr := ":" + strconv.Itoa(int((*config).HTTPPort))

	handler := handlers.CombinedLoggingHandler(os.Stdout, &HTTPSRedirector{config: config})

	log.SetOutput(os.Stdout)
	//log.Println("Redirecting http from: " + config.HTTPHost)

	err := http.ListenAndServe(addr, handler)

	if err != nil {
		c <- err
	}
}

func main() {
	dirname = getDirname()
	loadEnv()
	config := getConfig()

	pmux := createProxyMux(config)

	c := make(chan error)
	go redirectToHTTPS(c, config)
	go listen(c, pmux, config)

	if dev {
		println("Development server at https://localhost:8443")
	}

	log.Fatal(<-c)

}
