package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/xederro/WorklogManager/config"
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
	client := &http.Client{}

	body, err := json.Marshal(*w)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/worklog?adjustEstimate=auto", *i.Self), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: config.SESSIONID,
	})

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.Status != "201" {
		return errors.New(string(rBody))
	}

	return nil
}
