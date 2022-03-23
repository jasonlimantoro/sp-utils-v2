package createcard

import (
	"context"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"github.com/sirupsen/logrus"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
)

type Module interface {
	Do(ctx context.Context, Args *Args) error
}

type module struct {
	logger  *logrus.Logger
	manager task.Manager
}

func (m module) Do(ctx context.Context, args *Args) error {
	card, err := m.manager.CreateInList(
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

	m.logger.WithFields(logrus.Fields{
		"name": card.Name,
		"url":  card.URL,
	}).Info("created")

	return nil
}

func NewModule(manager task.Manager, logger *logrus.Logger) *module {
	return &module{manager: manager, logger: logger}
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
