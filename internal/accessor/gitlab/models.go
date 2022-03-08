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

type ListMergeRequestRequest struct {
	// ID is the Project.ID
	ID int `url:"id"`
	// State can be `opened`, `closed`, `locked` or `merged`
	State string `url:"state,omitempty"`
	// OrderBy can be `created_at`, `title`, `updated_at`, defaults to `created_at`
	OrderBy string `url:"order_by,omitempty"`
	// Sort can be `asc` or `desc`
	Sort string `url:"sort,omitempty"`

	AuthorID       int    `url:"author_id,omitempty"`
	AuthorUsername string `url:"author_username,omitempty"`
}

type MergeRequest struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
	WebURL       string `json:"web_url"`
	State        string `json:"state"`
}
