package repository

type Repository struct {
	ProjectID int
	Name      string
}

type MergeRequest struct {
	Title        string
	WebURL       string
	TargetBranch string
	SourceBranch string
	State        string
}
