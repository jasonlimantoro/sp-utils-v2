package repository

import (
	"context"
	"fmt"
	"strings"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
)

type Manager interface {
	GetByName(ctx context.Context, name string) (*Repository, error)
	CreateMergeRequest(
		ctx context.Context,
		projectID int,
		sourceBranch string,
		targetBranch string,
		description string,
		jiraTicketIDs []string,
	) (*MergeRequest, error)
	ListMergeRequests(ctx context.Context, projectID int, jiraTicketIDs []string) ([]*MergeRequest, error)
}

type manager struct {
	accessor gitlab.Accessor
}

func NewManager(accessor gitlab.Accessor) *manager {
	return &manager{accessor: accessor}
}

func (m manager) GetByName(ctx context.Context, name string) (*Repository, error) {
	res, err := m.accessor.GetProjectByName(ctx, name)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return &Repository{
		ProjectID: res.ID,
		Name:      res.Name,
	}, nil
}

func (m manager) CreateMergeRequest(
	ctx context.Context,
	projectID int,
	sourceBranch string,
	targetBranch string,
	description string,
	jiraTicketIDs []string,
) (*MergeRequest, error) {
	var jiraTicketIDsString string
	for _, jiraTicketId := range jiraTicketIDs {
		jiraTicketIDsString += fmt.Sprintf("[%s]", strings.ToUpper(jiraTicketId))
	}

	title := fmt.Sprintf("%s[%s] %s", jiraTicketIDsString, targetBranch, description)

	res, err := m.accessor.CreateMergeRequest(ctx, &gitlab.CreateMergeRequestRequest{
		ID:           projectID,
		SourceBranch: sourceBranch,
		TargetBranch: targetBranch,
		Title:        title,
	})
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return &MergeRequest{
		res,
	}, nil
}

func (m manager) ListMergeRequests(ctx context.Context, projectID int, jiraTicketIDs []string) ([]*MergeRequest, error) {
	myMergeRequests, err := m.accessor.ListMergeRequests(ctx, &gitlab.ListMergeRequestRequest{
		ID:             projectID,
		State:          "opened",
		AuthorUsername: "jason.limantoro",
	})

	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	res := []*MergeRequest{}
	for _, mr := range myMergeRequests {
		if titleContainsJiraTickets(mr.Title, jiraTicketIDs) {
			res = append(res, &MergeRequest{mr})
		}
	}

	return res, nil
}

func titleContainsJiraTickets(title string, jiraTicketIDs []string) bool {
	for _, jiraTicketID := range jiraTicketIDs {
		if strings.Contains(strings.ToUpper(title), strings.ToUpper(jiraTicketID)) {
			return true
		}
	}

	return false
}
