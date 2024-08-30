package jira

import (
	"encoding/json"
	"errors"
	"fmt"
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

func GetIssues() (*Issues, error) {
	body, err := Jira{}.Request("GET", fmt.Sprintf("%s/search?jql=assignee in(currentUser())AND sprint in openSprints()", UrlBase), nil)
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
