package app

import (
	"mind-service/eventemitter"
	"mind-service/proto/gen/go/mind/v1/mindv1connect"
	"mind-service/service"
	"mind-service/service/interceptors"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var Version = "unknown" // overwritten by build pipeline

// App contains all dependencies required to run the application.
// Dependencies are initialized using a memoized, lazy loaded function. This reduces the need to pass around variables,
// build dependencies in a strict order, and most notably; it allows individual dependencies to be replaced by mocks
// in tests.
//
// The pattern goes:
//
// 1: Define a field on App that will store the builder:
//
//	type App struct {
//		MyDependency func() some.MyDependency
//	}
//
// 2: Assign the builder to the field. Usually you'll want the same instance to be returned every time. This can be
// done by wrapping the builder with the memoize() helper.
//
//	func New(c Config) *App {
//		a := &App{Config: c}
//		...
//		a.MyDependency = memoize(a.newMyDependency)
//		...
//	}
//
// 3: Define the builder method. The naming convention is `new<FieldName>()`.
// The method cannot take parameters. It can access config or other dependencies through the `a *App` receiver.
// When accessing other dependencies, use the memoized field, not the builder method directly.
//
//	func (a *App) newMyDependency() some.MyDependency {
//		return some.MyDependency{
//			OtherDependency: a.OtherDependency(),
//		}
//	}
//
// 4: Retrieve the dependency where required.
//
//	d := a.MyDependency()
type App struct {
	Config Config
	Logger zerolog.Logger

	Service        func() service.Service
	ServiceHandler func() http.Handler
	Server         func() *http.Server
	EventEmitter   func() eventemitter.EventEmitter
	Authenticator  func() func(string) bool
}

func New(c Config) *App {
	a := &App{Config: c}
	a.configureLogger()

	a.Service = memoize(a.newService)
	a.ServiceHandler = memoize(a.newServiceHandler)
	a.Server = memoize(a.newServer)
	a.EventEmitter = memoize(a.newEventEmitter)
	a.Authenticator = memoize(a.newAuthenticator)

	return a
}

func (a *App) newService() service.Service {
	return service.Service{
		EventEmitter: a.EventEmitter(),
	}
}

func (a *App) newEventEmitter() eventemitter.EventEmitter {
	return eventemitter.EventEmitter{Secret: a.Config.Secret}
}

func (a *App) newServiceHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(mindv1connect.NewGreetServiceHandler(
		a.Service(),
		connect.WithInterceptors(
			interceptors.NewAuthInterceptor(a.Authenticator()),
		),
	))
	return h2c.NewHandler(mux, &http2.Server{})
}

func (a *App) newAuthenticator() func(token string) bool {
	return func(token string) bool {
		return false
	}
}

func (a *App) newServer() *http.Server {
	return &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", a.Config.ServerPort),
		Handler: a.ServiceHandler(),
	}
}

func (a *App) configureLogger() {
	if !a.Config.LogJson {
		log.Logger = log.Logger.Output(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	}

	level, err := zerolog.ParseLevel(a.Config.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid log level configuration")
	}
	log.Logger = log.Logger.Level(level)

	zerolog.DefaultContextLogger = &log.Logger
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
