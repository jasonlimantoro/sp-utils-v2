package task

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/trello"
)

type Manager interface {
	CreateInList(ctx context.Context, listName string, title string, jiraLink string, epicLink string, TDLink string, PRDLink string) error
	CreateList(ctx context.Context, name string) error
}

type manager struct {
	accessor trello.Accessor
}

func NewManager(accessor trello.Accessor) *manager {
	return &manager{accessor: accessor}
}

func (m manager) CreateInList(ctx context.Context, listName string, title string, jiraLink string, epicLink string, TDLink string, PRDLink string) error {
	var list *trello.List
	lists, err := m.accessor.GetList(ctx, trello.BoardID)
	if err != nil {
		return errlib.WrapFunc(err)
	}
	for _, v := range lists {
		if v.Name == listName {
			list = v
		}
	}

	description := m.buildDescription(jiraLink, epicLink, TDLink, PRDLink)
	name := fmt.Sprintf("[%s]%s", strings.ToUpper(jiraTicketIDFromLink(jiraLink)), title)
	_, err = m.accessor.CreateCard(ctx, list.ID, name, description)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

func (m manager) buildDescription(jiraLink string, epicLink string, TDLink string, PRDLink string) string {
	infos := []string{}
	infos = append(infos, fmt.Sprintf("Jira: %s", jiraLink))
	if epicLink != "" {
		infos = append(infos, fmt.Sprintf("Epic: %s", epicLink))
	}
	if TDLink != "" {
		infos = append(infos, fmt.Sprintf("TD: %s", TDLink))
	}
	if PRDLink != "" {
		infos = append(infos, fmt.Sprintf("PRD: %s", PRDLink))
	}

	return strings.Join(infos, "\n")
}

func jiraTicketIDFromLink(link string) string {
	r := regexp.MustCompile(`https://jira.shopee.io/browse/(\w+-\d+)`)
	match := r.FindStringSubmatch(link)
	if len(match) < 1 {
		return ""
	}

	return match[1]
}

func (m manager) CreateList(ctx context.Context, name string) error {
	lists, err := m.accessor.GetList(ctx, trello.BoardID)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	firstListPos := lists[0].Pos
	lastWorkingDayPos := lists[1].Pos
	desiredPos := firstListPos + (lastWorkingDayPos-firstListPos)/2

	_, err = m.accessor.CreateList(ctx, trello.BoardID, name, desiredPos)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}
