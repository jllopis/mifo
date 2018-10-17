package mifo

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/codahale/metrics"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jllopis/mifo/log"
	"github.com/soheilhy/cmux"
)

type options struct {
	cors                  bool
	metrics               bool
	logRequests           bool
	outgoingHeaderMatcher runtime.HeaderMatcherFunc
}

type Option func(*options)

func WithCORS() Option {
	return func(o *options) {
		o.cors = true
	}
}

func WithMetrics() Option {
	return func(o *options) {
		o.metrics = true
	}
}

// DefaultHeaderMatcher is used to pass http request headers to/from gRPC context. This adds permanent HTTP header
// keys (as specified by the IANA) to gRPC context with grpcgateway- prefix. HTTP headers that start with
// 'Grpc-Metadata-' are mapped to gRPC metadata after removing prefix 'Grpc-Metadata-'.
func DefaultHeaderMatcher(key string) (string, bool) {
	if isPermanentHTTPHeader(key) {
		return runtime.MetadataPrefix + key, true
	} else if strings.HasPrefix(strings.ToLower(key), strings.ToLower(runtime.MetadataHeaderPrefix)) {
		return key[len(runtime.MetadataHeaderPrefix):], true
	}
	return "", false
}

// isPermanentHTTPHeader checks whether hdr belongs to the list of
// permenant request headers maintained by IANA.
// http://www.iana.org/assignments/message-headers/message-headers.xml
func isPermanentHTTPHeader(hdr string) bool {
	switch hdr {
	case
		"Accept",
		"Accept-Charset",
		"Accept-Language",
		"Accept-Ranges",
		"Authorization",
		"Cache-Control",
		"Content-Type",
		"Cookie",
		"Date",
		"Expect",
		"From",
		"Host",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Schedule-Tag-Match",
		"If-Unmodified-Since",
		"Max-Forwards",
		"Origin",
		"Pragma",
		"Referer",
		"User-Agent",
		"Via",
		"Warning":
		return true
	}
	return false
}

func WithOutgoingHeaderMatcher(fn runtime.HeaderMatcherFunc) Option {
	return func(o *options) {
		o.outgoingHeaderMatcher = fn
	}
}

func LogRequests() Option {
	return func(o *options) {
		o.logRequests = true
	}
}

func (g *GrpcServer) WithListener(l net.Listener) *GrpcServer {
	// Crear el muxer cmux
	g.CmuxSrv = newCMux(l)
	// Match connections in order:
	// First grpc, and otherwise HTTP.
	g.GrpcListener = g.CmuxSrv.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// Any significa que no hay coincidencia previa
	// En nuestro caso, no es grpc así que debe ser http.
	g.HttpListener = DefaultServer.CmuxSrv.Match(cmux.Any())
	return g
}

func (g *GrpcServer) SetName(name string) *GrpcServer {
	g.Name = name
	return g
}

func (g *GrpcServer) SetID(id string) *GrpcServer {
	g.ID = id
	return g
}

func (g *GrpcServer) SetPort(port int) *GrpcServer {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Err("[SetPort] cant create listener", "error", err.Error())
		return g
	}
	g.Port = strconv.Itoa(port)
	return g.WithListener(l)
}

func (g *GrpcServer) SetMaxConcurrentStreams(n int) *GrpcServer {
	g.MaxConcurrentStreams = n
	return g
}

func (g *GrpcServer) UseReflection() *GrpcServer {
	g.GrpcReflection = true
	return g
}

func (g *GrpcServer) UseGrpcGw(opts ...Option) *GrpcServer {
	grpcGwOpts := &options{}
	for _, opt := range opts {
		opt(grpcGwOpts)
	}

	g.grpcGwOpts = *grpcGwOpts

	outgoingHeaderMatcher := g.grpcGwOpts.outgoingHeaderMatcher
	g.GrpcGwMux = runtime.NewServeMux(runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcher))
	g.HttpMux.Handle("/", g.indexHandler(g.GrpcGwMux))
	return g
}

func (g *GrpcServer) String() string {
	return fmt.Sprintf("gRPC Microservice\n=================\n  Name: %s\n  ID: %s\n", g.Name, g.ID)
}

func (g *GrpcServer) indexHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g.grpcGwOpts.logRequests {
			log.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		}

		if g.grpcGwOpts.metrics {
			// Registrar llamada REST
			metrics.Counter("rest.requests").Add()
		}

		if g.grpcGwOpts.cors {
			enableCors(&w)

			if r.Method == "OPTIONS" {
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
