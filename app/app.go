package app

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"mind-service/proto/gen/go/mind/v1/mindv1connect"
	"mind-service/service"
	"net/http"
)

type App struct {
	Config Config
	Mux    *http.ServeMux
}

func New(c Config) App {
	s := App{Config: c}

	logging(s)

	s.Mux = server(service.Service{})

	return s
}

func server(svc service.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(mindv1connect.NewGreetServiceHandler(svc))
	return mux
}

func (s App) Start() {
	log.Info().
		Int("port", s.Config.ServerPort).
		Msg("Starting service")

	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(s.Mux, &http2.Server{}),
	)
}

func logging(s App) {
	if !s.Config.LogJson {
		log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	}
	zerolog.DefaultContextLogger = &log.Logger
}
