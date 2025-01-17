package grpc

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

var (
	methodCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_method_counter",
		},
		[]string{"method", "status"},
	)

	methodDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "grpc_method_duration",
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(methodCounter)
	prometheus.MustRegister(methodDuration)
}

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	var status string
	if err == nil {
		status = "success"
	} else {
		status = "error"
	}
	methodCounter.WithLabelValues(info.FullMethod, status).Inc()
	methodDuration.WithLabelValues(info.FullMethod, status).Observe(time.Since(start).Seconds())

	return resp, err

}
