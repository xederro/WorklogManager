package jira

type Fields struct {
	AggregateTimeOriginalEstimate *int      `json:"aggregatetimeoriginalestimate,omitempty"`
	Assignee                      *Person   `json:"assignee,omitempty"`
	Reporter                      *Person   `json:"reporter,omitempty"`
	Creator                       *Person   `json:"creator,omitempty"`
	IssueType                     *Name     `json:"issuetype,omitempty"`
	Description                   *string   `json:"description,omitempty"`
	Summary                       *string   `json:"summary,omitempty"`
	Priority                      *Name     `json:"priority,omitempty"`
	Status                        *Name     `json:"status,omitempty"`
	TimeOriginalEstimate          *int      `json:"timeoriginalestimate,omitempty"`
	AggregateTimeEstimate         *int      `json:"aggregatetimeestimate,omitempty"`
	Progress                      *Progress `json:"progress,omitempty"`
	AggregateProgress             *Progress `json:"aggregateprogress,omitempty"`
	TimeSpent                     *int      `json:"timespent,omitempty"`
	AggregateTimeSpent            *int      `json:"aggregatetimespent,omitempty"`
	Workratio                     *int      `json:"workratio,omitempty"`
}
