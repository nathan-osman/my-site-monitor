package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nathan-osman/my-site-monitor/db"
)

const (
	contextUser = "user"

	sessionName   = "session"
	sessionUserID = "userID"
)

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

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, sessionName)
	v, _ := session.Values[sessionUserID]
	if v != nil {
		u := &db.User{}
		if err := s.conn.First(u, v).Error; err == nil {
			writeJSON(w, u, http.StatusOK)
			return
		}
	}
	writeError(w, http.StatusForbidden)
}
