package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gabrielforster/voting/commom/telemetry"
	"github.com/gabrielforster/voting/poll/handlers"
	pollService "github.com/gabrielforster/voting/poll/poll"
	pollDatabase "github.com/gabrielforster/voting/poll/poll/database"
	"go.opentelemetry.io/otel/codes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	telemetrymiddleware "github.com/go-chi/telemetry"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := httplog.NewLogger("poll", httplog.Options{
		JSON: true,
	})

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer db.Close()

	ctx := context.Background()
	otel, err := telemetry.NewJaeger(ctx, "poll")
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer otel.Shutdown(ctx)

	pRepo := pollDatabase.NewVoteMySQL(db, otel)
	pService := pollService.NewService(pRepo, otel)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(telemetrymiddleware.Collector(telemetrymiddleware.Config{
		AllowAny: true,
	}, []string{"/v1"}))

	r.Post("/v1/poll", handlers.CreatePoll(ctx, pService, otel))
	// r.Get("/v1/poll", validateToken(ctx, otel))
	//
	// r.Post("/v1/vote", userAuth(ctx, pService, otel))
	// r.Get("/v1/vote", userAuth(ctx, pService, otel))

	http.Handle("/", r)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      http.DefaultServeMux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
}
