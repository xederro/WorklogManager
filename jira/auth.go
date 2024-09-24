package jira

import (
	"fmt"
	"net/http"
)

type auth interface {
	addAuth(req *http.Request)
}

type token string

func (t token) addAuth(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))
}
