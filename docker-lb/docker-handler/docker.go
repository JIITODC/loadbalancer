package dockerhandler

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/Shreyas220/loadbalancer/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerHandler struct {
	LastServedIndex int
	Docker          docker
	Logger          log.Logger
	Label           string
}

type docker struct {
	DockerClient *client.Client
	ContainerIDs map[string]struct{}
	Servers      []*models.Server
}

func InitalizeDockerHandler() DockerHandler {

	DockerClient, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatalln(err)
	}
	Dock := &docker{}
	Dock.DockerClient = DockerClient

	return DockerHandler{
		LastServedIndex: 0,
		Docker:          *Dock,
		Logger:          *log.New(os.Stdout, "DOCKER-LB: ", 2),
	}
}

// takes a container ID as input and returns info about the service running
// in it
func (dh *DockerHandler) GetServiceInfo(containerID string) (*models.Server, error) {
	if dh.Docker.DockerClient == nil {
		return nil, errors.New("no docker client")
	}

	// inspect containers
	inspect, err := dh.Docker.DockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}

	// collect info on services
	server := &models.Server{}

	// labels tell us if a service has to be load balanced or not
	labels := inspect.Config.Labels

	// fetch label name
	if _, ok := labels[dh.Label]; !ok {
		return nil, nil
	}

	server.ServiceName = labels[dh.Label]
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
		if _, ok := dh.Docker.ContainerIDs[activeContainerID]; !ok {
			newContainers[activeContainerID] = struct{}{}
		}
	}

	return newContainers
}

func (dh *DockerHandler) GetDeletedDockerContainers(containers map[string]struct{}) map[string]struct{} {
	deletedContainers := make(map[string]struct{})
	for globalContainerID := range dh.Docker.ContainerIDs {
		if _, ok := containers[globalContainerID]; !ok {
			deletedContainers[globalContainerID] = struct{}{}
			delete(dh.Docker.ContainerIDs, globalContainerID)
		}
	}

	dh.Docker.ContainerIDs = containers

	return deletedContainers
}

// lists docker contairnes
func (dh *DockerHandler) ListDockerContainers() (map[string]struct{}, error) {
	containerIDs := make(map[string]struct{})
	if containerList, err := dh.Docker.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{}); err == nil {
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
