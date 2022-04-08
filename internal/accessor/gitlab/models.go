package gitlab

import "time"

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

type GetBranchRequest struct {
	ProjectID int    `url:"-"`
	Search    string `url:"search"`
}

type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		ID             string      `json:"id"`
		ShortID        string      `json:"short_id"`
		CreatedAt      time.Time   `json:"created_at"`
		ParentIds      interface{} `json:"parent_ids"`
		Title          string      `json:"title"`
		Message        string      `json:"message"`
		AuthorName     string      `json:"author_name"`
		AuthorEmail    string      `json:"author_email"`
		AuthoredDate   time.Time   `json:"authored_date"`
		CommitterName  string      `json:"committer_name"`
		CommitterEmail string      `json:"committer_email"`
		CommittedDate  time.Time   `json:"committed_date"`
		WebURL         string      `json:"web_url"`
	} `json:"commit"`
	Merged             bool   `json:"merged"`
	Protected          bool   `json:"protected"`
	DevelopersCanPush  bool   `json:"developers_can_push"`
	DevelopersCanMerge bool   `json:"developers_can_merge"`
	CanPush            bool   `json:"can_push"`
	Default            bool   `json:"default"`
	WebURL             string `json:"web_url"`
}
