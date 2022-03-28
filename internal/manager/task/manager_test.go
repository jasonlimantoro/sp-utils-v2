package task

import (
	"reflect"
	"testing"
	"time"
)

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

func Test_generateWeekStringMap(t *testing.T) {
	type args struct {
		startDate time.Time
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{
			name: "normal",
			args: args{
				startDate: time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
			},
			want: map[string]bool{
				"28-Mar-2022": true,
				"29-Mar-2022": true,
				"30-Mar-2022": true,
				"31-Mar-2022": true,
				"01-Apr-2022": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateWeekStringMap(tt.args.startDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateWeekStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
