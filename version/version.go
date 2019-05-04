package version

const (
	// Name of the service
	Name = "MifoProject"
	// APIVersion of the API this service provide
	APIVersion = "v1"
)

var (
	// Version is the current echo version in SemVer format
	Version = "v0.0.1"
	// BuildDate represents the date this service was built
	BuildDate string
	// GitCommit is the hash of the git commit
	GitCommit string
)
