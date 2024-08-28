package jira

type Fields struct {
	Aggregatetimeoriginalestimate *int      `json:"aggregatetimeoriginalestimate,omitempty"`
	Assignee                      *Person   `json:"assignee,omitempty"`
	Reporter                      *Person   `json:"reporter,omitempty"`
	Creator                       *Person   `json:"creator,omitempty"`
	Issuetype                     *Name     `json:"issuetype,omitempty"`
	Description                   *string   `json:"description,omitempty"`
	Summary                       *string   `json:"summary,omitempty"`
	Priority                      *Name     `json:"priority,omitempty"`
	Status                        *Name     `json:"status,omitempty"`
	Timeoriginalestimate          *int      `json:"timeoriginalestimate,omitempty"`
	Aggregatetimeestimate         *int      `json:"aggregatetimeestimate,omitempty"`
	Progress                      *Progress `json:"progress,omitempty"`
	Aggregateprogress             *Progress `json:"aggregateprogress,omitempty"`
	Timespent                     *int      `json:"timespent,omitempty"`
	Aggregatetimespent            *int      `json:"aggregatetimespent,omitempty"`
	Workratio                     *int      `json:"workratio,omitempty"`
}
