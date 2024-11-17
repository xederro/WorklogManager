package jira

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	client        = &http.Client{}
	UrlBase       = ""
	jiraAuth auth = nil
)

type Jira struct{}

func (j *Jira) SetUrlBase(urlBase string) {
	UrlBase = urlBase
}

func (j *Jira) SetAuth(authToken string) error {
	jiraAuth = token(authToken)
	if authToken == "test" {
		return nil
	}

	_, err := j.Request("GET", fmt.Sprintf("%s/myself", UrlBase), nil)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jira) SetAuthRepeat(times int, wait time.Duration, authToken string) error {
	for {
		err := j.SetAuth(authToken)
		if err != nil {
			if times == 0 {
				return err
			}
			times--
			time.Sleep(wait)
			continue
		}
		break
	}

	return nil
}

func (j *Jira) Request(method, url string, body io.Reader) ([]byte, error) {
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

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusAccepted {
		return []byte{}, errors.New(resp.Status)
	}

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return rBody, nil
}

func (j *Jira) addHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	jiraAuth.addAuth(req)
}

func (j *Jira) getAuth() auth {
	return jiraAuth
}
