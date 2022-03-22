package createcard

import (
	"context"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
)

type Module interface {
	Do(ctx context.Context, Args *Args) error
}

type module struct {
	manager task.Manager
}

func (m module) Do(ctx context.Context, args *Args) error {
	err := m.manager.CreateInList(
		ctx,
		args.ListName,
		args.Title,
		args.JiraLink,
		args.EpicLink,
		args.TDLink,
		args.PRDLink,
	)

	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

func NewModule(manager task.Manager) *module {
	return &module{manager: manager}
}

type Args struct {
	Title    string
	ListName string
	JiraLink string
	EpicLink string
	PRDLink  string
	TDLink   string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	a.Title = flags["title"]
	a.ListName = flags["list-name"]
	a.JiraLink = flags["jira-link"]
	a.EpicLink = flags["epic-link"]
	a.PRDLink = flags["prd-link"]
	a.TDLink = flags["td-link"]

	return a
}
