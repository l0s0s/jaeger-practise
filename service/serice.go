package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const (
	ctxTimeout = time.Second * 3
)

type HTTPServer struct {
	name    string
	nextURL string
	tp      trace.TracerProvider
}

func NewHTTPServer(name, nextURL string, tp trace.TracerProvider) *HTTPServer {
	return &HTTPServer{
		name:    name,
		nextURL: nextURL,
		tp:      tp,
	}
}

func (h *HTTPServer) Talk(c *gin.Context) {
	if h.nextURL == "" {
		c.JSON(http.StatusOK, "chain end")

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)

	defer cancel()

	client := http.Client{}

	_, span := h.tp.Tracer("http-request").Start(ctx, h.nextURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.nextURL, &bytes.Reader{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
		return
	}

	span.End()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
		return
	}

	c.Status(resp.StatusCode)
}

func (h *HTTPServer) BindRoutes(g *gin.RouterGroup) {
	g.GET(fmt.Sprintf("/service/%s/api/v1/talk", h.name), h.Talk)
}
