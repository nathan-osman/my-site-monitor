package server

import (
	"github.com/nathan-osman/api2go-resource"
	"github.com/nathan-osman/my-site-monitor/db"
)

func (s *Server) prepareSite(p *resource.Params) error {
	switch p.Action {
	case resource.FindAll:
		p.DB = p.DB.Order("name")
	case resource.Create:
		u, _ := p.Request.Context.Get(contextUser)
		p.Obj.(*db.Site).UserID = u.(*db.User).ID
	}
	return nil
}
