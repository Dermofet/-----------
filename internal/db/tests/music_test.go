package db

import (
	"bytes"
	"context"
	"fmt"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_source_GetAll(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get all music",
			args: args{ctx: ctx},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"name",
						"release_date",
						"file_name",
						"size",
						"duration",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"Song1",
						time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
						"Song1.mp3",
						uint64(500),
						"2:47",
					).
					AddRow(
						uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						"Song2",
						time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						"Song2.mp3",
						uint64(900),
						"3:23",
					)
				f.db.ExpectQuery("SELECT * FROM music").WillReturnRows(rows)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
				{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				},
			},
			wantErr: false,
		},
		{
			name: "Bad request to database at music.GetAll",
			args: args{ctx: ctx},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("...").WillReturnError(fmt.Errorf("can't exec query"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				db: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := musicSource.GetAll(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_source_Get(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}

	type args struct {
		ctx     context.Context
		musicId uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    *entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get music",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f fields) {
				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"release_date",
					"file_name",
					"size",
					"duration",
				}).
					AddRow(
						uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						"Song1",
						time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						"Song1.mp3",
						uint64(900),
						"3:23",
					)
				f.db.ExpectQuery("SELECT * FROM music WHERE id = $1").WithArgs(a.musicId.String()).WillReturnRows(rows)
			},
			want: &entity.MusicDB{
				Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				Name:     "Song1",
				Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
				FileName: "Song1.mp3",
				Size:     uint64(900),
				Duration: "3:23",
			},
			wantErr: false,
		},
		{
			name: "Not music in base at music.Get",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f fields) {
				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"release_date",
					"file_name",
					"size",
					"duration",
				})
				f.db.ExpectQuery("SELECT * FROM music WHERE id = $1").WithArgs(a.musicId.String()).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error: can't exec query at music.Get",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("...").WithArgs(a.musicId.String()).WillReturnError(fmt.Errorf("can't exec query"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				db: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := musicSource.Get(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_source_GetAndSortByPopular(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}

	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get all sorted by popular music",
			args: args{ctx: ctx},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"name",
						"release_date",
						"file_name",
						"size",
						"duration",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"Song1",
						time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
						"Song1.mp3",
						uint64(500),
						"2:47",
					).
					AddRow(
						uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						"Song2",
						time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						"Song2.mp3",
						uint64(900),
						"3:23",
					)
				f.db.ExpectQuery("SELECT m.id AS id, m.name AS name " +
					"FROM music m " +
					"LEFT JOIN user_music um ON um.music_id = m.id " +
					"GROUP BY m.id, m.name " +
					"ORDER BY COALESCE(COUNT(um.music_id), 0) DESC;").WillReturnRows(rows)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
				{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				},
			},
			wantErr: false,
		},
		{
			name: "Bad request to database at music.GetAndSortByPopular",
			args: args{ctx: ctx},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("...").WillReturnError(fmt.Errorf("can't exec query"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				db: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := musicSource.GetAndSortByPopular(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}

		})
	}
}

