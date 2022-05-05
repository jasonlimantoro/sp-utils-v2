package getweeklyupdates

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/lib"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/logger"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
)

const (
	DefaultDraftTemplate = `**What I have done this week**

{{ range $jira, $updates := .UpdatesMap -}}
- [{{ $jira.Title }}]({{ $jira.Link }})
{{range $updates }}  - {{.}}
{{ end }}{{ end }}
**What I will do next working week**

{{ range $jira, $updates := .UpdatesMap -}}
- [{{ $jira.Title }}]({{ $jira.Link }})
{{ end }}`
)

var (
	ErrDraftAlreadyExists = errors.New("err_draft_already_exists")
)

type Module interface {
	Do(ctx context.Context, args *Args) error
}

type module struct {
	taskmanager task.Manager
	logger      logger.Logger
}

func NewModule(taskmanager task.Manager, logger logger.Logger) *module {
	return &module{taskmanager: taskmanager, logger: logger}
}
func (m module) Do(ctx context.Context, args *Args) error {
	mondayDate := getMondayDate(time.Now(), args.DeltaWeek)
	weeklyUpdates, err := m.taskmanager.GetWeeklyUpdates(ctx, mondayDate)
	if err != nil {
		return errlib.WrapFunc(err)
	}
	substitution := SubstitutionPayload{
		weeklyUpdates,
	}

	var out io.Writer
	if args.OutputDirPath != "" {
		eowDateString := lib.GetWeekday(time.Now(), args.DeltaWeek, time.Friday).Format("2006-01-02")
		filePath := filepath.Join(args.OutputDirPath, fmt.Sprintf("%s.md", eowDateString))

		_, err := os.Stat(filePath)
		if !errors.Is(err, os.ErrNotExist) {
			return errlib.WrapFunc(errlib.WithFields(ErrDraftAlreadyExists, errlib.Fields{
				"filepath": filePath,
			}))
		}

		file, err := os.Create(filePath)
		if err != nil {
			return errlib.WrapFunc(err)
		}
		out = file
	} else {
		out = os.Stdout
	}

	if err := renderMessage(substitution, args.TemplateFilePath, out); err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

func renderMessage(payload SubstitutionPayload, templatePath string, out io.Writer) error {
	var t *template.Template
	if templatePath != "" {
		templateFullPath, _ := filepath.Abs(templatePath)
		t = template.Must(template.ParseFiles(templateFullPath))
	} else {
		t = template.Must(template.New("draft-template").Parse(DefaultDraftTemplate))
	}

	err := t.Execute(out, payload)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

type SubstitutionPayload struct {
	UpdatesMap map[task.JiraIssue][]string
}

type Args struct {
	DeltaWeek        int
	TemplateFilePath string
	OutputDirPath    string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	deltaWeekVal := flags["delta-week"]
	deltaWeekInt, _ := strconv.Atoi(deltaWeekVal)
	a.DeltaWeek = deltaWeekInt

	a.TemplateFilePath = flags["template"]
	a.OutputDirPath = flags["out"]

	return a
}
