package service_test

import (
	"context"
	"mind-service/app"
	mindv1 "mind-service/proto/gen/go/mind/v1"
	"mind-service/proto/gen/go/mind/v1/mindv1connect"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	a := app.New(app.Config{})
	a.Authenticator = func() func(string) bool {
		return fakeAuthenticator
	}
	srv := httptest.NewServer(a.ServiceHandler())

	client := mindv1connect.NewGreetServiceClient(http.DefaultClient, srv.URL)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&mindv1.GreetRequest{Name: "Jane"}),
	)

	require.NoError(t, err)
	assert.Equal(t, "Hello, Jane!", res.Msg.Greeting)
}

func fakeAuthenticator(_ string) bool {
	return true
}
