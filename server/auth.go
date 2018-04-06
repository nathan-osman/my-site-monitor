package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nathan-osman/my-site-monitor/db"
)

type key int

const contextUser key = iota

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userAuth struct {
	Conn *db.Conn
}

func (a *userAuth) Authenticate(r *http.Request) (interface{}, interface{}, error) {
	l := &loginParams{}
	if err := json.NewDecoder(r.Body).Decode(l); err != nil {
		return nil, nil, err
	}
	u := &db.User{}
	if err := a.Conn.Where("username", l.Username).First(&u).Error; err != nil {
		return nil, nil, nil
	}
	if err := u.Authenticate(l.Password); err != nil {
		return nil, nil, nil
	}
	return u.ID, u, nil
}

func (a *userAuth) Initialize(r *http.Request, i interface{}) (*http.Request, error) {
	u := &db.User{}
	if err := a.Conn.First(u, i).Error; err != nil {
		return nil, err
	}
	return r.WithContext(
		context.WithValue(r.Context(), contextUser, u),
	), nil
}
