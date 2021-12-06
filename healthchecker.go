package main

import (
	"context"
	"log"
	"time"

	"github.com/procyon-projects/chrono"
)

func StartHealthCheck() {
	taskScheduler := chrono.NewDefaultTaskScheduler()
	log.Printf("before loop")
	
		_, err := taskScheduler.ScheduleAtFixedRate(func(ctx context.Context) {
			for _, server := range serverList {
			name,healthy := server.checkHealth()
			if healthy {
				log.Printf("🟢'%s' is healthy!",name )
			} else {
				log.Printf("🔴'%s' is unhealthy!",name)
			}
		}
		println("")

		}, 2 * time.Second)

			if err != nil {
				println("error \n")
			}
	

	
}


