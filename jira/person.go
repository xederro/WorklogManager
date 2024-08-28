package jira

type Person struct {
	Self         *string `json:"self,omitempty"`
	Name         *string `json:"name,omitempty"`
	Key          *string `json:"key,omitempty"`
	EmailAddress *string `json:"emailAddress,omitempty"`
	DisplayName  *string `json:"displayName,omitempty"`
}
