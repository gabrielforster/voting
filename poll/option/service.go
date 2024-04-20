package option

import (
	"context"
	"database/sql"

	"github.com/gabrielforster/voting/commom/telemetry"
)

type UseCase interface {
	CreateOptions(ctx context.Context, option []*Option, tx *sql.Tx) error
	GetOption(ctx context.Context, optionId int) (*Option, error)
	ListOptions(ctx context.Context, pollId int) (*[]Option, error)
}

type Service struct {
	repo      Repository
	telemetry telemetry.Telemetry
}

func NewService(repo Repository, telemetry telemetry.Telemetry) *Service {
	return &Service{
		repo:      repo,
		telemetry: telemetry,
	}
}
func (s *Service) CreateOptions(ctx context.Context, options []*Option, tx *sql.Tx) error {
	ctx, span := s.telemetry.Start(ctx, "service.create_options")
	defer span.End()

    for _, o := range options {
        err := s.repo.Store(ctx, o, tx)
        if err != nil {
            return err
        }
    }

	return nil
}
func (s *Service) GetOption(ctx context.Context, optionId int) (*Option, error) {
	ctx, span := s.telemetry.Start(ctx, "service.get_option")
	defer span.End()

	return nil, nil
}

func (s *Service) ListOptions(ctx context.Context, pollId int) (*[]Option, error) {
	ctx, span := s.telemetry.Start(ctx, "service.list_options")
	defer span.End()

	return nil, nil
}
