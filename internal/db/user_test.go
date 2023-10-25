package db

import (
	"context"
	"database/sql"
	"fmt"
	"go-test-grpc-http/internal/entity"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Test_source_GetUserById(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
		id  *entity.UserID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: GetUserById source: user found",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: &entity.UserDB{
				ID:         uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				FirstName:  "John",
				LastName:   "Doe",
				SecondName: "DoeD",
				Age:        30,
				Email:      "doe@example.com",
				Phone:      "+1111111111",
				Password:   "qwerty1234",
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"first_name",
						"last_name",
						"second_name",
						"age",
						"email",
						"phone",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"Doe",
						"DoeD",
						30,
						"doe@example.com",
						"+1111111111",
						"qwerty1234",
					)
				f.db.ExpectQuery("SELECT * FROM users WHERE id = $1").WithArgs(a.id.String()).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "success: GetUserById source: can't find user",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("SELECT * FROM users WHERE id = $1").WithArgs(a.id.String()).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "error: GetUserById source: can't exec query",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").WithArgs(a.id.String()).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: GetUserById source: can't scan user",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("SELECT * FROM users WHERE id = $1").WithArgs(a.id.String()).WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			s := &source{
				db: sqlx.NewDb(db, "sqlmock"),
			}

			tt.setup(tt.args, f)

			got, err := s.GetUserById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("source.GetUserById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_source_CreateUser(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx  context.Context
		user *entity.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: CreateUser source: user created",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"first_name",
						"last_name",
						"second_name",
						"age",
						"email",
						"phone",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"Doe",
						"DoeD",
						30,
						"doe@example.com",
						"+1111111111",
						"qwerty1234",
					)
				f.db.ExpectQuery("INSERT INTO users (id, first_name, last_name, second_name, age, email, phone, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error: CreateUser source: can't exec query",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: CreateUser source: can't scan user",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("INSERT INTO users (id, first_name, last_name, second_name, age, email, phone, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)").
					WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			s := &source{
				db: sqlx.NewDb(db, "sqlmock"),
			}

			tt.setup(tt.args, f)

			got, err := s.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				_, err = uuid.Parse(got.String())
				if err != nil {
					t.Errorf("source.CreateUser() = %v, want uuid", got)
				}
			}
		})
	}
}

func Test_source_UpdateUser(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx  context.Context
		id   *entity.UserID
		user *entity.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: Update source: user updated",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        31,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			want: &entity.UserDB{
				ID:         uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				FirstName:  "John",
				LastName:   "Doe",
				SecondName: "DoeD",
				Age:        31,
				Email:      "doe@example.com",
				Phone:      "+1111111111",
				Password:   "qwerty1234",
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"first_name",
						"last_name",
						"second_name",
						"age",
						"email",
						"phone",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"Doe",
						"DoeD",
						31,
						"doe@example.com",
						"+1111111111",
						"qwerty1234",
					)
				f.db.ExpectQuery("UPDATE users SET first_name = $1, last_name = $2, second_name = $3, age = $4, email = $5, phone = $6, password = $7 WHERE id = $8").
					WithArgs("John", "Doe", "DoeD", 31, "doe@example.com", "+1111111111", "qwerty1234", uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error: Update source: can't exec query",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs("John", "Doe", "DoeD", 31, "doe@example.com", "+1111111111", "qwerty1234", uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: Update source: can't scan user",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("UPDATE users SET first_name = $1, last_name = $2, second_name = $3, age = $4, email = $5, phone = $6, password = $7 WHERE id = $8").
					WithArgs("John", "Doe", "DoeD", 31, "doe@example.com", "+1111111111", "qwerty1234", uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			s := &source{
				db: sqlx.NewDb(db, "sqlmock"),
			}

			tt.setup(tt.args, f)

			got, err := s.UpdateUser(tt.args.ctx, tt.args.id, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("source.UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_source_DeleteUser(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
		id  *entity.UserID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: Delete source: user deleted",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("DELETE FROM users WHERE id = $1").
					WithArgs(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			wantErr: false,
		},
		{
			name: "error: DeleteUser source: can't exec query",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: UpDEleteUser source: can't scan user",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("DELETE FROM users WHERE id = ?").
					WithArgs(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			s := &source{
				db: sqlx.NewDb(db, "sqlmock"),
			}

			tt.setup(tt.args, f)

			if err := s.DeleteUser(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("source.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
