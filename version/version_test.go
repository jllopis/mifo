package version

import (
	"testing"

	semver "github.com/hashicorp/go-version"
)

type testValues struct {
	value      string
	Version    string
	Major      int64
	Minor      int64
	Patch      int64
	PreRelease string
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
		{"0.0.1-212-gbf8a411", "v0.0.1-212-gbf8a411", 0, 0, 1, "212-gbf8a411", "", "f8a411"},
		{"v0.0.1", "0.0.1", 0, 0, 1, "", "", ""},
		{"v0.0.1-dirty", "v0.0.1-dirty", 0, 0, 1, "dirty", "", ""},
	}
)

func TestSemVer(t *testing.T) {
	for _, test := range tests {
		version := test.value
		semver, err := semver.NewVersion(version)
		if err != nil {
			t.Error("error creating semver.NewVersion", err)
		}
		semverSegments := semver.Segments64()
		// testSegments := semver.Segments64()
		if semverSegments[0] != test.Major {
			t.Error("For", test.value,
				"Got", semverSegments[0],
				"Expected", test.Major)
		}
		if semverSegments[1] != test.Minor {
			t.Error("For", test.value,
				"Got", semverSegments[1],
				"Expected", test.Minor)
		}
		if semverSegments[2] != test.Patch {
			t.Error("For", test.value,
				"Got", semverSegments[2],
				"Expected", test.Patch)
		}
		if semver.Prerelease() != test.PreRelease {
			t.Error("For", test.value,
				"Got", semver.Prerelease(),
				"Expected", test.PreRelease)
		}
		if semver.Metadata() != test.Metadata {
			t.Error("For", test.value,
				"Got", semver.Metadata(),
				"Expected", test.Metadata)
		}
	}
}
