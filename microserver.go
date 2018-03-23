package mifo

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jllopis/mifo/log"
	"github.com/segmentio/ksuid"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer holds the components that build up a Microserver
type GrpcServer struct {
	Name         string
	Address      string
	Port         string
	ID           string
	Version      string
	CmuxSrv      cmux.CMux
	GrpcListener net.Listener
	GrpcSrv      *grpc.Server
	HttpListener net.Listener
	HttpMux      *http.ServeMux
	HttpSrv      *http.Server
	// Interceptors
	UnaryInter               []grpc.UnaryServerInterceptor
	StreamInter              []grpc.StreamServerInterceptor
	InterceptorInitializers  []InterceptorInitializer
	GrpcServicesRegistrators []ServiceRegistrator
	GrpcReflection           bool
}

// InterceptorInitializer defines an interface with the methods that
// should be implemented to setup an interceptor.
//
// InterceptorInitializers will be added to GrpcServer.InterceptorInitializers slice
// and when the server is created, this array will be traversed and for every item the
// method Init(*GrpcServer) will be called.
//
// This will register the intercepto into the GrpcSrv.
type InterceptorInitializer interface {
	Init(*GrpcServer)
}

type MicroServer interface {
	SetName(string) MicroServer
	SetID(string) MicroServer
	UseUnaryInterceptor(grpc.UnaryServerInterceptor)
	UseStreamInterceptor(grpc.StreamServerInterceptor)
	RegisterInterceptorInitializer(InterceptorInitializer)
	// Handle(http.Handler) error
	Register(ServiceRegistrator)
	// Deregister() error
	Serve() error
	Stop() error
	String() string
}

var (
	DefaultAddress             = ":0"
	DefaultPort                = "9000"
	DefaultName                = "grpc-microserver"
	DefaultVersion             = "1.0.0"
	DefaultID                  = ksuid.New().String()
	DefaultHttpMux             = http.NewServeMux()
	DefaultServer  *GrpcServer = &GrpcServer{
		Name:    DefaultName,
		ID:      DefaultID,
		Port:    DefaultPort,
		CmuxSrv: newCMux(nil),
		HttpMux: DefaultHttpMux,
		HttpSrv: &http.Server{
			Handler: DefaultHttpMux,
		},
	}
)

func NewMicroServer() *GrpcServer {
	if DefaultServer == nil {
		DefaultServer = &GrpcServer{
			Name:    DefaultName,
			Port:    DefaultPort,
			ID:      DefaultID,
			CmuxSrv: newCMux(nil),
			HttpMux: DefaultHttpMux,
			HttpSrv: &http.Server{
				Handler: DefaultHttpMux,
			},
		}
	}
	// Match connections in order:
	// First grpc, and otherwise HTTP.
	DefaultServer.GrpcListener = DefaultServer.CmuxSrv.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// Any significa que no hay coincidencia previa
	// En nuestro caso, no es grpc así que debe ser http.
	DefaultServer.HttpListener = DefaultServer.CmuxSrv.Match(cmux.Any())
	return DefaultServer
}

func newCMux(l net.Listener) cmux.CMux {
	var err error
	if l == nil {
		// Create the main listener.
		l, err = net.Listen("tcp", ":"+DefaultPort)
		if err != nil {
			log.Err("cant create CMux", "error", err.Error())
			return nil
		}
	}

	// Crear el muxer cmux
	return cmux.New(l)
}

// UseUnaryInterceptor adds a unary interceptor to the RPC server
func (g *GrpcServer) UseUnaryInterceptor(inter grpc.UnaryServerInterceptor) {
	g.UnaryInter = append(g.UnaryInter, inter)
}

// UseStreamInterceptor adds a stream interceptor to the RPC server
func (g *GrpcServer) UseStreamInterceptor(inter grpc.StreamServerInterceptor) {
	g.StreamInter = append(g.StreamInter, inter)
}

func (g *GrpcServer) RegisterInterceptorInitializer(i InterceptorInitializer) {
	g.InterceptorInitializers = append(g.InterceptorInitializers, i)
}

type ServiceRegistrator func(*grpc.Server)

func (g *GrpcServer) Register(sr ServiceRegistrator) {
	g.GrpcServicesRegistrators = append(g.GrpcServicesRegistrators, sr)
}

func (g *GrpcServer) Serve() error {
	g.GrpcSrv = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(g.UnaryInter...),
		grpc_middleware.WithStreamServerChain(g.StreamInter...),
	)
	g.initializeInterceptors()

	g.registerGrpcServices()

	if g.GrpcReflection {
		reflection.Register(g.GrpcSrv)
	}

	// Use the muxed listeners for your servers.
	go g.GrpcSrv.Serve(g.GrpcListener)
	if g.HttpSrv != nil {
		go g.HttpSrv.Serve(g.HttpListener)
	}
	// Start serving!
	return g.CmuxSrv.Serve()
}

func (g *GrpcServer) initializeInterceptors() {
	for _, i := range g.InterceptorInitializers {
		i.Init(g)
	}
}

func (g *GrpcServer) registerGrpcServices() {
	for _, s := range g.GrpcServicesRegistrators {
		s(g.GrpcSrv)
	}
}

func (g *GrpcServer) Stop() {
	log.Info("Stopping services...")
	g.GrpcSrv.GracefulStop()
	g.HttpSrv.Shutdown(context.Background())
}
