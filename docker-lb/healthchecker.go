package main

import (
	"time"
)

func StartHealthCheck() {
	for {
		for _, server := range Docker.Servers {
			name, healthy := server.CheckHealth()
			if healthy {
				Logger.Printf("%s is healthy!", name)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
