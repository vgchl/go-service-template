package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/skovtunenko/graterm"
)

func (a *App) Start() {
	log.Info().
		Str("version", Version).
		Interface("config", a.Config).
		Msg("Starting service")

	term, ctx := graterm.NewWithSignals(context.Background(), os.Interrupt, syscall.SIGTERM)
	term.SetLogger(&log.Logger)

	go func() {
		err := a.Server().ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	term.WithOrder(graterm.Order(1)).WithName("http server").
		Register(a.Config.TerminationTimeout, func(ctx context.Context) {
			err := a.Server().Shutdown(ctx)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to gracefully shut down server")
			}
		})

	err := term.Wait(ctx, a.Config.TerminationTimeout)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to terminate gracefully within timeout")
	}
	log.Info().Msg("Application has shut down.")
}

func (a *App) OnStop() {
	r := recover()
	if r != nil {
		log.WithLevel(zerolog.PanicLevel).
			Interface(zerolog.ErrorFieldName, r).
			Msg("Application terminated due to panic")
	}
	os.Exit(1)
}
