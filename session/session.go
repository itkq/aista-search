package session

import (
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

var (
	Store *sessions.CookieStore
	Name  string
)

func Configure() {
	ssecret := os.Getenv("AISTA_SEARCH_SESSION_SECRET")
	if ssecret == "" {
		ssecret = "episode solo"
	}

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
