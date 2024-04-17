package database

import (
	"context"
	"database/sql"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/poll/poll"
	"go.opentelemetry.io/otel/codes"
	// "go.opentelemetry.io/otel/codes"
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

func (r *PollMySQL) Store(ctx context.Context, p *poll.Poll) error {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()
	stmt, err := r.db.Prepare(`
        insert into polls (title, description, hash, created_by)
        values (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}

    // TODO: use p.CreatedBy (validate user session on service)
	rows, err := stmt.Exec(p.Title, p.Description, p.Slug, 12893)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

    // TODO: insert options (insert poll and options with sql transaction)

    insertedId, err := rows.LastInsertId()
    if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
    }

    p.Id = int(insertedId)

	return nil
}

