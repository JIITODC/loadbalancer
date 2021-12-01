package main

import (
	"fmt"
	"log"
	"net/http"
)

//server pool
var (
	serverList  []*serve
	lastServedIndex = 0
)

func main() {
	serverList = []*serve { 
		newServer("Server-1","http://127.0.0.1:3001"),
		newServer("Server-2","http://127.0.0.1:3002"),
		newServer("Server-3","http://127.0.0.1:3003"),
		newServer("Server-4","http://127.0.0.1:3004"),
	}
	http.HandleFunc("/",ForwardRequest)
	println("length of serverlist is ",len(serverList))
	go StartHealthCheck() 
	log.Fatal(http.ListenAndServe(":8000",nil))
}

func ForwardRequest(res http.ResponseWriter, req *http.Request){
	server , err := getHealthyServer()
	if err != nil {
		fmt.Fprintf(res,"Couldnt process request: %s", err.Error())
	}
	server.ReverseProxy.ServeHTTP(res,req)
}

func getHealthyServer() (*serve, error) {
	for i:=0; i< len(serverList);i++ {
		server:= getServer()
		if server.Health {
			return server,nil
		}
	}
	return nil , fmt.Errorf("no healthy hosts")
}

func getServer() *serve {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastServedIndex = nextIndex
	return server
}

