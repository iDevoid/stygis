package user

import (
	"reflect"
	"testing"
)

func TestInitializeDomain(t *testing.T) {
	type args struct {
		persistence Persistence
		caching     Caching
		repository  Repository
	}
	tests := []struct {
		name string
		args args
		want Usecase
	}{
		{
			name: "success",
			args: args{
				persistence: nil,
				caching:     nil,
				repository:  nil,
			},
			want: &service{
				userPersistence: nil,
				userCaching:     nil,
				userRepository:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeDomain(tt.args.persistence, tt.args.caching, tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
