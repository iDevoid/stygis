package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
	mock_user "github.com/iDevoid/stygis/mocks/user"
)

func TestUserInit(t *testing.T) {
	type args struct {
		cache       user.Caching
		persistence user.Persistence
	}
	tests := []struct {
		name string
		args args
		want user.Repository
	}{
		{
			name: "success",
			args: args{
				cache:       nil,
				persistence: nil,
			},
			want: &userRepo{
				cache:       nil,
				persistence: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInit(tt.args.cache, tt.args.persistence); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_DataProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name     string
		args     args
		want     *model.User
		wantErr  bool
		initMock func() (user.Caching, user.Persistence)
	}{
		{
			name: "exists in cache",
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
			initMock: func() (user.Caching, user.Persistence) {
				caching := mock_user.NewMockCaching(ctrl)
				caching.EXPECT().Get(gomock.Any(), int64(1)).Return(
					&model.User{
						ID:        1,
						Email:     "clyf@email.com",
						Username:  "clyf",
						CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
						LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
						Status:    1,
					},
					nil,
				)
				return caching, nil
			},
		},
		{
			name: "all errors",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			wantErr: true,
			initMock: func() (user.Caching, user.Persistence) {
				caching := mock_user.NewMockCaching(ctrl)
				caching.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, errors.New("ERROR"))
				db := mock_user.NewMockPersistence(ctrl)
				db.EXPECT().FindByID(gomock.Any(), int64(1)).Return(nil, errors.New("ERROR"))
				return caching, db
			},
		},
		{
			name: "use db to get data and save to caching",
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
			initMock: func() (user.Caching, user.Persistence) {
				caching := mock_user.NewMockCaching(ctrl)
				caching.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, errors.New("ERROR"))
				db := mock_user.NewMockPersistence(ctrl)
				db.EXPECT().FindByID(gomock.Any(), int64(1)).Return(
					&model.User{
						ID:        1,
						Email:     "clyf@email.com",
						Username:  "clyf",
						CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
						LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
						Status:    1,
					},
					nil,
				)
				caching.EXPECT().Save(gomock.Any(), &model.User{
					ID:        1,
					Email:     "clyf@email.com",
					Username:  "clyf",
					CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					Status:    1,
				}).Return(nil)
				return caching, db
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache, mockPersistence := tt.initMock()
			ur := &userRepo{
				cache:       mockCache,
				persistence: mockPersistence,
			}
			got, err := ur.DataProfile(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.DataProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.DataProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
