package servergo

import (
	"log"
	"net/http"
)



func server3(){

	http.HandleFunc("/",printHello)
	log.Fatal(http.ListenAndServe(":8003", nil))

}

func printHello(res http.ResponseWriter, req *http.Request) {
	println("server 3 says hello")
}