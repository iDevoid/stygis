package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/stygis/internal/constants/model"
	mock_user "github.com/iDevoid/stygis/mocks/user"
)

func Test_service_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		initMock func() (Persistence, Caching)
	}{
		{
			name: "error persistence create",
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				},
			},
			wantErr: true,
			initMock: func() (Persistence, Caching) {
				mockedPersis := mock_user.NewMockPersistence(ctrl)
				mockedPersis.EXPECT().Create(gomock.Any(), &model.User{
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				}).Return(nil, errors.New("ERROR"))
				return mockedPersis, nil
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				},
			},
			initMock: func() (Persistence, Caching) {
				mockedPersis := mock_user.NewMockPersistence(ctrl)
				mockedPersis.EXPECT().Create(gomock.Any(), &model.User{
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				}).Return(&model.User{
					ID:       1,
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				}, nil)
				mockedCache := mock_user.NewMockCaching(ctrl)
				mockedCache.EXPECT().Save(gomock.Any(), &model.User{
					ID:       1,
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "hashedpassword",
				}).Return(errors.New("ERROR"))
				return mockedPersis, mockedCache
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, c := tt.initMock()
			s := &service{
				userPersistence: p,
				userCaching:     c,
			}
			if err := s.Register(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("service.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
