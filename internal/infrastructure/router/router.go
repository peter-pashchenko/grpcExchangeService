package router

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/zap"
)

const (
	jaegerEndPoint = "/api/traces"
	promEndpoint   = "/metrics"
)

type Struct struct {
	r  *chi.Mux
	tp *sdktrace.TracerProvider

	log *zap.Logger
}

type Option func(*Struct)

func New(logger *zap.Logger, opts ...Option) *Struct {
	res := &Struct{r: chi.NewRouter(), log: logger}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (s *Struct) GetRouter() *chi.Mux {
	return s.r
}
func (s *Struct) TPShutdown(ctx context.Context) error {
	return s.tp.Shutdown(ctx)
}

func WithJaeger(httpHost string) Option {
	return func(s *Struct) {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://" + httpHost + jaegerEndPoint)))
		if err != nil {
			s.log.Error("error creating jaeger exporter",
				zap.Error(err))
		}

		s.log.Debug("jaeger exporter created")

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("exchageRateService"),
				),
			))

		otel.SetTracerProvider(tp)

		s.tp = tp

		s.log.Debug(fmt.Sprintf("trace provider is set,collecting with endpoint:%s", httpHost+jaegerEndPoint))
	}
}

func WithPrometheus() Option {
	return func(r *Struct) {
		r.r.Handle(promEndpoint, promhttp.Handler())
	}
}
