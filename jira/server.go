package jira

import (
	"errors"
	"fmt"
	"github.com/xederro/WorklogManager/config"
	"io"
	"net/http"
)

var (
	client        = &http.Client{}
	UrlBase       = config.BASEJIRALINK
	jiraAuth auth = nil
)

type Jira struct{}

func (j Jira) SetBasicAuth(username, password string) error {
	jiraAuth = &basicAuth{
		Pass:  password,
		Login: username,
	}

	r, err := j.Request("GET", fmt.Sprintf("%s/myself", UrlBase), nil)
	fmt.Println(string(r))
	if err != nil {
		return err
	}
	return nil
}

func (j Jira) SetTokenAuth(token string) error {
	jiraAuth = &tokenAuth{
		Token: token,
	}
	fmt.Println(1)
	_, err := j.Request("GET", fmt.Sprintf("%s/myself", UrlBase), nil)
	if err != nil {
		return err
	}
	return nil
}

func (j Jira) Request(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte{}, err
	}

	j.addHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []byte{}, errors.New(resp.Status)
	}

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return rBody, nil
}

func (j Jira) addHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	jiraAuth.addToken(req)
}
