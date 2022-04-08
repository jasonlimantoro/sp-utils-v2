package repository

type Repository struct {
	ProjectID int
	Name      string
}

type Branch struct {
	Name   string `json:"name"`
	Merged bool   `json:"merged"`
	WebURL string `json:"web_url"`
}

type MergeRequest struct {
	Title        string
	WebURL       string
	TargetBranch string
	SourceBranch string
	State        string
}
