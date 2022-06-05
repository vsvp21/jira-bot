package jira

type Issue struct {
	Fields Fields `json:"fields"`
}

type Fields struct {
	Project     Project   `json:"project,omitempty"`
	IssueType   IssueType `json:"issuetype,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Description string    `json:"description,omitempty"`
}

type Project struct {
	ID string `json:"id,omitempty"`
}

type IssueType struct {
	ID string `json:"id,omitempty"`
}

type CreatedIssue struct {
	ID   string
	Key  string
	Self string
}
