package session

import (
	"net/http"
)

type Manager struct {
	store  map[string]bool
	name   string
	maxAge int
}

func NewManager(maxAge int) *Manager {
	return &Manager{
		store:  make(map[string]bool),
		name:   "MimaSessionID",
		maxAge: maxAge,
	}
}

func (manager *Manager) NewSession(sid string) http.Cookie {
	return http.Cookie{
		Name:     manager.name,
		Value:    sid,
		Path:     "/", // important
		MaxAge:   manager.maxAge,
		HttpOnly: true,
	}
}

func (manager *Manager) Add(w http.ResponseWriter, sid string) {
	session := manager.NewSession(sid)
	http.SetCookie(w, &session)
	manager.store[sid] = true
}

func (manager *Manager) Check(r *http.Request) bool {
	cookie, err := r.Cookie(manager.name)
	if err != nil || cookie.Value == "" || !manager.store[cookie.Value] {
		return false
	}
	return true
}

func (manager *Manager) DeleteSID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.name)
	if err == nil {
		manager.store[cookie.Value] = false
	}
	session := http.Cookie{
		Name:     manager.name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &session)
}
