package jira

import (
	"encoding/json"
	"errors"
	_ "embed"
	"fmt"
)
//go:embed data/exampleIssues.json
var exampleIssues []byte

type Issues struct {
	Issues []*Issue `json:"issues,omitempty"`
}

type Issue struct {
	ID     *string `json:"id,omitempty"`
	Self   *string `json:"self,omitempty"`
	Key    *string `json:"key,omitempty"`
	Fields *Fields `json:"fields,omitempty"`
}

func (j *Jira) GetIssues() (*Issues, error) {
	if a, ok := j.getAuth().(token); ok && a == "test" {
		return j.getTestIssues()
	}

	query := "?jql=assignee%20in(currentUser())AND%20status!=closed"
	body, err := j.Request("GET", fmt.Sprintf("%s/search%s", UrlBase, query), nil)
	if err != nil {
		return nil, err
	}

	i := &Issues{}
	err = json.Unmarshal(body, i)
	if err != nil {
		return nil, err
	}

	if len(i.Issues) == 0 {
		return nil, errors.New("no issues found")
	}

	return i, nil
}

func (j *Jira) getTestIssues() (*Issues, error) {
	i := &Issues{}
	err := json.Unmarshal(exampleIssues, i)
	if err != nil {
		return nil, err
	}

	if len(i.Issues) == 0 {
		return nil, errors.New("no issues found")
	}

	return i, nil
}
