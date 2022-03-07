package gitlab

type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
	WebURL        string `json:"web_url"`
}

type CreateMergeRequestRequest struct {
	// ID is the Project.ID
	ID           int    `json:"id"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type MergeRequest struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
	WebURL       string `json:"web_url"`
}
