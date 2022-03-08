package listmergerequest

import (
	"context"
	"fmt"
	"strings"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
)

type Module interface {
	Do(ctx context.Context, args *Args) error
}

type module struct {
	repositorydm repository.Manager
}

func NewModule(repositorydm repository.Manager) *module {
	return &module{repositorydm: repositorydm}
}

func (m module) Do(ctx context.Context, args *Args) error {
	repository, err := m.repositorydm.GetByName(ctx, args.Repository)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	matchingMergeRequests, err := m.repositorydm.ListMergeRequests(ctx, repository.ProjectID, args.JiraTicketIDs)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	for _, mr := range matchingMergeRequests {
		fmt.Printf("%s: %s -> %s\n", mr.Title, mr.WebURL, mr.State)
	}

	return nil
}

type Args struct {
	Repository    string
	JiraTicketIDs []string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	repositoryVal, _ := flags["repository"]
	jiraVal, _ := flags["jira"]

	a.Repository = repositoryVal
	a.JiraTicketIDs = strings.Split(jiraVal, ",")

	return a
}
