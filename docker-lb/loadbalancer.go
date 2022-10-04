package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	dockerhandler "github.com/Shreyas220/loadbalancer/docker-lb/docker-handler"
	"github.com/Shreyas220/loadbalancer/docker-lb/utils"
	"github.com/Shreyas220/loadbalancer/models"
)

// 1. Get all container IDs in a map[string]struct{}
// 1.a If new container, get their service's info in a server struct
// 1.b If container died (healthcheck), don't send requests to it

var (
	dh      = dockerhandler.InitalizeDockerHandler()
	SigChan chan os.Signal
)

func main() {

	// check if DockerHandler exists
	dh.Logger.Println("Started to monitor docker services")
	go StartHealthCheck()

	// stop if termination signal received
	go loadbalancer()

	http.HandleFunc("/", ForwardRequest)
	dh.Logger.Fatalln(http.ListenAndServe(":8000", nil))
}

func loadbalancer() {
	SigChan := utils.GetOSSigChannel()
	for {
		select {
		case <-SigChan:
			dh.Logger.Println("Stopping load balancer")
			os.Exit(0)

		default:
			containers, err := dh.ListDockerContainers()
			if err != nil {
				return
			}

			// if number of stored container IDs is equal to number of container IDs
			// returned by the API, no containers added/deleted
			if len(containers) == len(dh.Docker.ContainerIDs) {
				time.Sleep(time.Millisecond * 10)
				continue
			}

			newContainers := dh.GetNewDockerContainers(containers)
			deletedContainers := dh.GetDeletedDockerContainers(containers)

			if len(newContainers) > 0 {
				for containerID := range newContainers {
					serviceInfo, err := dh.GetServiceInfo(containerID)
					if err != nil {
						log.Println("WARNING", err)
					}
					dh.Docker.Servers = append(dh.Docker.Servers, serviceInfo)
				}
			}

			if len(deletedContainers) > 0 {
				for containerID := range deletedContainers {
					delete(dh.Docker.ContainerIDs, containerID)
					for _, server := range dh.Docker.Servers {
						if server.ContainerID == containerID {
							dh.Logger.Printf("STOPPED MONTIORING %s AS IT IS UNHEALTHY!", server.Name)
							server = nil
						}
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 50)
	}

}

func ForwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := GetHealthyServer()
	if err != nil {
		fmt.Fprintf(res, "Couldnt process request: %s", err.Error())
	}
	server.ReverseProxy.ServeHTTP(res, req)
}

func GetHealthyServer() (*models.Server, error) {
	for i := 0; i < len(dh.Docker.ContainerIDs); i++ {
		server := GetNextServer()
		if server.Healthy {
			return server, nil
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}

func GetNextServer() *models.Server {
	nextIndex := (dh.LastServedIndex + 1) % len(dh.Docker.ContainerIDs)
	server := dh.Docker.Servers[nextIndex]
	dh.LastServedIndex = nextIndex
	return server
}
