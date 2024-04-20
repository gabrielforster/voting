package poll

import (
	"context"
)

type Repository interface {
    Store(ctx context.Context, poll *Poll, options []string, user_id string) error
}
