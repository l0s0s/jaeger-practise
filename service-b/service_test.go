package serviceb_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	serviceb "github.com/l0s0s/jaeger-practise/service-b"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer_Talk(t *testing.T) {
	for _, tc := range []struct {
		testName             string // name of test
		expectedStatusCode   int    // staus code which service need to return
		availebleStatusCodes []int  // availeble status codes for return
	}{
		{
			"success: status OK",
			http.StatusOK,
			[]int{http.StatusOK},
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r := gin.New()

			h := serviceb.NewHTTPServer(tc.availebleStatusCodes...)

			h.BindRoutes(&r.RouterGroup)

			ts := httptest.NewServer(r)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+"/service-b/api/v1/talk", &bytes.Reader{})
			require.NoError(t, err)

			client := http.Client{}

			res, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}
