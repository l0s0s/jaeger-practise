package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/l0s0s/jaeger-practise/config"
	"github.com/l0s0s/jaeger-practise/service"
	"go.opentelemetry.io/otel/exporters/jaeger"
)

func tracerProvider(c config.Config) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Jaeger.URL)))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.Service.Name),
			attribute.String("ID", c.Service.ID),
		)),
	)

	return tp, nil
}

func main() {
	c := config.Parse()

	log := zerolog.New(os.Stderr).Level(zerolog.ErrorLevel)

	tp, err := tracerProvider(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to init trace provider")
	}

	otel.SetTracerProvider(tp)

	otel.SetLogger(logr.New(logr.Discard().GetSink()))

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown trace provider")
		}
	}()

	s := service.NewHTTPServer(c.Service.Name, c.NextURL, tp)

	r := gin.New()

	r.Use(otelgin.Middleware(c.Service.Name, otelgin.WithTracerProvider(tp)))
	s.BindRoutes(&r.RouterGroup)

	if err := r.Run(c.Service.Port); err != nil {
		log.Error().Err(err).Msg("failed to listen and serve")
	}
}
