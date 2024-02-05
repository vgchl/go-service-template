package app

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"mind-service/proto/gen/go/mind/v1/mindv1connect"
	"mind-service/service"
	"net/http"
)

type App struct {
	Config  Config
	Logger  zerolog.Logger
	Service service.Service
	Server  *http.ServeMux
}

func New(c Config) App {
	a := App{Config: c}
	a.Logger = newLogger(a)
	a.Service = service.Service{}
	a.Server = newServer(a.Service)
	return a
}

func newServer(svc service.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(mindv1connect.NewGreetServiceHandler(svc))
	return mux
}

func (a App) Start() {
	log.Info().
		Int("port", a.Config.ServerPort).
		Msg("Starting service")

	err := http.ListenAndServe(
		fmt.Sprint("0.0.0.0:", a.Config.ServerPort),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(a.Server, &http2.Server{}),
	)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("Failed to start server")
	}
}

func newLogger(s App) zerolog.Logger {
	if !s.Config.LogJson {
		log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	}
	zerolog.DefaultContextLogger = &log.Logger
	return log.Logger
}
