package main

import (
	"time"
)

func StartHealthCheck() {
	for {
		for _, server := range dh.Docker.Servers {
			name, healthy := server.CheckHealth()
			if healthy {
				dh.Logger.Printf("%s is healthy!", name)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
