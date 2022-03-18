package reviewmergerequest

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_renderMessage(t *testing.T) {
	type args struct {
		payload      SubstitutionPayload
		templatePath string
		out          io.Writer
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantOutString string
	}{
		{
			name: "with footer",
			args: args{
				payload: SubstitutionPayload{
					ReviewerUsername: "shannon.wong",
					Description:      "Test MR",
					JiraLink:         "https://jira.shopee.io/browse/SPOT-36226",
					MergeRequests: []SubstitutionMergeRequest{
						{
							RepoName:     "bridge",
							TargetBranch: "master",
							Link:         "https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/914",
						},
						{
							RepoName:     "bridge",
							TargetBranch: "uat",
							Link:         "https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/913",
						},
					},
					Footer: "Can also publish to master topic?",
				},
				templatePath: absPath("code-review.tpl"),
				out:          &bytes.Buffer{},
			},
			wantErr: false,
			wantOutString: `Hi @shannon.wong, please review the following:
**Test MR**

- bridge|master: https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/914
- bridge|uat: https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/913

Jira: https://jira.shopee.io/browse/SPOT-36226

Can also publish to master topic?

Thank you! :capoo-thanks:
`,
		},
		{
			name: "without footer",
			args: args{
				payload: SubstitutionPayload{
					ReviewerUsername: "shannon.wong",
					Description:      "Test MR",
					JiraLink:         "https://jira.shopee.io/browse/SPOT-36226",
					MergeRequests: []SubstitutionMergeRequest{
						{
							RepoName:     "bridge",
							TargetBranch: "master",
							Link:         "https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/914",
						},
						{
							RepoName:     "bridge",
							TargetBranch: "uat",
							Link:         "https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/913",
						},
					},
				},
				templatePath: absPath("code-review.tpl"),
				out:          &bytes.Buffer{},
			},
			wantErr: false,
			wantOutString: `Hi @shannon.wong, please review the following:
**Test MR**

- bridge|master: https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/914
- bridge|uat: https://git.garena.com/shopee/pl/marketplace-payment/-/merge_requests/913

Jira: https://jira.shopee.io/browse/SPOT-36226

Thank you! :capoo-thanks:
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := renderMessage(tt.args.payload, tt.args.templatePath, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("renderMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			resBuf := tt.args.out.(*bytes.Buffer)
			assert.Equal(t, tt.wantOutString, resBuf.String())
		})
	}
}

func absPath(path string) string {
	res, _ := filepath.Abs(path)

	return res
}

func Test_cleanTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				title: "[SPOT-35581][Master] Add spl_order_price inside OrderAmountDetail",
			},
			want: "Add spl_order_price inside OrderAmountDetail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanTitle(tt.args.title); got != tt.want {
				t.Errorf("cleanTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
