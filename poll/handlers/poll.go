package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/poll/poll"

	"github.com/go-chi/httplog"
	"go.opentelemetry.io/otel/codes"
)

func CreatePoll(ctx context.Context, pService poll.UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		ctx, span := otel.Start(ctx, "service.create_poll")
		defer span.End()

		user_id := r.Context().Value("user_id")

		var param struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Options     []string `json:"options"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}

		var p poll.Poll
		err = pService.CreatePoll(ctx, &p, param.Options, user_id.(string))
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}

		var result struct {
			Id       string `json:"indentifier"`
			Location string `json:"location"`
		}
		result.Id = p.Hash
		result.Location = "poll/" + p.Hash

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}

		return
	}
}
