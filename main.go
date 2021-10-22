package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	serverList = []string{
		"",
		"",
	}
	lastServedIndex = 0
)

func main(){

	http.HandleFunc("/",forwardRequest)
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
	url,_ := url.Parse(serverList[lastServedIndex])
	lastServedIndex = nextIndex
	return url
}
/*


*/