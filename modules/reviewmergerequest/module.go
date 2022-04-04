package reviewmergerequest

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
)

const (
	DefaultCodeReviewMessageTemplate = "modules/reviewmergerequest/code-review.tpl"
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
	substitutionMergeRequests := []SubstitutionMergeRequest{}
	for _, repoShortName := range args.Repositories {
		repoData, err := m.repositorydm.GetByName(ctx, repository.RepoToPathMapping[repoShortName])
		if err != nil {
			return errlib.WrapFunc(err)
		}

		mergeRequests, err := m.repositorydm.ListMergeRequests(ctx, repoData.ProjectID, args.JiraTicketIDs, "opened", "")
		if err != nil {
			return errlib.WrapFunc(err)
		}

		for _, mr := range mergeRequests {
			substitutionMergeRequests = append(substitutionMergeRequests, SubstitutionMergeRequest{
				Title:        mr.Title,
				TargetBranch: mr.TargetBranch,
				Link:         mr.WebURL,
				RepoName:     repoShortName,
			})
		}

	}

	substitutionPayload := constructSubstitutionPayload("", substitutionMergeRequests)

	templatePath, _ := filepath.Abs(args.TemplateFilePath)
	if err := renderMessage(substitutionPayload, templatePath, os.Stdout); err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

func constructSubstitutionPayload(reviewer string, smrs []SubstitutionMergeRequest) SubstitutionPayload {
	if len(smrs) == 0 {
		return SubstitutionPayload{}
	}

	firstMergeRequest := smrs[0]

	finalReviewer := reviewer
	if finalReviewer == "" {
		finalReviewer = repoToRecommendedReviewerMapping[firstMergeRequest.RepoName][0]
	}

	jiraMap := map[string]string{}
	for _, smr := range smrs {
		for _, jiraTicketID := range smr.GetRelatedJiraTickets() {
			if _, ok := jiraMap[jiraTicketID]; !ok {
				jiraMap[jiraTicketID] = buildJiraLink(jiraTicketID)
			}
		}
	}

	jiraLinks := []string{}
	for _, link := range jiraMap {
		jiraLinks = append(jiraLinks, link)
	}

	return SubstitutionPayload{
		ReviewerUsername: reviewerToMattermostUsernameMapping[finalReviewer],
		Description:      cleanTitle(firstMergeRequest.Title),
		MergeRequests:    smrs,
		JiraLink:         strings.Join(jiraLinks, ","),
	}
}

func buildJiraLink(jiraTicketID string) string {
	return fmt.Sprintf("https://jira.shopee.io/browse/%s", jiraTicketID)
}

func cleanTitle(title string) string {
	r := regexp.MustCompile(`\[.*\]`)
	removed := r.ReplaceAllString(title, "")
	return strings.TrimSpace(removed)
}

func renderMessage(payload SubstitutionPayload, templatePath string, out io.Writer) error {
	t := template.Must(template.ParseFiles(templatePath))

	err := t.Execute(out, payload)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

type SubstitutionPayload struct {
	ReviewerUsername string
	Description      string
	MergeRequests    []SubstitutionMergeRequest
	JiraLink         string
	Footer           string
}

type SubstitutionMergeRequest struct {
	Title        string
	RepoName     string
	TargetBranch string
	Link         string
}

func (s SubstitutionMergeRequest) GetRelatedJiraTickets() []string {
	re := regexp.MustCompile(`\[(\w+-\w+)\]`)
	result := []string{}

	matches := re.FindAllStringSubmatch(s.Title, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			if len(m) > 0 {
				result = append(result, m[1])
			}
		}
	}

	return result
}

type Args struct {
	Repositories     []string
	JiraTicketIDs    []string
	TemplateFilePath string
}

func (a *Args) FromMap(flag map[string]string) *Args {
	a.Repositories = strings.Split(flag["repository"], ",")
	a.JiraTicketIDs = strings.Split(flag["jira"], ",")
	a.TemplateFilePath = flag["template"]

	return a
}
