package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/stygis/internal/constants/model"
	mock_user "github.com/iDevoid/stygis/mocks/user"
)

func Test_service_Profile(t *testing.T) {
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
		initMock func() Repository
	}{
		{
			name: "error repository",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			wantErr: true,
			initMock: func() Repository {
				mocked := mock_user.NewMockRepository(ctrl)
				mocked.EXPECT().DataProfile(gomock.Any(), int64(1)).Return(nil, errors.New("ERROR"))
				return mocked
			},
		},
		{
			name: "success",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			want: &model.User{
				ID:       1,
				Username: "clyf",
				Email:    "clyf@email.com",
				Password: "hashedpassword",
			},
			initMock: func() Repository {
				mocked := mock_user.NewMockRepository(ctrl)
				mocked.EXPECT().DataProfile(gomock.Any(), int64(1)).Return(&model.User{
					ID:       1,
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				}, nil)
				return mocked
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				userRepository: tt.initMock(),
			}
			got, err := s.Profile(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Profile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Profile() = %v, want %v", got, tt.want)
			}
		})
	}
}
