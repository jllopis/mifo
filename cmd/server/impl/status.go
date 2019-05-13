package impl

import (
	"context"
	"time"

	"bitbucket.org/acbapis/acbapis/common"
	"bitbucket.org/acbapis/acbapis/status"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jllopis/mifo/server/grpc/middleware"
	restMw "github.com/jllopis/mifo/server/rest/middleware"
	"github.com/jllopis/mifo/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// StatusService provides a service that offers basic status functionality:
// - Time -> returns UTC epoch time in nanoseconds precision
// - Version -> returns the current version of the service
// - Status -> returns the current status of the service
type StatusService struct {
	middleware.FuncGrpcService
}

func NewStatusService() *StatusService {
	ss := &StatusService{}
	ss.FuncGrpcService = func(srv *grpc.Server) {
		status.RegisterStatusServiceServer(srv, ss)
	}
	return ss
}

type StatusRestService struct {
	restMw.FuncRestService
}

func NewStatusRestService() *StatusRestService {
	ss := &StatusRestService{}
	ss.FuncRestService = func(ctx context.Context, r *runtime.ServeMux, addr string, opts ...grpc.DialOption) {
		status.RegisterStatusServiceHandlerFromEndpoint(ctx, r, addr, opts)
	}
	return ss
}

// GetServerTime returns the current UTC server time in nanoseconds
func (st *StatusService) GetServerTime(ctx context.Context, empty *common.EmptyMessage) (*status.ServerTimeMessage, error) {
	return &status.ServerTimeMessage{Value: time.Now().UTC().UnixNano()}, nil
}

// GetVersion returns the current API Version. It is a direct mapping from go-version "github.com/hashicorp/go-version.Version
func (st *StatusService) GetVersion(ctx context.Context, empty *common.EmptyMessage) (*common.Version, error) {
	headerVal := "max-age=600, s-maxage=600"
	grpc.SetHeader(ctx, metadata.Pairs("Grpc-Metadata-Cache-Control", headerVal))
	grpc.SetHeader(ctx, metadata.Pairs("Grpc-Metadata-X-App-Cache-Control", headerVal))

	return &common.Version{
		Version:    version.Version,
		APIVersion: version.APIVersion,
		GitCommit:  version.GitCommit,
		BuildDate:  version.BuildDate,
	}, nil
}

func (st *StatusService) GetGlobalServiceStatus(ctx context.Context, empty *common.EmptyMessage) (*status.ServerStatusMessage, error) {
	headerVal := "max-age=600, s-maxage=600"
	grpc.SetHeader(ctx, metadata.Pairs("Grpc-Metadata-Cache-Control", headerVal))
	grpc.SetHeader(ctx, metadata.Pairs("Grpc-Metadata-X-App-Cache-Control", headerVal))

	serverStatus := &status.ServerStatusMessage{
		Status: status.ServerStatus_OK,
	}

	return serverStatus, nil
}
