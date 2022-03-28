package task

import (
	"fmt"
	"regexp"
)

type Task struct {
	Name string
	URL  string
}

type List struct {
	ID   string
	Name string
	Pos  float64
}

type JiraIssue struct {
	ID    string
	Title string
	Link  string
}

func (j JiraIssue) FromCardTitle(title string) JiraIssue {
	re := regexp.MustCompile(`\[(\w+-\w+)\](.+)`)
	matches := re.FindStringSubmatch(title)
	if len(matches) > 1 {
		return JiraIssue{ID: matches[1], Title: matches[2], Link: buildJiraLink(matches[1])}
	}

	return JiraIssue{}
}

func buildJiraLink(jiraTicketID string) string {
	return fmt.Sprintf("https://jira.shopee.io/browse/%s", jiraTicketID)
}
