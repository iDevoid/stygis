package repository

import (
	"reflect"
	"testing"

	"github.com/iDevoid/stygis/internal/constant/model"
)

func TestUserInit(t *testing.T) {
	type args struct {
		systemEncryptKey string
	}
	tests := []struct {
		name string
		args args
		want UserRepository
	}{
		{
			name: "success",
			args: args{
				systemEncryptKey: "24cacf5004bf68ae9daad19a5bba391d85ad1cb0b31366e89aec86fad0ab16cb",
			},
			want: &userRepository{
				systemEncryptKey: "24cacf5004bf68ae9daad19a5bba391d85ad1cb0b31366e89aec86fad0ab16cb",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInit(tt.args.systemEncryptKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_Encrypt(t *testing.T) {
	type fields struct {
		systemEncryptKey string
	}
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				systemEncryptKey: "asd",
			},
			args: args{
				user: &model.User{
					Email:    "clyf@example.com",
					Password: "password",
				},
			},
			wantErr: true,
		}, {
			name: "success",
			fields: fields{
				systemEncryptKey: "24cacf5004bf68ae9daad19a5bba391d85ad1cb0b31366e89aec86fad0ab16cb",
			},
			args: args{
				user: &model.User{
					Email:    "clyf@example.com",
					Password: "password",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &userRepository{
				systemEncryptKey: tt.fields.systemEncryptKey,
			}
			if err := ur.Encrypt(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
