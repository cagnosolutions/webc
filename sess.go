package web

import (
	"net/http"
	"sync"
	"time"
)

const SID = "GOSESSID"

type Store struct {
	rate     int
	sessions map[string]*Session
	sync.RWMutex
}

func StoreInstance(rate int) *Store {
	store := &Store{
		rate:     rate,
		sessions: make(map[string]*Session, 0),
	}
	store.GC()
	return store
}

func GetId(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(SID)
	if err != nil || cookie.Value == "" {
		return UUID4(), false
	}
	return cookie.Value, true
}

func (s *Store) InitSession(uuid string, w http.ResponseWriter) *Session {
	session := SessionInstance()
	s.Lock()
	s.sessions[uuid] = session
	s.Unlock()
	return session
}

func InitCookie(uuid, w) {
	cookie := http.Cookie{
		Name:     SID,
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Duration(s.rate) * time.Second),
		MaxAge:   s.rate,
	}
	http.SetCookie(w, &cookie)
}

func (s *Store) GetSession(w http.ResponseWriter, r *http.Request) *Session {
	uuid, ok := GetId(r)
	InitCookie(uuid, w)
	if ok {
		s.RLock()
		if session, ok := s.session[uuid]; ok {
			s.RUnlock()
			return session
		}
		s.RUnlock()
	}
	return InitSession(uuid, w)
}
