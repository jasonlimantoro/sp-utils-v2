package reviewmergerequest

import "testing"

func Test_getRepoName(t *testing.T) {
	type args struct {
		repo string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with alias",
			args: args{
				repo: "shopee/pl/marketplace-payment",
			},
			want: "bridge",
		},
		{
			name: "no alias",
			args: args{
				repo: "shopee/marketplace-payments/channel",
			},
			want: "channel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepoName(tt.args.repo); got != tt.want {
				t.Errorf("getRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
