package database

import (
	"context"
	"database/sql"

	"github.com/gabrielforster/voting/auth/user"
	"github.com/gabrielforster/voting/commom/telemetry"

	"go.opentelemetry.io/otel/codes"
)

// DBConn mysql repo
type DBConn struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

func NewConn(db *sql.DB, telemetry telemetry.Telemetry) *DBConn {
	return &DBConn{
		db:        db,
		telemetry: telemetry,
	}
}

func (r *DBConn) Get(ctx context.Context, email string) (*user.User, error) {
	ctx, span := r.telemetry.Start(ctx, "mysql")
	defer span.End()
	stmt, err := r.db.Prepare(`select id, email, password, first_name, last_name from user where email = ?`)
	if err != nil {
		return nil, err
	}
	var u user.User
	rows, err := stmt.Query(email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}
	return &u, nil
}
