package database

import (
	"context"
	"database/sql"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/poll/option"
	"github.com/gabrielforster/voting/poll/poll"
	"go.opentelemetry.io/otel/codes"
)

type PollMySQL struct {
	db             *sql.DB
	telemetry      telemetry.Telemetry
	optionsService option.UseCase
}

func NewPollMySQL(db *sql.DB, telemetry telemetry.Telemetry, optionsService option.UseCase) *PollMySQL {
	return &PollMySQL{
		db:             db,
		telemetry:      telemetry,
		optionsService: optionsService,
	}
}

func (r *PollMySQL) Store(ctx context.Context, p *poll.Poll, options []string, user_id string) error {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()

	tx, err := r.db.Begin()
	stmt, err := tx.Prepare(`
        insert into polls (title, description, hash, created_by)
        values (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}

	rows, err := stmt.Exec(p.Title, p.Description, p.Hash, user_id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var parsedOptions []*option.Option

	for _, o := range options {
		parsedOptions = append(parsedOptions, &option.Option{
			Title:  o,
			PollId: p.Id,
		})
	}

	err = r.optionsService.CreateOptions(ctx, parsedOptions, tx)
	if err != nil {
		tx.Rollback()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	insertedId, err := rows.LastInsertId()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	p.Id = int(insertedId)

	return nil
}
