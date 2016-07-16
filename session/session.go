package session

import (
	"aista-search/config"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	Store *sessions.CookieStore
	Name  string
)

func Configure() {
	ssecret := config.GetEnv("AISTA_SEARCH_SESSION_SECRET", "episode solo")
	Store = sessions.NewCookieStore([]byte(ssecret))
	Name = "aista-search"
}

func Instance(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, Name)
	return session
}

func Empty(session *sessions.Session) {
	for k := range session.Values {
		delete(session.Values, k)
	}
}
