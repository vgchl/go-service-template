package service

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"mind-service/eventemitter"

	v1 "mind-service/proto/gen/go/mind/v1"
)

type Service struct {
	EventEmitter eventemitter.EventEmitter
}

func (s Service) Greet(ctx context.Context, req *connect.Request[v1.GreetRequest]) (*connect.Response[v1.GreetResponse], error) {
	log.Info().Interface("headers", req.Header()).Msg("Request headers")
	res := connect.NewResponse(&v1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
