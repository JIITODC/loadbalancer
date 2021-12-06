package main

import (
	"log"
	"net/http"
)


func main(){
	//server 1
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	
	if (r.Method == "HEAD"){
		print("Request from health check ")
		println("\n")
	}else if (r.Method == "POST"){
		log.Printf("server 2 says hello....MEHTOD = POST "  )
	}else if (r.Method == "GET"){
		log.Printf("server 2 says hello....MEHTOD = GET "  )
	}else if (r.Method == "PUT"){
		log.Printf("server 2 says hello....MEHTOD = PUT "  )
	}else {
		log.Printf("server 2 says hello " , r.Method )
	}
	})

	server := &http.Server{
		Addr:    "127.0.0.1:3002",
		Handler: myMux,
	}
	
	server.ListenAndServe()
	println("server 1 running on port 3001")
	
	println("\n")

}

