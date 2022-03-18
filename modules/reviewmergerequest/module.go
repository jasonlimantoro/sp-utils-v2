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
	repoData, err := m.repositorydm.GetByName(ctx, args.Repository)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	mergeRequests, err := m.repositorydm.ListMergeRequests(ctx, repoData.ProjectID, args.JiraTicketIDs, "opened")
	if err != nil {
		return errlib.WrapFunc(err)
	}

	substitutionPayload := constructSubstitutionPayload(args.Repository, "", mergeRequests)

	templatePath, _ := filepath.Abs(args.TemplateFilePath)
	if err := renderMessage(substitutionPayload, templatePath, os.Stdout); err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

func constructSubstitutionPayload(repo string, reviewer string, mergeRequests []*repository.MergeRequest) SubstitutionPayload {
	res := SubstitutionPayload{}

	finalReviewer := reviewer
	if finalReviewer == "" {
		finalReviewer = repoToRecommendedReviewerMapping[repo][0]
	}
	res.ReviewerUsername = reviewerToMattermostUsernameMapping[finalReviewer]

	firstMergeRequests := mergeRequests[0]
	res.Description = cleanTitle(firstMergeRequests.Title)

	jiraMap := map[string]string{}
	repoName := getRepoName(repo)

	for _, mr := range mergeRequests {
		res.MergeRequests = append(res.MergeRequests, SubstitutionMergeRequest{
			RepoName:     repoName,
			TargetBranch: mr.TargetBranch,
			Link:         mr.WebURL,
		})
		for _, jiraTicketID := range mr.GetRelatedJiraTickets() {
			if _, ok := jiraMap[jiraTicketID]; !ok {
				jiraMap[jiraTicketID] = buildJiraLink(jiraTicketID)
			}
		}
	}

	jiraLinks := []string{}
	for _, link := range jiraMap {
		jiraLinks = append(jiraLinks, link)
	}
	res.JiraLink = strings.Join(jiraLinks, ",")

	return res
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
	RepoName     string
	TargetBranch string
	Link         string
}

type Args struct {
	Repository       string
	JiraTicketIDs    []string
	TemplateFilePath string
}

func (a *Args) FromMap(flag map[string]string) *Args {
	repositoryVal := flag["repository"]
	jiraVal := flag["jira"]
	templateVal := flag["template"]

	a.Repository = repositoryVal
	a.JiraTicketIDs = strings.Split(jiraVal, ",")
	a.TemplateFilePath = templateVal

	return a
}
