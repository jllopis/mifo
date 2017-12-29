package mifo

import (
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jllopis/mifo/log"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusInterceptorConfig struct {
	EndPoint            string
	EnableTimeHistogram bool
}

func (p *PrometheusInterceptorConfig) Init(g *GrpcServer) {
	log.Info("called Init on Prometheus Interceptor")
	if p.EnableTimeHistogram {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}
	grpc_prometheus.Register(g.GrpcSrv)
}

func (g *GrpcServer) UsePrometheus(conf *PrometheusInterceptorConfig) *GrpcServer {
	// Default interceptors, [prometheus, opentracing]
	g.UseUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor)
	g.UseStreamInterceptor(grpc_prometheus.StreamServerInterceptor)

	if conf != nil && conf.EndPoint != "" {
		g.HttpMux.Handle(conf.EndPoint, prometheus.Handler())
	} else {
		g.HttpMux.Handle("/metrics", prometheus.Handler())
	}

	g.RegisterInterceptoInitializer(conf)

	return g
}
