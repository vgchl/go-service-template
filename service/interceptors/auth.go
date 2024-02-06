package interceptors

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

type authenticator func(string) bool

func NewAuthInterceptor(authenticate authenticator) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			token := req.Header().Get("authorization")
			if !authenticate(token) {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("unauthenticated"),
				)
			}
			return next(ctx, req)
		}
	}
	return interceptor
}
