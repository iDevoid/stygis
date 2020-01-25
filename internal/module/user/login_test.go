package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/stygis/internal/constants/model"
	mock_user "github.com/iDevoid/stygis/mocks/user"
)

func Test_service_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name     string
		args     args
		want     int64
		wantErr  bool
		initMock func() (Persistence, Caching)
	}{
		{
			name: "error persistence find",
			args: args{
				ctx:      context.TODO(),
				email:    "clyf@email.com",
				password: "hashedpassword",
			},
			wantErr: true,
			initMock: func() (Persistence, Caching) {
				mockedPersis := mock_user.NewMockPersistence(ctrl)
				mockedPersis.EXPECT().Find(gomock.Any(), "clyf@email.com", "hashedpassword").Return(nil, errors.New("ERROR"))
				return mockedPersis, nil
			},
		},
		{
			name: "error caching save",
			args: args{
				ctx:      context.TODO(),
				email:    "clyf@email.com",
				password: "hashedpassword",
			},
			want:    1,
			wantErr: true,
			initMock: func() (Persistence, Caching) {
				mockedPersis := mock_user.NewMockPersistence(ctrl)
				mockedPersis.EXPECT().Find(gomock.Any(), "clyf@email.com", "hashedpassword").Return(&model.User{
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
		{
			name: "success",
			args: args{
				ctx:      context.TODO(),
				email:    "clyf@email.com",
				password: "hashedpassword",
			},
			want: 1,
			initMock: func() (Persistence, Caching) {
				mockedPersis := mock_user.NewMockPersistence(ctrl)
				mockedPersis.EXPECT().Find(gomock.Any(), "clyf@email.com", "hashedpassword").Return(&model.User{
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
				}).Return(nil)
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
			got, err := s.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
