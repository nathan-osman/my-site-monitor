package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/manyminds/api2go"
	"github.com/nathan-osman/api2go-resource"
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/sirupsen/logrus"
)

// Server provides the API for interacting with the database and serves the
// static files that comprise the UI.
type Server struct {
	listener net.Listener
	conn     *db.Conn
	store    *sessions.CookieStore
	log      *logrus.Entry
	stopped  chan bool
}

// New creates and initializes the server.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		r = mux.NewRouter()
		s = &Server{
			listener: l,
			conn:     cfg.Conn,
			store:    sessions.NewCookieStore([]byte(cfg.SecretKey)),
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
		api    = api2go.NewAPI("api")
		server = http.Server{
			Handler: r,
		}
	)
	api.AddResource(&db.User{}, &resource.Resource{
		DB:    s.conn.DB,
		Type:  &db.User{},
		Hooks: []resource.Hook{s.requireLogin},
	})
	api.AddResource(&db.Site{}, &resource.Resource{
		DB:     s.conn.DB,
		Type:   &db.Site{},
		Hooks:  []resource.Hook{s.requireLogin, s.prepareSite},
		Fields: []string{"id"},
	})
	api.AddResource(&db.Outage{}, &resource.Resource{
		DB:    s.conn.DB,
		Type:  &db.Outage{},
		Hooks: []resource.Hook{s.requireLogin},
	})
	r.HandleFunc("/api/login", s.login)
	r.HandleFunc("/api/logout", s.logout)
	r.PathPrefix("/api/").Handler(api.Handler())
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
