package user_test

import (
	"context"
	"github.com/gabrielforster/voting/auth/user"
	"github.com/gabrielforster/voting/auth/user/mocks"
	tmocks "github.com/gabrielforster/voting/commom/telemetry/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	otel := tmocks.NewTelemetry(t)
	span := tmocks.NewSpan(t)
	span.On("RecordError", mock.Anything).Return(nil)
	span.On("SetStatus", mock.Anything, mock.Anything).Return(nil)
	span.On("End").Return(nil)
	otel.On("Start", ctx, "validatePassword").Return(ctx, span)

	s := user.NewService(repo, otel)
	u := &user.User{
		Email:    "rochafrgabriel@gmail.com",
		Password: "7110eda4d09e062aa5e4a390b0a572ac0d2c0220",
	}
	t.Run("invalid password", func(t *testing.T) {
		err := s.ValidatePassword(ctx, u, "invalid")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid password", err.Error())
	})
	t.Run("valid password", func(t *testing.T) {
		err := s.ValidatePassword(ctx, u, "12345")
		assert.Nil(t, err)
	})
}
