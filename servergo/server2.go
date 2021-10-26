package servergo

import (
	"log"
	"net/http"
)



func server2(){

	http.HandleFunc("/",printHello)
	log.Fatal(http.ListenAndServe(":8002", nil))

}

func printHello(res http.ResponseWriter, req *http.Request) {
	println("server 2 says hello")
}