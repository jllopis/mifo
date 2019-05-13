package middleware

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Shared options for the logger, with a custom gRPC code to log level function.
var grpcZapOptions = []grpc_zap.Option{
	grpc_zap.WithLevels(codeToLevel),
}

type Logger struct {
	FuncGrpcService
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		FuncGrpcService: func(s *grpc.Server) {
			grpc_zap.ReplaceGrpcLogger(logger)
		},
	}
}

// codeToLevel redirects OK to DEBUG level logging instead of INFO
// This is example how you can log several gRPC code results
func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		// It is DEBUG
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

func NewUnaryLogger(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, grpcZapOptions...),
	)
}

func NewStreamLogger(logger *zap.Logger) grpc.StreamServerInterceptor {
	return grpc_middleware.ChainStreamServer(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, grpcZapOptions...),
	)
}
