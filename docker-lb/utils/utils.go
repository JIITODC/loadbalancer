package utils

import (
	"os"
	"os/signal"
	"syscall"
)

//server pool

const (
	DockerSocket = "unix:///var/run/docker.sock"
)

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
