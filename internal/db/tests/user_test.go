package db

import (
	"context"
	"database/sql"
	"fmt"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

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
					Username: "John",
					Password: "qwerty1234",
				},
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"username",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"qwerty1234",
					)
				f.db.ExpectQuery("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error: CreateUser source: can't exec query",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "qwerty1234",
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
					Username: "John",
					Password: "qwerty1234",
				},
			},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)").
					WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := usersSource.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != uuid.Nil {
				_, err = uuid.Parse(got.String())
				if err != nil {
					t.Errorf("source.CreateUser() = %v, want uuid", got)
				}
			}
		})
	}
}

func Test_source_GetUserById(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
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
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: &entity.UserDB{
				ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Username: "John",
				Password: "qwerty1234",
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"username",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
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
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
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
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
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
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
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
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := usersSource.GetUserById(tt.args.ctx, tt.args.id)
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

func Test_source_GetUserByUsername(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: GetUserByUsername source: user found",
			args: args{
				ctx:      context.Background(),
				username: "John",
			},
			want: &entity.UserDB{
				ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Username: "John",
				Password: "qwerty1234",
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"username",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"qwerty1234",
					)
				f.db.ExpectQuery("SELECT * FROM users WHERE username = $1").WithArgs(a.username).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "success: GetUserByUsername source: can't find user",
			args: args{
				ctx:      context.Background(),
				username: "John",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("SELECT * FROM users WHERE username = $1").WithArgs(a.username).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "error: GetUserByUsername source: can't exec query",
			args: args{
				ctx:      context.Background(),
				username: "John",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").WithArgs(a.username).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: GetUserByUsername source: can't scan user",
			args: args{
				ctx:      context.Background(),
				username: "John",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("SELECT * FROM users WHERE username = $1").WithArgs(a.username).WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := usersSource.GetUserByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("source.GetUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_source_UpdateUser(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx    context.Context
		userDB *entity.UserDB
		user   *entity.UserCreate
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
				userDB: &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "Bob",
					Password: "qwerty1234",
				},
				user: &entity.UserCreate{
					Username: "John",
					Password: "qwerty1234",
				},
			},
			want: &entity.UserDB{
				ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Username: "John",
				Password: "qwerty1234",
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"username",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"qwerty1234",
					)
				f.db.ExpectQuery("UPDATE users SET username = $1, password = $2 WHERE id = $3").
					WithArgs("John", "qwerty1234", uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error: Update source: can't exec query",
			args: args{
				ctx: context.Background(),
				userDB: &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "Bob",
					Password: "qwerty1234",
				},
				user: &entity.UserCreate{
					Username: "John",
					Password: "qwerty1234",
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs("John", "qwerty1234", uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: Update source: can't scan user",
			args: args{
				ctx: context.Background(),
				userDB: &entity.UserDB{
					ID:       uuid.Nil,
					Username: "Bob",
					Password: "qwerty1234",
				},
				user: &entity.UserCreate{
					Username: "John",
					Password: "qwerty1234",
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("UPDATE users SET username = $1, password = $2 WHERE id = $3").
					WithArgs("John", "qwerty1234", uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't scan user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := usersSource.UpdateUser(tt.args.ctx, tt.args.userDB, tt.args.user)
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
		id  uuid.UUID
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
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: nil,
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"username",
						"password",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"John",
						"qwerty1234",
					)

				f.db.ExpectQuery("DELETE FROM users WHERE id = $1").
					WithArgs(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error: DeleteUser source: can't exec query",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")).
					WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			if err := usersSource.DeleteUser(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("source.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_source_LikeTrack(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx     context.Context
		userId  uuid.UUID
		trackId uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    error
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: LikeTrack source: user likes track",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("INSERT INTO user_music (user_id, music_id) VALUES ($1, $2)").
					WithArgs(
						a.userId,
						a.trackId,
					).WillReturnRows(sqlmock.NewRows([]string{}))
			},
			wantErr: false,
		},
		{
			name: "error: LikeTrack source: can't exec query",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs(
						a.userId,
						a.trackId,
					).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			if err := usersSource.LikeTrack(tt.args.ctx, tt.args.userId, tt.args.trackId); (err != nil) != tt.wantErr {
				t.Errorf("source.LikeTrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_source_DislikeTrack(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx     context.Context
		userId  uuid.UUID
		trackId uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    error
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: DislikeTrack source: user dislikes track",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("DELETE FROM user_music WHERE user_id = $1 AND music_id = $2").
					WithArgs(
						a.userId,
						a.trackId,
					).WillReturnRows(sqlmock.NewRows([]string{}))
			},
			wantErr: false,
		},
		{
			name: "error: DislikeTrack source: can't exec query",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs(
						a.userId,
						a.trackId,
					).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			if err := usersSource.DislikeTrack(tt.args.ctx, tt.args.userId, tt.args.trackId); (err != nil) != tt.wantErr {
				t.Errorf("source.DislikeTrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func MustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func Test_source_ShowLikedTracks(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.MusicDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: ShowLikedTracks source: show user's liked track",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
					Name:     "name",
					Release:  MustParseTime("2006-01-02", "2006-01-02"),
					FileName: "file_name",
					Size:     1000,
					Duration: "duration",
				},
			},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"name",
						"release_date",
						"file_name",
						"size",
						"duration",
					}).AddRow(
					uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
					"name",
					MustParseTime("2006-01-02", "2006-01-02"),
					"file_name",
					"1000",
					"duration",
				)

				f.db.ExpectQuery("SELECT music.* FROM music JOIN user_music ON music.id = user_music.music_id WHERE user_music.user_id = $1").
					WithArgs(
						a.id,
					).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error: ShowLikedTracks source: can't exec query",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery("ABRA-CADABRA").
					WithArgs(
						a.id,
					).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "error: ShowLikedTracks source: can't scan rows",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: nil,
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"name",
						"release_date",
						"file_name",
						"size",
						"duration",
					}).AddRow(
					uuid.MustParse("499afbff-7ff4-41e8-9f4d-9856669cca63"),
					"name",
					"puk-puk",
					"file_name",
					1000,
					"duration",
				)

				f.db.ExpectQuery("SELECT music.* FROM music JOIN user_music ON music.id = user_music.music_id WHERE user_music.user_id = $1").
					WithArgs(
						a.id,
					).WillReturnRows(rows).WillReturnError(fmt.Errorf("can't scan rows"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			usersSource := db.NewUserSourсe(db.NewSource(sqlx.NewDb(source, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := usersSource.ShowLikedTracks(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.ShowLikedTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("source.ShowLikedTracks() = %v, want %v", got, tt.want)
			}
		})
	}
}
