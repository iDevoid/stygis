package persistence

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/iDevoid/cptx"
	"github.com/iDevoid/stygis/internal/constant/model"
	"github.com/iDevoid/stygis/internal/constant/query"
)

func TestProfileInit(t *testing.T) {
	type args struct {
		db cptx.Database
	}
	tests := []struct {
		name string
		args args
		want ProfilePersistence
	}{
		{
			name: "success",
			args: args{
				db: nil,
			},
			want: &profilePersistence{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProfileInit(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProfileInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_profilePersistence_InsertProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		db cptx.Database
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
			name: "error",
			fields: fields{
				db: func() cptx.Database {
					_mockedMain := cptx.NewMockMainDB(ctrl)
					_mockedMain.EXPECT().ExecuteMustTx(gomock.Any(), query.ProfileInsert,
						map[string]interface{}{
							"id":          int64(1),
							"username":    "clyf",
							"full_name":   "clyf",
							"status":      0,
							"create_time": time.Date(2020, 12, 21, 12, 12, 12, 0, time.UTC),
						},
					).Return(nil, errors.New("ERROR"))
					mocked := cptx.NewMockDatabase(ctrl)
					mocked.EXPECT().Main().Return(_mockedMain)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					ID:         1,
					Username:   "clyf",
					CreateTime: time.Date(2020, 12, 21, 12, 12, 12, 0, time.UTC),
				},
			},
			wantErr: true,
		}, {
			name: "success",
			fields: fields{
				db: func() cptx.Database {
					_mockedMain := cptx.NewMockMainDB(ctrl)
					_mockedMain.EXPECT().ExecuteMustTx(gomock.Any(), query.ProfileInsert,
						map[string]interface{}{
							"id":          int64(1),
							"username":    "clyf",
							"full_name":   "clyf",
							"status":      0,
							"create_time": time.Date(2020, 12, 21, 12, 12, 12, 0, time.UTC),
						},
					).Return(nil, nil)
					mocked := cptx.NewMockDatabase(ctrl)
					mocked.EXPECT().Main().Return(_mockedMain)
					return mocked
				}(),
			},
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					ID:         1,
					Username:   "clyf",
					CreateTime: time.Date(2020, 12, 21, 12, 12, 12, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp := &profilePersistence{
				db: tt.fields.db,
			}
			if err := pp.InsertProfile(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("profilePersistence.InsertProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
