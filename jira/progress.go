package jira

type Progress struct {
	Progress *int `json:"progress,omitempty"`
	Total    *int `json:"total,omitempty"`
	Percent  *int `json:"percent,omitempty"`
}