func Test_source_GetAllSortByTime(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}

	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get all sorted by reliase music",
			args: args{ctx: ctx},
			setup: func(a args, f fields) {
				rows := sqlmock.
					NewRows([]string{
						"id",
						"name",
						"release_date",
						"file_name",
						"size",
						"duration",
					}).
					AddRow(
						uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						"Song1",
						time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
						"Song1.mp3",
						uint64(500),
						"2:47",
					).
					AddRow(
						uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						"Song2",
						time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						"Song2.mp3",
						uint64(900),
						"3:23",
					)
				f.db.ExpectQuery("SELECT * FROM music ORDER BY release_date").WillReturnRows(rows)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
				{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				},
			},
			wantErr: false,
		},
		{
			name: "Bad request to database at music.GetAllSortByTime",
			args: args{ctx: ctx},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("...").WillReturnError(fmt.Errorf("can't exec query"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				db: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			got, err := musicSource.GetAllSortByTime(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_source_Create(t *testing.T) {
	type fields struct {
		sqlmock sqlmock.Sqlmock
	}

	type args struct {
		ctx     context.Context
		musicDB *entity.MusicDB
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "Create music",
			args: args{
				ctx: ctx,
				musicDB: &entity.MusicDB{
					Id:       uuid.Nil,
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
			},
			setup: func(a args, f fields) {
				reader := bytes.NewReader([]byte("1111111111111111"))
				uuid.SetRand(reader)
				uuid.SetClockSequence(0)

				a.musicDB.Id = uuid.MustParse("31313131-3131-4131-b131-313131313131")
				rows := sqlmock.NewResult(1, 1)
				f.sqlmock.ExpectExec("INSERT INTO music (id, name, release_date, file_name, size, duration) VALUES ($1, $2, $3, $4, $5, $6)").
					WithArgs(a.musicDB.Id, a.musicDB.Name, a.musicDB.Release, a.musicDB.FileName, a.musicDB.Size, a.musicDB.Duration).
					WillReturnResult(rows)
			},
			wantErr: false,
		}, {
			name: "Bad database request",
			args: args{
				ctx: ctx,
				musicDB: &entity.MusicDB{
					Id:       uuid.Nil,
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
			},
			setup: func(a args, f fields) {
				reader := bytes.NewReader([]byte("1111111111111111"))
				uuid.SetRand(reader)
				uuid.SetClockSequence(0)

				a.musicDB.Id = uuid.MustParse("31313131-3131-4131-b131-313131313131")
				f.sqlmock.ExpectExec("...").
					WithArgs(a.musicDB.Id, a.musicDB.Name, a.musicDB.Release, a.musicDB.FileName, a.musicDB.Size, a.musicDB.Duration).
					WillReturnError(fmt.Errorf("Bad request"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				sqlmock: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			err = musicSource.Create(tt.args.ctx, tt.args.musicDB)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_source_Update(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}

	type args struct {
		ctx     context.Context
		musicDb *entity.MusicDB
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "Update music with file",
			args: args{
				ctx: ctx,
				musicDb: &entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
			},
			setup: func(a args, f fields) {
				rows := sqlmock.NewResult(1, 1)
				f.db.ExpectExec("UPDATE music SET name = $2, release_date = $3, file_name = $4, size = $5, duration = $6 WHERE id = $1").
					WithArgs(a.musicDb.Id, a.musicDb.Name, a.musicDb.Release, a.musicDb.FileName, a.musicDb.Size, a.musicDb.Duration).
					WillReturnResult(rows)
			},
			wantErr: false,
		},
		{
			name: "Update music without file",
			args: args{
				ctx: ctx,
				musicDb: &entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "",
					Size:     0,
					Duration: "",
				},
			},
			setup: func(a args, f fields) {
				rows := sqlmock.NewResult(1, 1)
				f.db.ExpectExec("UPDATE music SET name = $2, release_date = $3 WHERE id = $1").
					WithArgs(a.musicDb.Id, a.musicDb.Name, a.musicDb.Release).WillReturnResult(rows)
			},
			wantErr: false,
		},
		{
			name: "Bad request to database at music.Create",
			args: args{
				ctx: ctx,
				musicDb: &entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
			},
			setup: func(a args, f fields) {
				f.db.ExpectQuery("...").
					WithArgs(a.musicDb.Id, a.musicDb.Name, a.musicDb.Release, a.musicDb.FileName, a.musicDb.Size, a.musicDb.Duration).
					WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				db: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			err = musicSource.Update(tt.args.ctx, tt.args.musicDb)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_source_Delete(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
	}

	type args struct {
		ctx     context.Context
		musicId uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "Delete music",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f fields) {
				f.db.ExpectExec("DELETE FROM music WHERE id = $1").WithArgs(a.musicId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Bad request to database at music.Delete",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f fields) {
				f.db.ExpectExec("...").WithArgs(a.musicId).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer database.Close()

			f := fields{
				db: mock,
			}

			musicSource := db.NewMusicSource(db.NewSource(sqlx.NewDb(database, "sqlmock")))

			tt.setup(tt.args, f)

			err = musicSource.Delete(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
