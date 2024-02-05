package service

import (
	"connectrpc.com/connect"
	"context"
	v1 "mind-svc/proto/gen/go/mind/v1"
)

type Service struct {
}

func (h Service) Greet(ctx context.Context, req *connect.Request[v1.GreetRequest]) (*connect.Response[v1.GreetResponse], error) {
	//TODO implement me
	panic("implement me")
}
