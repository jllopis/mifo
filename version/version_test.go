package version

import (
	"testing"

	"github.com/coreos/go-semver/semver"
)

type testValues struct {
	value      string
	Version    string
	Major      int64
	Minor      int64
	Patch      int64
	PreRelease semver.PreRelease
	Metadata   string
	GitCommit  string
}

var (
	// {string, [version, major, minor, Patch, PreRelease, Metadata, GitCommit}
	tests = []testValues{
		{"0.0.1-212-gbf8a411", "0.0.1-212-gbf8a411", 0, 0, 1, "212-gbf8a411", "", "f8a411"},
		{"0.0.1", "0.0.1", 0, 0, 1, "", "", ""},
		{"0.0.1-dirty", "0.0.1-dirty", 0, 0, 1, "dirty", "", ""},
		{"0.0.1-212-gbf8a411-dirty", "0.0.1-212-gbf8a411-dirty", 0, 0, 1, "212-gbf8a411-dirty", "", "f8a411"},
	}
)

func TestSemVer(t *testing.T) {
	for _, test := range tests {
		version := test.value
		semver, err := semver.NewVersion(version)
		if err != nil {
			t.Error("error creating semver.NewVersion", err)
		}
		if semver.Major != test.Major {
			t.Error("For", test.value,
				"Got", semver.Major,
				"Expected", test.Major)
		}
		if semver.Minor != test.Minor {
			t.Error("For", test.value,
				"Got", semver.Minor,
				"Expected", test.Minor)
		}
		if semver.Patch != test.Patch {
			t.Error("For", test.value,
				"Got", semver.Patch,
				"Expected", test.Patch)
		}
		if semver.PreRelease != test.PreRelease {
			t.Error("For", test.value,
				"Got", semver.PreRelease,
				"Expected", test.PreRelease)
		}
		if semver.Metadata != test.Metadata {
			t.Error("For", test.value,
				"Got", semver.Metadata,
				"Expected", test.Metadata)
		}
	}
}
