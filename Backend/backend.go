package main

import (
	"log"
	"net/http"
)

func main() {
	//server 1
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { log.Printf("server 1 says hello ", r.Method) })

	server := &http.Server{
		Addr:    "127.0.0.1:3001",
		Handler: myMux,
	}

	go server.ListenAndServe()
	println("server 1 running on port 3001")
	//server 2
	myMux2 := http.NewServeMux()
	myMux2.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { log.Printf("server 2 says hello ", r.Method) })

	server2 := &http.Server{
		Addr:    "127.0.0.1:3002",
		Handler: myMux2,
	}
	go server2.ListenAndServe()
	println("server 2 running on port 3002")

	//server 3
	myMux3 := http.NewServeMux()
	myMux3.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { log.Printf("server 3 says hello ", r.Method) })

	server3 := &http.Server{
		Addr:    "127.0.0.1:3003",
		Handler: myMux3,
	}
	go server3.ListenAndServe()
	println("server 3 running on port 3003")

	//server4
	myMux4 := http.NewServeMux()
	myMux4.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { log.Printf("server 4 says hello ", r.Method) })

	server4 := &http.Server{
		Addr:    "127.0.0.1:3004",
		Handler: myMux4,
	}
	println("server 4 running on port 3004")
	server4.ListenAndServe()

	println("")

}
