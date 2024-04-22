package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gabrielforster/voting/commom/telemetry"

	"go.opentelemetry.io/otel/codes"
)

func ValidateToken(ctx context.Context, telemetry telemetry.Telemetry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			_, span := telemetry.Start(ctx, "ValidateToken")
			defer span.End()

			errMessage := "Unauthorized"
			token := r.Header.Get("Authorization")

			if token == "" {
				err := errors.New("Unauthorized")
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errMessage)
				return
			}

			payload := `{
			"token": "` + token + `"
            }`

			req, err := http.Post(os.Getenv("AUTH_URL")+"/v1/validate_token", "application/json", strings.NewReader(payload))
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errMessage)
				return
			}

			defer req.Body.Close()

			type result struct {
				Id string `json:"user_id"`
			}
			var res result
			err = json.NewDecoder(req.Body).Decode(res)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errMessage)
				return
			}

			newCtx := context.WithValue(r.Context(), "user_id", res.Id)
			next.ServeHTTP(rw, r.WithContext(newCtx))
		}

		return http.HandlerFunc(fn)
	}
}

func respondWithError(w http.ResponseWriter, code int, e string, message string) {
	respondWithJSON(w, code, map[string]string{"code": strconv.Itoa(code), "error": e, "message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
