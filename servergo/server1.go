package servergo

import (
	"log"
	"net/http"
)



func server1(){

	http.HandleFunc("/",printHello)
	log.Fatal(http.ListenAndServe(":8001", nil))

}

func printHello(res http.ResponseWriter, req *http.Request) {
	println("server 1 says hello")
}