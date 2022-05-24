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
			} else {
				Logger.Printf("%s is unhealthy!", name)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
