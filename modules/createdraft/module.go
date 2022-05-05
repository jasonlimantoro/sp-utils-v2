package createdraft

import (
	"context"
	"io/ioutil"
	"strings"
	"time"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"github.com/gomarkdown/markdown"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/email"
)

const (
	SubjectTemplate    = "{prefix} {name} {today}"
	Name               = "Jason Gunawan Limantoro"
	PrefixWeeklyReport = "[Weekly Report]"
)

type Module interface {
	Do(ctx context.Context, args *Args) error
}

type module struct {
	emaildm email.Manager
}

func NewModule(emaildm email.Manager) *module {
	return &module{emaildm: emaildm}
}

func (m module) Do(ctx context.Context, args *Args) error {
	todayString := time.Now().Format("2006/01/02")
	subject := strings.NewReplacer(
		"{prefix}", PrefixWeeklyReport,
		"{name}", Name,
		"{today}", todayString,
	).Replace(SubjectTemplate)

	content, err := ioutil.ReadFile(args.InputFile)
	if err != nil {
		return errlib.WrapFunc(err)
	}
	html := string(markdown.ToHTML(content, nil, nil))

	if err := m.emaildm.CreateDraft(ctx, &email.CreateDraftRequest{
		Subject:     subject,
		To:          "buith@sea.com,roslim@sea.com",
		CC:          "limx@sea.com",
		ContentType: "text/html; charset=UTF-8",
		Body:        html,
	}); err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

type Args struct {
	InputFile string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	a.InputFile = flags["input-file"]

	return a
}
