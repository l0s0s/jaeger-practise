package servicea

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	contextTimeout = time.Second * 3
)

type HTTPServer struct {
	talkEdpoint string
}

func NewHTTPServer(talkEdpoint string) *HTTPServer {
	return &HTTPServer{talkEdpoint: talkEdpoint}
}

func (h *HTTPServer) Talk(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)

	defer cancel()

	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.talkEdpoint, &bytes.Reader{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
		return
	}

	c.Status(resp.StatusCode)
}

func (h *HTTPServer) BindRoutes(g *gin.RouterGroup) {
	g.GET("/service-a/api/v1/talk", h.Talk)
}
