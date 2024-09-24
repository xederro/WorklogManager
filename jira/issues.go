package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Issues struct {
	Issues []*Issue `json:"issues,omitempty"`
}

type Issue struct {
	ID     *string `json:"id,omitempty"`
	Self   *string `json:"self,omitempty"`
	Key    *string `json:"key,omitempty"`
	Fields *Fields `json:"fields,omitempty"`
}

func (j Jira) GetIssues() (*Issues, error) {
	if a, ok := j.getAuth().(token); ok && a == "test" {
		return j.getTestIssues()
	}

	body, err := j.Request("GET", fmt.Sprintf("%s/search?jql=assignee in(currentUser())AND status!=closed", UrlBase), nil)
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

func (j Jira) getTestIssues() (*Issues, error) {
	body, err := os.ReadFile("data/exampleIssues.json")
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
