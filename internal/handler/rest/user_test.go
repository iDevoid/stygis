package rest

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
	mock_user "github.com/iDevoid/stygis/mocks/user"
	"github.com/savsgio/atreugo/v10"
	"github.com/valyala/fasthttp"
)

func TestHandleUser(t *testing.T) {
	type args struct {
		usecase user.Usecase
		hash    hash.Hash
	}
	tests := []struct {
		name string
		args args
		want user.Handler
	}{
		{
			name: "success",
			args: args{
				usecase: nil,
				hash:    nil,
			},
			want: &userService{
				usecase: nil,
				hash:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleUser(tt.args.usecase, tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test(t *testing.T) {
	type args struct {
		ctx *atreugo.RequestCtx
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx: &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				},
			},
			want:       "Hello World!",
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userService{}
			if err := us.Test(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("userService.Test() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := string(tt.args.ctx.Response.Body())
			if got != tt.want {
				t.Errorf("Response.Body() = %v, want %v", got, tt.want)
			}
			gotStatus := tt.args.ctx.Response.StatusCode()
			if gotStatus != tt.wantStatus {
				t.Errorf("Response.StatusCode() = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_userService_CreateNewAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name       string
		want       string
		wantStatus int
		keyCookie  string
		wantCookie string
		wantErr    bool
		initMock   func() (*atreugo.RequestCtx, user.Usecase)
	}{
		{
			name:       "bad username",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				ctx.Request.PostArgs().Add("password", "nothashedpassword")
				return ctx, nil
			},
		},
		{
			name:       "bad email",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("username", "clyf")
				ctx.Request.PostArgs().Add("password", "nothashedpassword")
				return ctx, nil
			},
		},
		{
			name:       "bad password",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("username", "clyf")
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				return ctx, nil
			},
		},
		{
			name:       "error usecase register",
			wantStatus: 500,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("username", "clyf")
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				ctx.Request.PostArgs().Add("password", "nothashedpassword")

				mockedUsecase := mock_user.NewMockUsecase(ctrl)
				mockedUsecase.EXPECT().Register(gomock.Any(), &model.User{
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "ae7bc8bd67f4f5c6a911373677ec56f95246288f1130c62048703be38397bda7",
				}).Return(errors.New("ERROR"))
				return ctx, mockedUsecase
			},
		},
		{
			name:       "success",
			wantStatus: 200,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("username", "clyf")
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				ctx.Request.PostArgs().Add("password", "nothashedpassword")

				mockedUsecase := mock_user.NewMockUsecase(ctrl)
				mockedUsecase.EXPECT().Register(gomock.Any(), &model.User{
					Username: "clyf",
					Email:    "clyf@email.com",
					Password: "ae7bc8bd67f4f5c6a911373677ec56f95246288f1130c62048703be38397bda7",
				}).Return(nil)
				return ctx, mockedUsecase
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, usecase := tt.initMock()
			us := &userService{
				hash:    sha256.New(),
				usecase: usecase,
			}
			if err := us.CreateNewAccount(ctx); (err != nil) != tt.wantErr {
				t.Errorf("userService.CreateNewAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := string(ctx.Response.Body())
			if got != tt.want {
				t.Errorf("Response.Body() = %v, want %v", got, tt.want)
			}
			gotStatus := ctx.Response.StatusCode()
			if gotStatus != tt.wantStatus {
				t.Errorf("Response.StatusCode() = %v, want %v", gotStatus, tt.wantStatus)
			}
			if tt.keyCookie == "" {
				return
			}
			gotCookie := string(ctx.Response.Header.PeekCookie(tt.keyCookie))
			if !strings.Contains(got, fmt.Sprint("=", tt.wantCookie)) && tt.wantCookie != "" {
				t.Errorf("Response.Body() = %v, want %v", gotCookie, tt.wantCookie)
			}
		})
	}
}

func Test_userService_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name       string
		want       string
		wantStatus int
		keyCookie  string
		wantCookie string
		wantErr    bool
		initMock   func() (*atreugo.RequestCtx, user.Usecase)
	}{

		{
			name:       "bad email",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("password", "nothashedpassword")
				return ctx, nil
			},
		},
		{
			name:       "bad password",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				return ctx, nil
			},
		},
		{
			name:       "error usecase login",
			wantStatus: 500,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("username", "clyf")
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				ctx.Request.PostArgs().Add("password", "nothashedpassword")

				mockedUsecase := mock_user.NewMockUsecase(ctrl)
				mockedUsecase.EXPECT().Login(gomock.Any(), "clyf@email.com", "ae7bc8bd67f4f5c6a911373677ec56f95246288f1130c62048703be38397bda7").Return(int64(0), errors.New("ERROR"))
				return ctx, mockedUsecase
			},
		},
		{
			name:       "success",
			wantStatus: 200,
			keyCookie:  "unique_key_value",
			wantCookie: "1",
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.PostArgs().Add("username", "clyf")
				ctx.Request.PostArgs().Add("email", "clyf@email.com")
				ctx.Request.PostArgs().Add("password", "nothashedpassword")

				mockedUsecase := mock_user.NewMockUsecase(ctrl)
				mockedUsecase.EXPECT().Login(gomock.Any(), "clyf@email.com", "ae7bc8bd67f4f5c6a911373677ec56f95246288f1130c62048703be38397bda7").Return(int64(1), nil)
				return ctx, mockedUsecase
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, usecase := tt.initMock()
			us := &userService{
				hash:    sha256.New(),
				usecase: usecase,
			}
			if err := us.SignIn(ctx); (err != nil) != tt.wantErr {
				t.Errorf("userService.SignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := string(ctx.Response.Body())
			if got != tt.want {
				t.Errorf("Response.Body() = %v, want %v", got, tt.want)
			}
			gotStatus := ctx.Response.StatusCode()
			if gotStatus != tt.wantStatus {
				t.Errorf("Response.StatusCode() = %v, want %v", gotStatus, tt.wantStatus)

			}
			if tt.keyCookie == "" {
				return
			}
			gotCookie := string(ctx.Response.Header.PeekCookie(tt.keyCookie))
			cookieVals := strings.Split(gotCookie, "; ")
			if !strings.Contains(cookieVals[0], fmt.Sprint("=", tt.wantCookie)) && tt.wantCookie != "" {
				t.Errorf("Cookie = %v, want %v", gotCookie, tt.wantCookie)
			}
		})
	}
}

