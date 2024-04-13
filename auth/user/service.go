package user

import (
    "context"
    "crypto/sha1"
    "fmt"

    "github.com/gabrielforster/voting/commom/telemetry"

	"go.opentelemetry.io/otel/codes"
)

type UseCase interface {
    ValidateUser(ctx context.Context, email string, password string) error
}

type Service struct {
    repository Repository
    telemetry telemetry.Telemetry
}

func NewService(repository Repository, telemetry telemetry.Telemetry) *Service {
    return &Service{
        repository: repository,
        telemetry: telemetry,
    }
}

func (s *Service) ValidateUser (ctx context.Context, email string, password string) error {
    ctx, span := s.telemetry.Start(ctx, "service")
	defer span.End()
	u, err := s.repository.Get(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if u == nil {
		err := fmt.Errorf("invalid user")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return s.ValidatePassword(ctx, u, password)
}


func (s *Service) ValidatePassword (ctx context.Context, user *User, password string) error {
	ctx, span := s.telemetry.Start(ctx, "validatePassword")
	defer span.End()

	passwordHash := sha1.New()
	passwordHash.Write([]byte(password))
	p := fmt.Sprintf("%x", passwordHash.Sum(nil))

	if p != user.Password {
		err := fmt.Errorf("invalid password")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
