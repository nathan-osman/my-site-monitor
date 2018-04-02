package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server provides the API for interacting with the database and serves the
// static files that comprise the UI.
type Server struct {
	listener net.Listener
	log      *logrus.Entry
	stopped  chan bool
}

// New creates and initializes the server
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		r = mux.NewRouter()
		s = &Server{
			listener: l,
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
		server = http.Server{
			Handler: r,
		}
	)
	r.PathPrefix("/").Handler(http.FileServer(HTTP))
	go func() {
		defer close(s.stopped)
		defer s.log.Info("server stopped")
		s.log.Info("server started")
		if err := server.Serve(l); err != nil {
			s.log.Error(err)
		}
	}()
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stopped
}
