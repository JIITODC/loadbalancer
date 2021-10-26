package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	serverList = []string{
		"http://127.0.0.1:3001",
		"http://127.0.0.1:3002",
		"http://127.0.0.1:3003",
		"http://127.0.0.1:3004",
	}
	lastServedIndex = 1
)

func main(){

	http.HandleFunc("/",forwardRequest)
	println("loadbalancer started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func forwardRequest(res http.ResponseWriter, req *http.Request){
	url := getServer()
	//reverse proxy 
	rProxy := httputil.NewSingleHostReverseProxy(url)
	rProxy.ServeHTTP(res,req)
}

func getServer() *url.URL {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	url,err := url.Parse(serverList[lastServedIndex])
	if err != nil {
		log.Fatalf("failed to parse: %s", err)
		os.Exit(1)
	}
	lastServedIndex = nextIndex
	println("requestion sent to ",url)
	return url
}
/*


*/