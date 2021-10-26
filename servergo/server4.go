package servergo

import (
	"log"
	"net/http"
)



func server4(){

	http.HandleFunc("/",printHello)
	log.Fatal(http.ListenAndServe(":8004", nil))

}

func printHello(res http.ResponseWriter, req *http.Request) {
	println("server 4 says hello")
}