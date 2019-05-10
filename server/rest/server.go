package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jllopis/mifo/logger"
	"github.com/jllopis/mifo/option"
	"google.golang.org/grpc"
)

type RestServer struct {
	httpSrv  *http.Server
	httpMux  *http.ServeMux
	gwMux    *runtime.ServeMux
	listener net.Listener
}

func New(o *option.MsOptions) *RestServer {
	gm := runtime.NewServeMux(o.RestOptions...)
	for _, service := range o.RestServices {
		service.RegisterService(context.Background(), gm, o.Address, grpc.WithInsecure())
	}
	mux := http.NewServeMux()
	mux.Handle("/", o.RestMiddleware.Then(gm))

	return &RestServer{
		gwMux:   gm,
		httpMux: mux,
		httpSrv: &http.Server{
			Handler: mux,
		},
	}
}

// RunServer runs HTTP/REST gateway
func (r *RestServer) Serve(listen net.Listener) {
	if listen == nil {
		logger.Log.Error("rest.New got a nil listener")
	}
	r.listener = listen

	logger.Log.Info("starting HTTP/REST gateway on " + listen.Addr().String())
	r.httpSrv.Serve(listen)
}

// func (r *RestServer) GetRuntimeMux() *runtime.ServeMux {
// 	return r.gwMux
// }

func (r *RestServer) GetHttpMux() *http.ServeMux {
	return r.httpMux
}

func (r *RestServer) Shutdown() {
	fmt.Println("Shutting down http server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r.httpSrv.Shutdown(ctx)
	// r.listener.Close()
	fmt.Println("http server server shot down")
}
