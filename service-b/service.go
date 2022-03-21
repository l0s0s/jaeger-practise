package serviceb

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	availebleStatusCodes []int
}

func NewHTTPServer(availebleStatusCodes ...int) *HTTPServer {
	return &HTTPServer{
		availebleStatusCodes: availebleStatusCodes,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (h *HTTPServer) Talk(c *gin.Context) {
	code := h.availebleStatusCodes[rand.Intn(len(h.availebleStatusCodes))]

	c.Status(code)
}

func (h *HTTPServer) BindRoutes(g *gin.RouterGroup) {
	g.GET("/service-b/api/v1/talk", h.Talk)
}
