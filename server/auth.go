package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/manyminds/api2go"
	"github.com/nathan-osman/api2go-resource"
	"github.com/nathan-osman/my-site-monitor/db"
)

type key int

const (
	contextUser key = iota

	sessionName   = "session"
	sessionUserID = "userID"
)

func (s *Server) requireLogin(p *resource.Params) error {
	switch p.Action {
	case resource.Create, resource.Delete, resource.Update:
		session, _ := s.store.Get(p.Request.PlainRequest, sessionName)
		v, _ := session.Values[sessionUserID]
		if v != nil {
			u := &db.User{}
			if err := s.conn.First(u, v).Error; err == nil {
				p.Request.PlainRequest = p.Request.PlainRequest.WithContext(
					context.WithValue(
						p.Request.PlainRequest.Context(),
						contextUser,
						u,
					),
				)
				return nil
			}
		}
		return api2go.NewHTTPError(nil, "login required", http.StatusBadRequest)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func writeError(w http.ResponseWriter, statusCode int) {
	writeJSON(
		w,
		&struct {
			Error string `json:"error"`
		}{
			Error: http.StatusText(statusCode),
		},
		statusCode,
	)
}

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}
	p := &loginParams{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}
	u := &db.User{}
	if err := s.conn.
		First(u, "username = ?", p.Username).
		Error; err != nil || u.Authenticate(p.Password) != nil {
		writeError(w, http.StatusForbidden)
		return
	}
	session, _ := s.store.Get(r, sessionName)
	session.Values[sessionUserID] = u.ID
	session.Save(r, w)
	writeJSON(w, u, http.StatusOK)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, sessionName)
	delete(session.Values, sessionUserID)
	session.Save(r, w)
	writeJSON(w, &struct{}{}, http.StatusOK)
}
