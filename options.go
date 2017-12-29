package mifo

import (
	"fmt"
	"net"
	"strconv"

	"github.com/jllopis/mifo/log"
	"github.com/soheilhy/cmux"
)

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

func (g *GrpcServer) UseReflection() *GrpcServer {
	g.GrpcReflection = true
	return g
}

func (g *GrpcServer) String() string {
	return fmt.Sprintf("gRPC Microservice\n=================\n  Name: %s\n  ID: %s\n", g.Name, g.ID)
}
