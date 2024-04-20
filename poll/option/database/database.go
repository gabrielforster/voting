package database

import (
	"context"
	"database/sql"

	"github.com/gabrielforster/voting/commom/telemetry"
    // FIX: not finding utils on mod tidy
	"github.com/gabrielforster/voting/commom/utils"
	"github.com/gabrielforster/voting/poll/option"
	"go.opentelemetry.io/otel/codes"
)

type OptionMySQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

func NewOptionMySQL(db *sql.DB, telemetry telemetry.Telemetry) *OptionMySQL {
	return &OptionMySQL{
		db:        db,
		telemetry: telemetry,
	}
}

func (r *OptionMySQL) Store(ctx context.Context, option *option.Option, trx *sql.Tx) error {
	ctx, span := r.telemetry.Start(ctx, "mysql.option.store")
	defer span.End()

	prepareFunc := utils.Ternary(trx != nil, trx.Prepare, r.db.Prepare)

	stmt, err := prepareFunc(`
        insert into options (title, poll_id)
        values (?, ?)
    `)
	if err != nil {
		return err
	}

	rows, err := stmt.Exec(option.Title, option.PollId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	option.Id = int(id)

	return nil
}

func (r *OptionMySQL) Get(ctx context.Context, optionId int, trx *sql.Tx) (*option.Option, error) {
	ctx, span := r.telemetry.Start(ctx, "mysql.option.get")
	defer span.End()

	prepareFunc := utils.Ternary(trx != nil, trx.Prepare, r.db.Prepare)

	stmt, err := prepareFunc(`
        select id, title, poll_id
        from options
        where id = ?
    `)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(optionId)

	var opt option.Option
	err = row.Scan(&opt.Id, &opt.Title, &opt.PollId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &opt, nil
}

func (r *OptionMySQL) List(ctx context.Context, pollId int, trx *sql.Tx) ([]*option.Option, error) {
	ctx, span := r.telemetry.Start(ctx, "mysql.option.list")
	defer span.End()

	prepareFunc := utils.Ternary(trx != nil, trx.Prepare, r.db.Prepare)

	stmt, err := prepareFunc(`
        select id, title, poll_id
        from options
        where poll_id = ?
    `)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(pollId)
	if err != nil {
		return nil, err
	}

	var options []*option.Option
	for rows.Next() {
		var opt option.Option
		err = rows.Scan(&opt.Id, &opt.Title, &opt.PollId)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		options = append(options, &opt)
	}

	return options, nil
}
