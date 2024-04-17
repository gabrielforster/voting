package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/poll/vote"

	"go.opentelemetry.io/otel/codes"
)

type VoteMySQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

func NewVoteMySQL(db *sql.DB, telemetry telemetry.Telemetry) *VoteMySQL {
	return &VoteMySQL{
		db:        db,
		telemetry: telemetry,
	}
}

func (r *VoteMySQL) Store(ctx context.Context, v *vote.Vote) error {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()
	stmt, err := r.db.Prepare(`
		insert into vote (id, email, talk_name, score, created_at)
		values(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		v.ID,
		v.Email,
		v.TalkName,
		v.Score,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		span.RecordError(err)
		return err
	}
	err = stmt.Close()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
