package cache

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/mocks"
)

func TestUserInit(t *testing.T) {
	type args struct {
		conn *redis.Client
	}
	tests := []struct {
		name string
		args args
		want user.Caching
	}{
		{
			name: "success",
			args: args{
				conn: nil,
			},
			want: &userCache{
				connection: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInit(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userCache_Save(t *testing.T) {
	client, miniredis := mocks.RedisMock()

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		expecting func()
	}{
		{
			name: "redis exists",
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					ID:        1,
					Email:     "clyf@email.com",
					Username:  "clyf",
					CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					Status:    1,
				},
			},
			expecting: func() {
				miniredis.CheckGet(t, "user:1", `{"user_id":1,"username":"clyf","email":"clyf@email.com","created_at":"2020-01-12T12:12:12.000000012Z","last_login":"2020-01-12T12:12:12.000000012Z","status":1}`)
			},
		},
		{
			name: "redis expire",
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					ID:        1,
					Email:     "clyf@email.com",
					Username:  "clyf",
					CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					Status:    1,
				},
			},
			expecting: func() {
				miniredis.FastForward(time.Second * 60 * 60 * 25)
				if miniredis.Exists("user:1") {
					t.Error("This should not be existed anymore")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &userCache{
				connection: client,
			}
			if err := uc.Save(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.expecting()
		})
	}
	miniredis.Close()
}

func Test_userCache_Get(t *testing.T) {
	client, miniredis := mocks.RedisMock()

	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name     string
		args     args
		want     *model.User
		wantErr  bool
		initMock func()
	}{
		{
			name: "not exist",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			wantErr:  true,
			initMock: func() {},
		},
		{
			name: "exists",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			want: &model.User{
				ID:        1,
				Email:     "clyf@email.com",
				Username:  "clyf",
				CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
				LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
				Status:    1,
			},
			initMock: func() {
				miniredis.Set("user:1", `{"user_id":1,"username":"clyf","email":"clyf@email.com","created_at":"2020-01-12T12:12:12.000000012Z","last_login":"2020-01-12T12:12:12.000000012Z","status":1}`)
				miniredis.SetTTL("user:1", time.Duration(time.Second*60*60*24))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initMock()
			uc := &userCache{
				connection: client,
			}
			got, err := uc.Get(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userCache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
	miniredis.Close()
}

func Test_userCache_Delete(t *testing.T) {
	client, miniredis := mocks.RedisMock()

	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			wantErr: miniredis.Exists("user:1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &userCache{
				connection: client,
			}
			if err := uc.Delete(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("userCache.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	miniredis.Close()
}
