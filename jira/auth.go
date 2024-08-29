package jira

import (
	"net/http"
)

type auth interface {
	addToken(req *http.Request)
}

type basicAuth struct {
	Pass  string
	Login string
}

func (b *basicAuth) addToken(req *http.Request) {
	req.SetBasicAuth(b.Login, b.Pass)
}

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) addToken(req *http.Request) {
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: t.Token,
	})
}
