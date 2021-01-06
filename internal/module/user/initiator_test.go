package user

import (
	"reflect"
	"testing"

	"github.com/iDevoid/cptx"
	"github.com/iDevoid/stygis/internal/repository"
	"github.com/iDevoid/stygis/internal/storage/persistence"
)

func TestInitialize(t *testing.T) {
	type args struct {
		transaction    cptx.Transaction
		userRepo       repository.UserRepository
		userPersist    persistence.UserPersistence
		profilePersist persistence.ProfilePersistence
	}
	tests := []struct {
		name string
		args args
		want Usecase
	}{
		{
			name: "success",
			args: args{
				transaction:    nil,
				userRepo:       nil,
				userPersist:    nil,
				profilePersist: nil,
			},
			want: &service{
				transaction:    nil,
				userRepo:       nil,
				userPersist:    nil,
				profilePersist: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Initialize(tt.args.transaction, tt.args.userRepo, tt.args.userPersist, tt.args.profilePersist); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
