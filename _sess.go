package web

import (
	"net/http"
	"sync"
	"time"
)

const SID = "GOSESSID"

type ContextStore struct {
	rate     int
	contexts map[string]*Context
	sync.RWMutex
}

func ContextStoreInstance(rate int) *ContextStore {
	ctxStore := &ContextStore{
		rate:     rate,
		contexts: make(map[string]*Context, 0),
	}
	ctxStore.GC()
	return ctxStore
}

func (s *ContextStore) GetUnique(uuid string) *Context {
	context := ContextInstance()
	s.Lock()
	s.context[uuid] = context
	s.Unlock()
	return context
}

func (s *Store) GetContext(w http.ResponseWriter, r *http.Request) *Context {
	uuid, ok := GetId(r)  // try to get existing user id (via cookie)
	UpdateCookie(uuid, w) // create/update and set cookie
	if ok {               // if existing id...
		s.RLock()
		if context, ok := s.context[uuid]; ok { // update context...
			context.ts = time.Now()
			s.RUnlock()
			return context // return context
		}
		s.RUnlock()
	} // there was no existing id found...
	return s.GetContext(uuid) // create/update and return context
}

func (s *Store) GC() {
	self.Lock()
	defer self.Unlock()
	for uuid, context := range s.contexts {
		if (context.ts.Unix() + s.rate) < time.Now().Unix() {
			delete(s.contexts, uuid)
		} else {
			break
		}
	}
	time.AfterFunc(time.Duration(s.rate)*time.Second, func() {
		self.GC()
	})
}

func GetId(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(SID)
	if err != nil && err == http.ErrNoCookie || cookie.Value == "" {
		return UUID4(), false
	}
	return cookie.Value, true
}

func UpdateCookie(uuid, w) {
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
