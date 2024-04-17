package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/poll/poll"

	"go.opentelemetry.io/otel/codes"
)

type PollMySQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

func NewVoteMySQL(db *sql.DB, telemetry telemetry.Telemetry) *PollMySQL {
	return &PollMySQL{
		db:        db,
		telemetry: telemetry,
	}
}

func (r *PollMySQL) Store(ctx context.Context, p *poll.Poll) (string, error) {
	return "", nil
}
