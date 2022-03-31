package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"strings"
	"text/template"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gmail"
)

type Manager interface {
	CreateDraft(ctx context.Context, request *CreateDraftRequest) error
}

type manager struct {
	accessor gmail.Accessor
}

func NewManager(accessor gmail.Accessor) *manager {
	return &manager{accessor: accessor}
}

func (m manager) CreateDraft(ctx context.Context, request *CreateDraftRequest) error {
	buf := &bytes.Buffer{}

	if err := renderBody(*request, payloadTemplate, buf); err != nil {
		return errlib.WrapFunc(err)
	}

	body := buf.String()
	body = base64.StdEncoding.EncodeToString([]byte(body))
	body = strings.Replace(body, "/", "_", -1)
	body = strings.Replace(body, "+", "-", -1)
	body = strings.Replace(body, "=", "", -1)

	if err := m.accessor.CreateDraft(ctx, body); err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}

type SubstitutionPayload = CreateDraftRequest

const payloadTemplate = `To: {{ .To }}
CC: {{ .CC }}
BCC: {{ .BCC }}
Subject: {{ .Subject }}
Content-Type: {{ .ContentType }}

{{ .Body }}`

func renderBody(payload SubstitutionPayload, templateText string, out io.Writer) error {
	t := template.Must(template.New("email-payload").Parse(templateText))

	err := t.Execute(out, payload)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}
