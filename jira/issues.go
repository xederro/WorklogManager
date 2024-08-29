package jira

import (
	"encoding/json"
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

func GetIssues() *Issues {
	body, err := Jira{}.Request("GET", fmt.Sprintf("%s/search?jql=assignee in(currentUser())AND sprint in openSprints()", UrlBase), nil)
	if err != nil {
		return nil
	}

	i := &Issues{}
	err = json.Unmarshal(body, i)
	if err != nil {
		return nil
	}

	return i
}
