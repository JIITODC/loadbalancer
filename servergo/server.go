package main

import (
	"log"
	"net/http"
)


func main(){
	//server 1
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 1 says hello")})

	server := &http.Server{
		Addr:    "127.0.0.1:3001",
		Handler: myMux,
	}
	go server.ListenAndServe()

	//server 2
	myMux2 := http.NewServeMux()
	myMux2.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 2 says hello")})

	server2 := &http.Server{
		Addr:    "127.0.0.1:3002",
		Handler: myMux2,
	}
	go server2.ListenAndServe()

	//server 3
	myMux3 := http.NewServeMux()
	myMux3.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 3 says hello")})

	server3 := &http.Server{
		Addr:    "127.0.0.1:3003",
		Handler: myMux3,
	}
	go server3.ListenAndServe()

	//server4
	myMux4 := http.NewServeMux()
	myMux4.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 4 says hello")})

	server4 := &http.Server{
		Addr:    "127.0.0.1:3004",
		Handler: myMux4,
	}
	server4.ListenAndServe()




}

func server1() {
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 1 says hello")})

	server := &http.Server{
		Addr:    "127.0.0.1:3001",
		Handler: myMux,
	}
	server.ListenAndServe()

}

func server2() {
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 2 says hello")})

	server := &http.Server{
		Addr:    "127.0.0.1:3002",
		Handler: myMux,
	}
	server.ListenAndServe()

}
func server3() {
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {println("server 3 says hello")})

	server := &http.Server{
		Addr:    "127.0.0.1:3003",
		Handler: myMux,
	}
	server.ListenAndServe()

}
func server4() {
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {log.Println("server 4 says hello")})

	server := &http.Server{
		Addr:    "127.0.0.1:3004",
		Handler: myMux,
	}
	server.ListenAndServe()

}