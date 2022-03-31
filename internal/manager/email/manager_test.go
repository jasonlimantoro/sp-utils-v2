package email

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_renderBody(t *testing.T) {
	type args struct {
		payload      SubstitutionPayload
		templateText string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				payload: SubstitutionPayload{
					Subject:     "[Weekly Report] Jason Gunawan",
					To:          "abc.@shopee.com,def@shopee.com",
					CC:          "ghj@shopee.com",
					BCC:         "klm@shopee.com",
					ContentType: "text/html; charset=UTF-8",
					Body:        "hello world",
				},
				templateText: payloadTemplate,
			},
			wantOut: `To: abc.@shopee.com,def@shopee.com
CC: ghj@shopee.com
BCC: klm@shopee.com
Subject: [Weekly Report] Jason Gunawan
Content-Type: text/html; charset=UTF-8

hello world`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := renderBody(tt.args.payload, tt.args.templateText, out)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
