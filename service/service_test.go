package service_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/l0s0s/jaeger-practise/service"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func startMockServer() string {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	ts := httptest.NewServer(r)

	return ts.URL + "/"
}

func TestHTTPServer_Talk(t *testing.T) {
	for _, tc := range []struct {
		testName           string // name of test
		talkEndpoint       string // endpoint on which service would be make request
		expectedStatusCode int    // staus code which service need to return
	}{
		{
			"succes: Status OK",
			startMockServer(),
			http.StatusOK,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r := gin.New()

			h := service.NewHTTPServer("test", os.Getenv(tc.talkEndpoint), trace.NewNoopTracerProvider())

			h.BindRoutes(&r.RouterGroup)

			ts := httptest.NewServer(r)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+"/service/test/api/v1/talk", &bytes.Reader{})
			require.NoError(t, err)

			client := http.Client{}

			res, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}
