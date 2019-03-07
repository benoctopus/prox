package main

import (
	"crypto/rand"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

var dirname string
var cert string
var key string
var dev bool

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
	size, err := rand.Read(secret)
	if err != nil {
		log.Panic(err)
	}

	log.Println(size)
	return secret
}

//func test(w http.ResponseWriter, r *http.Request) {
//	file, _ := ioutil.ReadFile(path.Join(dirname, "index.html"))
//	_, err := fmt.Fprintf(w, string(file), csrf.Token(r))
//	if err != nil {
//		log.Fatal(err)
//	}
//}

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
}

func main() {
	dirname = getDirname()
	loadEnv()
	config := getConfig()
	secret := getSecret()
	addr := ":" + strconv.Itoa(int((*config).Port))

	protect := csrf.Protect(
		secret,
		csrf.Secure(!dev),
		csrf.RequestHeader("_csrf"),
	)

	var mux *http.ServeMux
	mux = &http.ServeMux{}

	//mux.HandleFunc("/test", test)
	//mux.HandleFunc("/test/post", testpost)

	createProxies(config, mux)

	log.Println("starting server at " + config.Host)
	err := http.ListenAndServeTLS(addr, cert, key, protect(mux))
	if err != nil {
		log.Panic(err)
	}
}
