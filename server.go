package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type serve struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Health       bool
}

func newServer(name string, urlStr string) *serve {

	u, _ := url.Parse(urlStr)
	rp := httputil.NewSingleHostReverseProxy(u)
	return &serve{
		Name:         name,
		URL:          urlStr,
		ReverseProxy: rp,
		Health:       true,
	}

}

func (s *serve) checkHealth() (string, bool) {
	resp, err := http.Head(s.URL)

	if err != nil {
		s.Health = false
		return s.Name, s.Health
	}

	if resp.StatusCode != http.StatusOK {
		s.Health = false
		return s.Name, s.Health
	}

	s.Health = true

	return s.Name, s.Health
}
