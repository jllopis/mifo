package middleware

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Reflect struct {
	FuncGrpcService
}

func NewReflection() *Reflect {
	return &Reflect{
		FuncGrpcService: func(s *grpc.Server) {
			reflection.Register(s)
		},
	}
}
