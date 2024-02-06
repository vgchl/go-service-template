package app

import (
	"errors"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (a *App) Start() {
	log.Info().
		Str("version", Version).
		Interface("config", a.Config).
		Msg("Starting service")

	err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", a.Config.ServerPort), a.ServiceHandler())
	if !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("Failed to start server")
	}
}
