package models

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Server struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Healthy      bool
	ServiceName  string
	ContainerID  string
}

var (
	Client = http.Client{
		Timeout: 1 * time.Second,
	}
)

func NewServer(name string, urlStr string) *Server {
	u, _ := url.Parse(fmt.Sprintf("http://%s", urlStr))
	rp := httputil.NewSingleHostReverseProxy(u)
	return &Server{
		Name:         name,
		URL:          u.String(),
		ReverseProxy: rp,
		Healthy:      true,
	}
}

func (s *Server) CheckHealth() (string, bool) {
	resp, err := Client.Get(fmt.Sprintf("%s/health", s.URL))

	if err != nil {
		s.Healthy = false
		return s.Name, s.Healthy
	}

	if resp.StatusCode != http.StatusOK {
		s.Healthy = false
		return s.Name, s.Healthy
	}

	s.Healthy = true

	return s.Name, s.Healthy
}
