package balancer

import (
	"fmt"
	"net/http"

	"github.com/Shreyas220/loadbalancer/healthcheck"
)

//server pool
var (
	serverList = []*healthcheck.Server{
		healthcheck.NewServer("http://127.0.0.1:3001"),
		healthcheck.NewServer("http://127.0.0.1:3002"),
		healthcheck.NewServer("http://127.0.0.1:3003"),
		healthcheck.NewServer("http://127.0.0.1:3004"),
	}
	lastServedIndex = 0
)

func ForwardRequest(res http.ResponseWriter, req *http.Request){
	server , err := getHealthyServer()
	if err != nil {
		fmt.Fprintf(res,"Couldnt process request: %s", err.Error())
	}
	server.ReverseProxy.ServeHTTP(res,req)
}

func getHealthyServer() (*healthcheck.Server, error) {
	for i:=0; i< len(serverList);i++ {
		server:= getServer()
		if server.Health {
			return server,nil
		}
	}
	return nil , fmt.Errorf("no healthy hosts")
}

func getServer() *healthcheck.Server {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastServedIndex = nextIndex

	return server
}

