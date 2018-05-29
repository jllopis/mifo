package impl

import (
	"context"
	"time"

	"bitbucket.org/acbapis/acbapis/common"
	"bitbucket.org/acbapis/acbapis/status"
	"github.com/jllopis/mifo/version"
)

// StatusService provides a service that offers basic status functionality:
// - Time -> returns UTC epoch time in nanoseconds precision
// - Version -> returns the current version of the service
// - Status -> returns the current status of the service
type StatusService struct{}

// GetServerTime returns the current UTC server time in nanoseconds
func (st *StatusService) GetServerTime(ctx context.Context, empty *common.EmptyMessage) (*status.ServerTimeMessage, error) {
	return &status.ServerTimeMessage{Value: time.Now().UTC().UnixNano()}, nil
}

// GetVersion returns the current API Version. It is a direct mapping from github.com/coreos/go-semver/semver.Version
func (st *StatusService) GetVersion(ctx context.Context, empty *common.EmptyMessage) (*common.Version, error) {
	return &common.Version{
		Version:    version.Version,
		APIVersion: version.APIVersion,
		GitCommit:  version.GitCommit,
		BuildDate:  version.BuildDate,
	}, nil
}

func (st *StatusService) GetGlobalServiceStatus(ctx context.Context, empty *common.EmptyMessage) (*status.ServerStatusMessage, error) {
	serverStatus := &status.ServerStatusMessage{
		Status: status.ServerStatus_OK,
	}

	return serverStatus, nil
}
