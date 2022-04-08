package createmergerequest

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
)

var (
	ErrSourceBranchNotFound = errors.New("err_source_branch_not_found")
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
	repositoryPath := repository.GetRepoPath(args.Repository)
	repoDetail, err := m.repositorydm.GetByName(ctx, repositoryPath)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	remoteSourceBranch, err := m.repositorydm.GetBranch(ctx, repoDetail.ProjectID, args.SourceBranch)
	if err != nil {
		return errlib.WrapFunc(err)
	}
	if remoteSourceBranch == nil {
		return errlib.WrapFunc(errlib.WithFields(ErrSourceBranchNotFound, errlib.Fields{
			"name": args.SourceBranch,
		}))
	}

	mergeRequests := []repository.MergeRequest{}

	for _, targetBranch := range args.TargetBranches {
		mergeRequest, err := m.repositorydm.CreateMergeRequest(
			ctx,
			repoDetail.ProjectID,
			args.SourceBranch,
			targetBranch,
			args.Description,
			args.JiraTicketIDs,
		)

		if err != nil {
			return errlib.WrapFunc(err)
		}

		mergeRequests = append(mergeRequests, *mergeRequest)
	}

	for _, mergeRequest := range mergeRequests {
		fmt.Printf("%s: %s\n", mergeRequest.Title, mergeRequest.WebURL)
	}

	return nil
}

type Args struct {
	Repository     string
	SourceBranch   string
	TargetBranches []string
	Description    string
	JiraTicketIDs  []string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	a.Repository = flags["repository"]
	a.SourceBranch = flags["source-branch"]

	targetBranchVal := flags["target-branch"]
	a.TargetBranches = strings.Split(targetBranchVal, ",")

	a.Description = flags["description"]

	jiraVal := flags["jira"]
	a.JiraTicketIDs = strings.Split(jiraVal, ",")

	return a
}
