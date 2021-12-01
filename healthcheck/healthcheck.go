package healthcheck

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/procyon-projects/chrono"
)

type Server struct {
	URL string 
	ReverseProxy *httputil.ReverseProxy
	Health bool
} 

func NewServer(urlStr string) *Server {

	u,_ := url.Parse(urlStr)
	rp := httputil.NewSingleHostReverseProxy(u)
	return &Server{
		URL:	urlStr,
		ReverseProxy: rp,
		Health: true,
	}

}

func (s *Server) checkHealth() {
	resp, err := http.Head(s.URL)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		s.Health = false
		return
	}
	s.Health = true
}

func (s *Server) startHealthCheck() {
	taskScheduler := chrono.NewDefaultTaskScheduler()

	_, err := taskScheduler.ScheduleAtFixedRate(func(ctx context.Context) {
		
		s.checkHealth()
	}, 5 * time.Second)

	if err == nil {
		log.Print("Task has been scheduled successfully.")
	}
}


