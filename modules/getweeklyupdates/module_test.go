package getweeklyupdates

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
)

func Test_renderMessage(t *testing.T) {
	type args struct {
		payload      SubstitutionPayload
		templatePath string
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
					UpdatesMap: map[task.JiraIssue][]string{
						task.JiraIssue{
							ID:    "SPOT-1234",
							Title: "Jira Task 1",
							Link:  "https://jira.shopee.io/browse/SPOT-1234",
						}: {"update 1", "update 2"},
						task.JiraIssue{
							ID:    "SPOT-2345",
							Title: "Jira Task 2",
							Link:  "https://jira.shopee.io/browse/SPOT-2345",
						}: {"update 3", "update 4"},
					},
				},
				templatePath: absPath("draft.tpl"),
			},
			wantOut: `**What I have done this week**

- [Jira Task 1](https://jira.shopee.io/browse/SPOT-1234)
  - update 1
  - update 2
- [Jira Task 2](https://jira.shopee.io/browse/SPOT-2345)
  - update 3
  - update 4

**What I will do next working week**

- [Jira Task 1](https://jira.shopee.io/browse/SPOT-1234)
- [Jira Task 2](https://jira.shopee.io/browse/SPOT-2345)
`,
			wantErr: false,
		},
		{
			name: "without file",
			args: args{
				payload: SubstitutionPayload{
					UpdatesMap: map[task.JiraIssue][]string{
						task.JiraIssue{
							ID:    "SPOT-1234",
							Title: "Jira Task 1",
							Link:  "https://jira.shopee.io/browse/SPOT-1234",
						}: {"update 1", "update 2"},
						task.JiraIssue{
							ID:    "SPOT-2345",
							Title: "Jira Task 2",
							Link:  "https://jira.shopee.io/browse/SPOT-2345",
						}: {"update 3", "update 4"},
					},
				},
				templatePath: "",
			},
			wantOut: `**What I have done this week**

- [Jira Task 1](https://jira.shopee.io/browse/SPOT-1234)
  - update 1
  - update 2
- [Jira Task 2](https://jira.shopee.io/browse/SPOT-2345)
  - update 3
  - update 4

**What I will do next working week**

- [Jira Task 1](https://jira.shopee.io/browse/SPOT-1234)
- [Jira Task 2](https://jira.shopee.io/browse/SPOT-2345)
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := renderMessage(tt.args.payload, tt.args.templatePath, out); (err != nil) != tt.wantErr {
				t.Errorf("renderMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

func absPath(path string) string {
	res, _ := filepath.Abs(path)

	return res
}
