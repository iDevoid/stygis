package persistence

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iDevoid/stygis/internal/constants/model"
	"github.com/iDevoid/stygis/internal/module/user"
	"github.com/iDevoid/stygis/mocks"
	"github.com/iDevoid/stygis/platform/postgres"
)

func TestUserInit(t *testing.T) {
	type args struct {
		db *postgres.Database
	}
	tests := []struct {
		name string
		args args
		want user.Persistence
	}{
		{
			name: "success",
			args: args{
				db: nil,
			},
			want: &userPersistence{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserInit(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userPersistence_Create(t *testing.T) {
	sqlxDB, mocked, db := mocks.PSQLMock()
	defer db.Close()

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		initMock func() *postgres.Database
	}{
		{
			name: "error beginx",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
			initMock: func() *postgres.Database {
				return &postgres.Database{
					Master: sqlxDB,
				}
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
			initMock: func() *postgres.Database {
				mocked.ExpectBegin()
				mocked.ExpectQuery(`
					INSERT INTO account \(
						email, 
						hash_password, 
						username, 
						created_at, 
						last_login, 
						status
					\) VALUES \(
						(.+),
						(.+),
						(.+),
						NOW\(\),
						NOW\(\),
						(.+)
					\)
					RETURNING id
				`).WillReturnRows(sqlmock.NewRows([]string{
					"id",
				}).AddRow(
					1,
				))
				mocked.ExpectCommit()

				return &postgres.Database{
					Master: sqlxDB,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := &userPersistence{
				db: tt.initMock(),
			}
			if err := up.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userPersistence.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userPersistence_FindByID(t *testing.T) {
	sqlxDB, mocked, db := mocks.PSQLMock()
	defer db.Close()

	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name     string
		args     args
		wantUser *model.User
		wantErr  bool
		initMock func() *postgres.Database
	}{
		{
			name: "sql no rows",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			wantUser: &model.User{},
			initMock: func() *postgres.Database {
				mocked.ExpectQuery(`
					SELECT 
						id, 
						username,  
						last_login, 
						status 
					FROM 
						account
					WHERE
						id = (.+) AND
						statuts = (.+)
				`).WillReturnError(sql.ErrNoRows)

				return &postgres.Database{
					Slave: sqlxDB,
				}
			},
		},
		{
			name: "success",
			args: args{
				ctx:    context.TODO(),
				userID: 1,
			},
			wantUser: &model.User{
				ID:        1,
				Username:  "clyf",
				LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
				Status:    1,
			},
			initMock: func() *postgres.Database {
				mocked.ExpectQuery(`
					SELECT 
						id, 
						username,  
						last_login, 
						status 
					FROM 
						account
					WHERE
						id = (.+) AND
						statuts = (.+)
				`).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"username",
					"last_login",
					"status",
				}).AddRow(
					1,
					"clyf",
					time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					1,
				))

				return &postgres.Database{
					Slave: sqlxDB,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := &userPersistence{
				db: tt.initMock(),
			}
			gotUser, err := up.FindByID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("userPersistence.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("userPersistence.FindByID() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func Test_userPersistence_Find(t *testing.T) {
	sqlxDB, mocked, db := mocks.PSQLMock()
	defer db.Close()

	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name     string
		args     args
		want     *model.User
		wantErr  bool
		initMock func() *postgres.Database
	}{
		{
			name: "success",
			args: args{
				ctx:      context.TODO(),
				email:    "clyf@email.com",
				password: "hashedpassword",
			},
			want: &model.User{},
			initMock: func() *postgres.Database {
				mocked.ExpectQuery(`
					SELECT 
						id, 
						email, 
						username, 
						created_at, 
						last_login, 
						status 
					FROM 
						account
					WHERE
						email = (.+) AND
						hash_password = (.+) AND
						statuts = (.+)
				`).WillReturnError(sql.ErrNoRows)

				return &postgres.Database{
					Slave: sqlxDB,
				}
			},
		},
		{
			name: "success",
			args: args{
				ctx:      context.TODO(),
				email:    "clyf@email.com",
				password: "hashedpassword",
			},
			want: &model.User{
				ID:        1,
				Email:     "clyf@email.com",
				Username:  "clyf",
				CreatedAt: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
				LastLogin: time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
				Status:    1,
			},
			initMock: func() *postgres.Database {
				mocked.ExpectQuery(`
					SELECT 
						id, 
						email, 
						username, 
						created_at, 
						last_login, 
						status 
					FROM 
						account
					WHERE
						email = (.+) AND
						hash_password = (.+) AND
						statuts = (.+)
				`).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"email",
					"username",
					"created_at",
					"last_login",
					"status",
				}).AddRow(
					1,
					"clyf@email.com",
					"clyf",
					time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					time.Date(2020, time.January, 12, 12, 12, 12, 12, time.UTC),
					1,
				))

				return &postgres.Database{
					Slave: sqlxDB,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := &userPersistence{
				db: tt.initMock(),
			}
			got, err := up.Find(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userPersistence.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userPersistence.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userPersistence_ChangePassword(t *testing.T) {
	sqlxDB, mocked, db := mocks.PSQLMock()
	defer db.Close()

	type args struct {
		ctx         context.Context
		newPassword string
		user        *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		initMock func() *postgres.Database
	}{
		{
			name: "error beginx",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
			initMock: func() *postgres.Database {
				return &postgres.Database{
					Master: sqlxDB,
				}
			},
		},
		{
			name: "success",
			args: args{
				ctx:         context.TODO(),
				newPassword: "newhashedpassword",
				user: &model.User{
					ID:       1,
					Email:    "clyf@email.com",
					Password: "oldhashedpassword",
				},
			},
			initMock: func() *postgres.Database {
				mocked.ExpectBegin()
				mocked.ExpectExec(`
					UPDATE 
						account 
					SET  
						hash_password = (.+)
					WHERE 
						id = (.+) AND
						email = (.+) AND
						hash_password = (.+) AND
						statuts = (.+)
				`).WillReturnResult(sqlmock.NewResult(1, 1))
				mocked.ExpectCommit()

				return &postgres.Database{
					Master: sqlxDB,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := &userPersistence{
				db: tt.initMock(),
			}
			if err := up.ChangePassword(tt.args.ctx, tt.args.newPassword, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userPersistence.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userPersistence_Delete(t *testing.T) {
	sqlxDB, mocked, db := mocks.PSQLMock()
	defer db.Close()

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		initMock func() *postgres.Database
	}{
		{
			name: "error beginx",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
			initMock: func() *postgres.Database {
				return &postgres.Database{
					Master: sqlxDB,
				}
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				user: &model.User{
					ID:       1,
					Email:    "clyf@email.com",
					Password: "oldhashedpassword",
				},
			},
			initMock: func() *postgres.Database {
				mocked.ExpectBegin()
				mocked.ExpectExec(`
					UPDATE 
						account 
					SET  
						status = (.+)
					WHERE 
						id = (.+) AND
						email = (.+) AND
						hash_password = (.+) AND
						statuts = (.+)
				`).WillReturnResult(sqlmock.NewResult(1, 1))
				mocked.ExpectCommit()

				return &postgres.Database{
					Master: sqlxDB,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := &userPersistence{
				db: tt.initMock(),
			}
			if err := up.Delete(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userPersistence.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
