package task

import "testing"

func Test_jiraTicketIDFromLink(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				link: "https://jira.shopee.io/browse/SPOT-35892",
			},
			want: "SPOT-35892",
		},
		{
			name: "invalid input",
			args: args{
				link: "https://jira.shopee.io/browse/abcd",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jiraTicketIDFromLink(tt.args.link); got != tt.want {
				t.Errorf("jiraTicketIDFromLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
