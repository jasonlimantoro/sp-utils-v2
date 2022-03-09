package reviewmergerequest

import (
	"path/filepath"
	"testing"
)

func Test_renderMessage(t *testing.T) {
	type args struct {
		payload      SubstitutionPayload
		templatePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "with footer",
			args: args{
				payload: SubstitutionPayload{
					ReviewerUsername: "shannon.wong",
					Description:      "Test MR",
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
			},
			wantErr: false,
		},
		{
			name: "without footer",
			args: args{
				payload: SubstitutionPayload{
					ReviewerUsername: "shannon.wong",
					Description:      "Test MR",
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := renderMessage(tt.args.payload, tt.args.templatePath); (err != nil) != tt.wantErr {
				t.Errorf("renderMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
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
