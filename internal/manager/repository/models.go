package repository

import "git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"

type Repository struct {
	ProjectID int
	Name      string
}

type MergeRequest struct {
	*gitlab.MergeRequest
}
