package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/cptx"
	"github.com/iDevoid/stygis/internal/constant/model"
	"github.com/iDevoid/stygis/internal/repository"
	"github.com/iDevoid/stygis/internal/storage/persistence"
	profile_mock "github.com/iDevoid/stygis/mocks/profile"
	user_mock "github.com/iDevoid/stygis/mocks/user"
)

func Test_service_Registration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		transaction    cptx.Transaction
		userRepo       repository.UserRepository
		userPersist    persistence.UserPersistence
		profilePersist persistence.ProfilePersistence
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "bad email",
			fields: fields{},
			args: args{
				user: &model.User{},
			},
			wantErr: true,
		}, {
			name: "error encrypt",
			fields: fields{
				userRepo: func() repository.UserRepository {
					mocked := user_mock.NewMockUserRepository(ctrl)
					mocked.EXPECT().Encrypt(gomock.Any()).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				user: &model.User{
					Email: "clyf@example.com",
				},
			},
			wantErr: true,
		}, {
			name: "error transaction",
			fields: fields{
				userRepo: func() repository.UserRepository {
					mocked := user_mock.NewMockUserRepository(ctrl)
					mocked.EXPECT().Encrypt(gomock.Any()).Return(nil)
					return mocked
				}(),
				transaction: func() cptx.Transaction {
					mocked := cptx.NewMockTransaction(ctrl)
					mocked.EXPECT().Begin(gomock.Any()).Return(nil, errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Email: "clyf@example.com",
				},
			},
			wantErr: true,
		}, {
			name: "error insert user",
			fields: fields{
				userRepo: func() repository.UserRepository {
					mocked := user_mock.NewMockUserRepository(ctrl)
					mocked.EXPECT().Encrypt(gomock.Any()).Return(nil)
					return mocked
				}(),
				transaction: func() cptx.Transaction {
					mocked := cptx.NewMockTransaction(ctrl)
					_tx := cptx.NewMockTx(ctrl)
					_tx.EXPECT().Rollback()
					mocked.EXPECT().Begin(gomock.Any()).Return(_tx, nil)
					return mocked
				}(),
				userPersist: func() persistence.UserPersistence {
					mocked := user_mock.NewMockUserPersistence(ctrl)
					mocked.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Email: "clyf@example.com",
				},
			},
			wantErr: true,
		}, {
			name: "error insert profile",
			fields: fields{
				userRepo: func() repository.UserRepository {
					mocked := user_mock.NewMockUserRepository(ctrl)
					mocked.EXPECT().Encrypt(gomock.Any()).Return(nil)
					return mocked
				}(),
				transaction: func() cptx.Transaction {
					mocked := cptx.NewMockTransaction(ctrl)
					_tx := cptx.NewMockTx(ctrl)
					_tx.EXPECT().Rollback()
					mocked.EXPECT().Begin(gomock.Any()).Return(_tx, nil)
					return mocked
				}(),
				userPersist: func() persistence.UserPersistence {
					mocked := user_mock.NewMockUserPersistence(ctrl)
					mocked.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
				profilePersist: func() persistence.ProfilePersistence {
					mocked := profile_mock.NewMockProfilePersistence(ctrl)
					mocked.EXPECT().InsertProfile(gomock.Any(), gomock.Any()).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Email: "clyf@example.com",
				},
			},
			wantErr: true,
		}, {
			name: "error commit",
			fields: fields{
				userRepo: func() repository.UserRepository {
					mocked := user_mock.NewMockUserRepository(ctrl)
					mocked.EXPECT().Encrypt(gomock.Any()).Return(nil)
					return mocked
				}(),
				transaction: func() cptx.Transaction {
					mocked := cptx.NewMockTransaction(ctrl)
					_tx := cptx.NewMockTx(ctrl)
					_tx.EXPECT().Rollback()
					_tx.EXPECT().Commit().Return(errors.New("ERROR"))
					mocked.EXPECT().Begin(gomock.Any()).Return(_tx, nil)
					return mocked
				}(),
				userPersist: func() persistence.UserPersistence {
					mocked := user_mock.NewMockUserPersistence(ctrl)
					mocked.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
				profilePersist: func() persistence.ProfilePersistence {
					mocked := profile_mock.NewMockProfilePersistence(ctrl)
					mocked.EXPECT().InsertProfile(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Email: "clyf@example.com",
				},
			},
			wantErr: true,
		}, {
			name: "success",
			fields: fields{
				userRepo: func() repository.UserRepository {
					mocked := user_mock.NewMockUserRepository(ctrl)
					mocked.EXPECT().Encrypt(gomock.Any()).Return(nil)
					return mocked
				}(),
				transaction: func() cptx.Transaction {
					mocked := cptx.NewMockTransaction(ctrl)
					_tx := cptx.NewMockTx(ctrl)
					_tx.EXPECT().Rollback()
					_tx.EXPECT().Commit().Return(nil)
					mocked.EXPECT().Begin(gomock.Any()).Return(_tx, nil)
					return mocked
				}(),
				userPersist: func() persistence.UserPersistence {
					mocked := user_mock.NewMockUserPersistence(ctrl)
					mocked.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
				profilePersist: func() persistence.ProfilePersistence {
					mocked := profile_mock.NewMockProfilePersistence(ctrl)
					mocked.EXPECT().InsertProfile(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					Email: "clyf@example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				transaction:    tt.fields.transaction,
				userRepo:       tt.fields.userRepo,
				userPersist:    tt.fields.userPersist,
				profilePersist: tt.fields.profilePersist,
			}
			if err := s.Registration(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("service.Registration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
