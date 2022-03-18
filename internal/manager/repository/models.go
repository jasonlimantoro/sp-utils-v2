package repository

import (
	"regexp"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
)

type Repository struct {
	ProjectID int
	Name      string
}

type MergeRequest struct {
	*gitlab.MergeRequest
}

func (m MergeRequest) GetRelatedJiraTickets() []string {
	re := regexp.MustCompile(`\[(\w+-\w+)\]`)
	result := []string{}

	matches := re.FindAllStringSubmatch(m.Title, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			if len(m) > 0 {
				result = append(result, m[1])
			}
		}
	}

	return result
}
