package app

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"mind-service/eventemitter"
	"mind-service/proto/gen/go/mind/v1/mindv1connect"
	"mind-service/service"
	"net/http"
)

var Version = "x" // overwritten by build pipeline

type App struct {
	Config         Config
	Logger         zerolog.Logger
	Service        service.Service
	ServiceHandler http.Handler
}

func New(c Config) App {
	a := App{Config: c}
	a.Logger = newLogger(a)
	a.Service = NewService(a)
	a.ServiceHandler = NewServiceHandler(a.Service)
	return a
}

func NewService(a App) service.Service {
	return service.Service{
		EventEmitter: newEventEmitter(a),
	}
}

func newEventEmitter(a App) eventemitter.EventEmitter {
	return eventemitter.EventEmitter{
		Secret: a.Config.Secret,
	}
}

func NewServiceHandler(svc service.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(mindv1connect.NewGreetServiceHandler(svc))
	return h2c.NewHandler(mux, &http2.Server{})
}

func (a App) Start() {
	log.Info().
		Str("version", Version).
		Interface("config", a.Config).
		Msg("Starting service")

	err := http.ListenAndServe(fmt.Sprint("0.0.0.0:", a.Config.ServerPort), a.ServiceHandler)
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
