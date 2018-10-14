package server

import (
	"net/http"

	"github.com/manyminds/api2go"
	"github.com/nathan-osman/api2go-resource"
	"github.com/nathan-osman/my-site-monitor/db"
)

func (s *Server) restricted(p *resource.Params) error {
	session, _ := s.store.Get(p.Request.PlainRequest, sessionName)
	v, _ := session.Values[sessionUserID]
	if v != nil {
		u := &db.User{}
		if err := s.conn.First(u, v).Error; err == nil {
			p.Request.Context.Set(
				contextUser,
				u,
			)
			return nil
		}
	}
	return api2go.NewHTTPError(nil, "login required", http.StatusForbidden)
}

func (s *Server) readOnly(p *resource.Params) error {
	switch p.Action {
	case resource.BeforeCreate, resource.BeforeDelete, resource.BeforeUpdate:
		return s.restricted(p)
	}
	return nil
}

func (s *Server) siteHook(p *resource.Params) error {
	switch p.Action {
	case resource.BeforeFindAll:
		p.DB = p.DB.Order("name")
	case resource.BeforeCreate:
		u, _ := p.Request.Context.Get(contextUser)
		p.Obj.(*db.Site).UserID = u.(*db.User).ID
	case resource.AfterCreate, resource.AfterUpdate:
		s.monitor.Trigger()
	}
	return nil
}
