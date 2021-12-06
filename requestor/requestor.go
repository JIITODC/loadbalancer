package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {

	for {
			payload := fmt.Sprintf(`{ "seconds": %v }`, randonDuration())
			body := strings.NewReader(payload)

			log.Printf("sending request with payload: %v", payload)
			res, err := http.Post("http://localhost:8000/", "application/json", body)

			badReq := res != nil && res.StatusCode >= 400
			if err != nil || badReq {
				log.Printf("requests failed, res: %v err: %v", res, err)
			}
		

		randomPause()
	}
}

func randomPause() {
	randDuration := time.Duration(rand.Intn(1)) * time.Second
	time.Sleep(randDuration + 250*time.Millisecond)
}

func randonDuration() int {
	return rand.Intn(5) + 3
}