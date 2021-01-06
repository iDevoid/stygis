package rest

import (
	"reflect"
	"testing"

	"github.com/iDevoid/stygis/internal/module/user"
)

func TestUserInit(t *testing.T) {
	type args struct {
		userCase user.Usecase
	}
	tests := []struct {
		name string
		args args
		want UserHandler
	}{
		{
			name: "success",
			args: args{
				userCase: nil,
			},
			want: &userHandler{
				userCase: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInit(tt.args.userCase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}
