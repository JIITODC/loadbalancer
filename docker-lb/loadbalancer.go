package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shreyas220/loadbalancer/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

//server pool
var (
	lastServedIndex = 0
	Docker          *DockerHandler
	Logger          *log.Logger
	SigChan         chan os.Signal
)

const (
	DockerSocket = "unix:///var/run/docker.sock"
)

type DockerHandler struct {
	DockerClient *client.Client
	ContainerIDs map[string]struct{}
	Servers      []*models.Server
}

// init Function
func init() {
	Docker = &DockerHandler{}

	DockerClient, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatalln(err)
	}

	Docker.DockerClient = DockerClient
	SigChan = GetOSSigChannel()
	Logger = log.New(os.Stdout, "DOCKER-LB: ", 2)
}

// 1. Get all container IDs in a map[string]struct{}
// 1.a If new container, get their service's info in a server struct
// 1.b If container died (healthcheck), don't send requests to it
func main() {

	// check if DockerHandler exists
	if Docker == nil {
		return
	}

	Logger.Println("Started to monitor docker services")
	go StartHealthCheck()

	// stop if termination signal received
	go func() {
		for {
			select {
			case <-SigChan:
				Logger.Println("Stopping load balancer")
				os.Exit(0)

			default:
				containers, err := Docker.ListDockerContainers()
				if err != nil {
					return
				}

				// if number of stored container IDs is equal to number of container IDs
				// returned by the API, no containers added/deleted
				if len(containers) == len(Docker.ContainerIDs) {
					time.Sleep(time.Millisecond * 10)
					continue
				}

				newContainers := Docker.GetNewDockerContainers(containers)
				deletedContainers := Docker.GetDeletedDockerContainers(containers)

				if len(newContainers) > 0 {
					for containerID := range newContainers {
						serviceInfo, err := Docker.GetServiceInfo(containerID)
						if err != nil {
							log.Println("WARNING")
						}
						Docker.Servers = append(Docker.Servers, serviceInfo)
					}
				}

				if len(deletedContainers) > 0 {
					for containerID := range deletedContainers {
						delete(Docker.ContainerIDs, containerID)
						for _, server := range Docker.Servers {
							if server.ContainerID == containerID {
								Logger.Printf("STOPPED MONTIORING %s AS IT IS UNHEALTHY!", server.Name)
								server = nil
							}
						}
					}
				}
			}

			time.Sleep(time.Millisecond * 50)
		}
	}()

	http.HandleFunc("/", ForwardRequest)
	Logger.Fatalln(http.ListenAndServe(":8000", nil))
}

// takes a container ID as input and returns info about the service running
// in it
func (dh *DockerHandler) GetServiceInfo(containerID string) (*models.Server, error) {
	if dh.DockerClient == nil {
		return &models.Server{}, errors.New("no docker client")
	}

	// inspect containers
	inspect, err := dh.DockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return &models.Server{}, err
	}

	// collect info on services
	server := &models.Server{}

	// labels tell us if a service has to be load balanced or not
	labels := inspect.Config.Labels
	loadBalance := labels["com.docker-lb.load-balance"]
	if loadBalance != "true" {
		return &models.Server{}, nil
	}

	server.ServiceName = labels["com.docker-lb.service-name"]

	server.Name = inspect.Name

	networkName := inspect.HostConfig.NetworkMode.NetworkName()
	endpointSettings := inspect.NetworkSettings.Networks[networkName]

	server.URL = endpointSettings.IPAddress
	server.ContainerID = inspect.ID

	server = models.NewServer(server.Name, server.URL)

	return server, nil
}

func (dh *DockerHandler) GetNewDockerContainers(containers map[string]struct{}) map[string]struct{} {
	newContainers := make(map[string]struct{})
	for activeContainerID := range containers {
		if _, ok := dh.ContainerIDs[activeContainerID]; !ok {
			newContainers[activeContainerID] = struct{}{}
		}
	}

	return newContainers
}

func (dh *DockerHandler) GetDeletedDockerContainers(containers map[string]struct{}) map[string]struct{} {
	deletedContainers := make(map[string]struct{})
	for globalContainerID := range dh.ContainerIDs {
		if _, ok := containers[globalContainerID]; !ok {
			deletedContainers[globalContainerID] = struct{}{}
			delete(dh.ContainerIDs, globalContainerID)
		}
	}

	dh.ContainerIDs = containers

	return deletedContainers
}

// lists docker contairnes
func (dh *DockerHandler) ListDockerContainers() (map[string]struct{}, error) {
	containerIDs := make(map[string]struct{})
	if containerList, err := Docker.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{}); err == nil {
		for _, dcontainer := range containerList {
			if dcontainer.ID == "" {
				continue
			}
			containerIDs[dcontainer.ID] = struct{}{}
		}
	} else {
		return nil, err
	}
	return containerIDs, nil
}

func ForwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := GetHealthyServer()
	if err != nil {
		fmt.Fprintf(res, "Couldnt process request: %s", err.Error())
	}
	server.ReverseProxy.ServeHTTP(res, req)
}

func GetHealthyServer() (*models.Server, error) {
	for i := 0; i < len(Docker.ContainerIDs); i++ {
		server := GetNextServer()
		if server.Healthy {
			return server, nil
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}

func GetNextServer() *models.Server {
	nextIndex := (lastServedIndex + 1) % len(Docker.ContainerIDs)
	server := Docker.Servers[nextIndex]
	lastServedIndex = nextIndex
	return server
}

func GetOSSigChannel() chan os.Signal {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	return c
}
