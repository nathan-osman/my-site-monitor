package server

import (
	"github.com/nathan-osman/api2go-resource"
	"github.com/nathan-osman/my-site-monitor/db"
)

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
