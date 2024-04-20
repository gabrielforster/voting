package poll

import (
	"context"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/commom/utils"

	"go.opentelemetry.io/otel/codes"
)

type UseCase interface {
	CreatePoll(ctx context.Context, poll *Poll) error
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
func (s *Service) CreatePoll(ctx context.Context, poll *Poll) error {
	ctx, span := s.telemetry.Start(ctx, "service")
	defer span.End()

	poll.Slug = utils.StringToHash(poll.Title)

	err := s.repo.Store(ctx, poll)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
