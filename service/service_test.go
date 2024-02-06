package service_test

import (
	"connectrpc.com/connect"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mind-service/app"
	"mind-service/proto/gen/go/mind/v1"
	"mind-service/proto/gen/go/mind/v1/mindv1connect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	svc := app.NewService(app.App{})
	srv := httptest.NewServer(app.NewServiceHandler(svc))

	client := mindv1connect.NewGreetServiceClient(http.DefaultClient, srv.URL)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&mindv1.GreetRequest{Name: "Jane"}),
	)

	require.NoError(t, err)
	assert.Equal(t, "Hello, Jane!", res.Msg.Greeting)
}
