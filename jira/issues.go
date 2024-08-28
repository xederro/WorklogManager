package jira

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/xederro/WorklogManager/config"
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
	client := &http.Client{}

	req, err := http.NewRequest("GET", config.BASEJIRALINK+"search?jql=assignee%20in(currentUser())AND%20sprint%20in%20openSprints()", nil)
	if err != nil {
		log.Fatalln(err, 0)
	}

	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: config.SESSIONID,
	})

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err, 1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err, 2)
	}

	i := &Issues{}
	json.Unmarshal(body, i)

	return i
}
