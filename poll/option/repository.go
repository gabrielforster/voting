package option

import (
	"context"
	"database/sql"
)

type Repository interface {
	Store(ctx context.Context, option *Option, trx *sql.Tx) error
	Get(ctx context.Context, optionId int, trx *sql.Tx) (*Option, error)
	List(ctx context.Context, pollId int, trx *sql.Tx) ([]*Option, error)
}
