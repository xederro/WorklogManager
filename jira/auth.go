package jira

import (
	"net/http"
)

type auth interface {
	addAuth(req *http.Request)
}

type basicAuth struct {
	Pass  string
	Login string
}

func (b *basicAuth) addAuth(req *http.Request) {
	req.SetBasicAuth(b.Login, b.Pass)
}

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) addAuth(req *http.Request) {
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: t.Token,
	})
}