func Test_userService_ShowProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name       string
		want       string
		wantStatus int
		keyCookie  string
		wantCookie string
		wantErr    bool
		initMock   func() (*atreugo.RequestCtx, user.Usecase)
	}{
		{
			name:       "error no cookie",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				return ctx, nil
			},
		},
		{
			name:       "error bad cookie",
			wantStatus: 400,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.Header.SetCookie("unique_key_value", "asd")
				return ctx, nil
			},
		},
		{
			name:       "error usecase profile",
			wantStatus: 500,
			wantErr:    true,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.Header.SetCookie("unique_key_value", "1")
				mockedUsecase := mock_user.NewMockUsecase(ctrl)
				mockedUsecase.EXPECT().Profile(gomock.Any(), int64(1)).Return(nil, errors.New("ERROR"))
				return ctx, mockedUsecase
			},
		},
		{
			name:       "success",
			want:       `{"user_id":1,"username":"clyf","email":"clyf@email.com","created_at":"2020-02-20T20:20:20.00000002Z","last_login":"2020-02-20T20:20:20.00000002Z","status":1}`,
			wantStatus: 200,
			initMock: func() (*atreugo.RequestCtx, user.Usecase) {
				ctx := &atreugo.RequestCtx{
					RequestCtx: &fasthttp.RequestCtx{},
				}
				ctx.Request.Header.SetCookie("unique_key_value", "1")
				mockedUsecase := mock_user.NewMockUsecase(ctrl)
				mockedUsecase.EXPECT().Profile(gomock.Any(), int64(1)).Return(&model.User{
					ID:        1,
					Username:  "clyf",
					Email:     "clyf@email.com",
					CreatedAt: time.Date(2020, time.February, 20, 20, 20, 20, 20, time.UTC),
					LastLogin: time.Date(2020, time.February, 20, 20, 20, 20, 20, time.UTC),
					Status:    1,
				}, nil)
				return ctx, mockedUsecase
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, usecase := tt.initMock()
			us := &userService{
				hash:    sha256.New(),
				usecase: usecase,
			}
			if err := us.ShowProfile(ctx); (err != nil) != tt.wantErr {
				t.Errorf("userService.ShowProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := string(ctx.Response.Body())
			if got != tt.want {
				t.Errorf("Response.Body() = %v, want %v", got, tt.want)
			}
			gotStatus := ctx.Response.StatusCode()
			if gotStatus != tt.wantStatus {
				t.Errorf("Response.StatusCode() = %v, want %v", gotStatus, tt.wantStatus)

			}
			if tt.keyCookie == "" {
				return
			}
			gotCookie := string(ctx.Response.Header.PeekCookie(tt.keyCookie))
			cookieVals := strings.Split(gotCookie, "; ")
			if !strings.Contains(cookieVals[0], fmt.Sprint("=", tt.wantCookie)) && tt.wantCookie != "" {
				t.Errorf("Cookie = %v, want %v", gotCookie, tt.wantCookie)
			}
		})
	}
}
