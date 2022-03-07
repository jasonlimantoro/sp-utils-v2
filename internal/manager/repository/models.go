package repository

type Repository struct {
	ProjectID int
	Name      string
}

type MergeRequest struct {
	ID     int
	Title  string
	WebURL string
}
