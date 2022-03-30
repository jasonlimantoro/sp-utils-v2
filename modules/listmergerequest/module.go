package listmergerequest

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
)

var (
	TargetBranchToCompare = []string{"test", "uat", "master", "staging"}
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
	repositoryData, err := m.repositorydm.GetByName(ctx, args.Repository)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	matchingMergeRequests, err := m.repositorydm.ListMergeRequests(
		ctx,
		repositoryData.ProjectID,
		args.JiraTicketIDs,
		args.State,
	)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	bySourceBranch := make(map[string][]*repository.MergeRequest)

	for _, mr := range matchingMergeRequests {
		if _, ok := bySourceBranch[mr.SourceBranch]; !ok {
			bySourceBranch[mr.SourceBranch] = []*repository.MergeRequest{mr}
		} else {
			bySourceBranch[mr.SourceBranch] = append(bySourceBranch[mr.SourceBranch], mr)
		}
	}

	for sourceBranch, mergeRequests := range bySourceBranch {
		for _, mr := range mergeRequests {
			fmt.Printf("%s|%s: %s (%s)\n", mr.TargetBranch, mr.Title, mr.WebURL, mr.State)
		}

		fmt.Println()

		fmt.Printf("Branch comparison links for %s: \n", sourceBranch)

		for _, targetBranch := range TargetBranchToCompare {
			fmt.Printf("%s|https://%s/%s/-/compare/%s...%s\n",
				targetBranch,
				gitlab.GitlabHost,
				args.Repository,
				targetBranch,
				url.QueryEscape(sourceBranch),
			)
		}
	}

	return nil
}

type Args struct {
	Repository    string
	JiraTicketIDs []string
	State         string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	repositoryVal, _ := flags["repository"]
	jiraVal, _ := flags["jira"]
	stateVal, _ := flags["state"]

	a.Repository = repositoryVal
	a.JiraTicketIDs = strings.Split(jiraVal, ",")
	a.State = stateVal

	return a
}
