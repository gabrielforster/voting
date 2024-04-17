package poll

import (
	"context"
)

type Repository interface {
	Store(ctx context.Context, v *Poll) error
}
