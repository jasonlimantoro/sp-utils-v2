package task

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/trello"
)

type Manager interface {
	CreateInList(ctx context.Context, listName string, title string, jiraLink string, epicLink string, TDLink string, PRDLink string) (*Task, error)
	CreateList(ctx context.Context, name string) (*List, error)
	GetWeeklyUpdates(ctx context.Context, startDate time.Time) (map[JiraIssue][]string, error)
}

type manager struct {
	accessor trello.Accessor
}

func NewManager(accessor trello.Accessor) *manager {
	return &manager{accessor: accessor}
}

func (m manager) CreateInList(ctx context.Context, listName string, title string, jiraLink string, epicLink string, TDLink string, PRDLink string) (*Task, error) {
	var list *trello.List
	lists, err := m.accessor.GetList(ctx, trello.BoardID)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}
	for _, v := range lists {
		if v.Name == listName {
			list = v
		}
	}

	description := m.buildDescription(jiraLink, epicLink, TDLink, PRDLink)
	name := fmt.Sprintf("[%s]%s", strings.ToUpper(jiraTicketIDFromLink(jiraLink)), title)
	card, err := m.accessor.CreateCard(ctx, list.ID, name, description)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return &Task{
		Name: card.Name,
		URL:  card.URL,
	}, nil
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

func (m manager) CreateList(ctx context.Context, name string) (*List, error) {
	lists, err := m.accessor.GetList(ctx, trello.BoardID)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	firstListPos := lists[0].Pos
	lastWorkingDayPos := lists[1].Pos
	desiredPos := firstListPos + (lastWorkingDayPos-firstListPos)/2

	list, err := m.accessor.CreateList(ctx, trello.BoardID, name, desiredPos)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return &List{
		Name: list.Name,
		Pos:  list.Pos,
	}, nil
}

func (m manager) GetWeeklyUpdates(ctx context.Context, startDate time.Time) (map[JiraIssue][]string, error) {
	lists, err := m.accessor.GetList(ctx, trello.BoardID)
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}
	weekStringsMap := generateWeekStringMap(startDate)
	jiraIDToUpdates := make(map[JiraIssue][]string)

	for i := range lists {
		list := lists[len(lists)-1-i]
		if _, ok := weekStringsMap[list.Name]; ok {
			cards, err := m.accessor.GetCards(ctx, list.ID)
			if err != nil {
				return nil, errlib.WrapFunc(err)
			}

			for _, card := range cards {
				jiraIssue := (JiraIssue{}).FromCardTitle(card.Name)
				comments, err := m.accessor.GetCardActions(ctx, &trello.GetCardActionsRequest{
					CardID: card.ID,
					Filter: trello.ActionCommentCard,
				})

				if err != nil {
					return nil, errlib.WrapFunc(err)
				}

				for _, comment := range comments {
					message := getMessageFromUpdateComment(comment.Data.Text)
					if len(message) > 0 {
						if _, ok := jiraIDToUpdates[jiraIssue]; !ok {
							jiraIDToUpdates[jiraIssue] = []string{message}
						} else {
							jiraIDToUpdates[jiraIssue] = append(jiraIDToUpdates[jiraIssue], message)
						}
					}
				}
			}
		}
	}

	return jiraIDToUpdates, nil
}

func getMessageFromUpdateComment(text string) string {
	re := regexp.MustCompile(`(?i)UPDATE: (.*)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func generateWeekStringMap(startDate time.Time) map[string]bool {
	res := map[string]bool{}

	for deltaDay := 0; deltaDay < 5; deltaDay++ {
		currentDate := startDate.AddDate(0, 0, deltaDay)
		currentDateString := currentDate.Format("02-Jan-2006")
		res[currentDateString] = true
	}

	return res
}
