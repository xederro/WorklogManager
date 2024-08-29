package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Worklogs struct {
	StartAt    *int       `json:"startAt,omitempty"`
	MaxResults *int       `json:"maxResults,omitempty"`
	Total      *int       `json:"total,omitempty"`
	Worklogs   []*Worklog `json:"worklogs,omitempty"`
}

type Worklog struct {
	Self             *string `json:"self,omitempty"`
	Author           *Person `json:"author,omitempty"`
	UpdateAuthor     *Person `json:"updateAuthor,omitempty"`
	Comment          *string `json:"comment,omitempty"`
	Created          *string `json:"created,omitempty"`
	Updated          *string `json:"updated,omitempty"`
	Started          *string `json:"started,omitempty"`
	TimeSpent        *string `json:"timeSpent,omitempty"`
	TimeSpentSeconds *int    `json:"timeSpentSeconds,omitempty"`
	ID               *string `json:"id,omitempty"`
	IssueID          *string `json:"issueId,omitempty"`
}

func (w *Worklog) AddToIssue(i *Issue) error {
	body, err := json.Marshal(*w)
	if err != nil {
		return err
	}

	_, err = Jira{}.Request("POST", fmt.Sprintf("%s/worklog?adjustEstimate=auto", *i.Self), bytes.NewReader(body))
	if err != nil {
		return err
	}

	return nil
}
