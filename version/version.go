// Copyright 2017 Joan Llopis. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// version package holds the latest version information following the
// semver spec (http://semver.org)
package version

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

var (
	// Name of the service
	Name = "unknown"
	// Version of the service
	Version = "0.0.1"
	// APIVersion of the API this service provide
	APIVersion = "unknown"
	// SemVer is the semantic version object
	SemVer *semver.Version
	// BuildDate represents the date this service was built
	BuildDate string
	// GitCommit is the hash of the git commit
	GitCommit string
)

func init() {
	v, err := semver.NewVersion(Version)
	if err == nil {
		SemVer = v
	} else {
		fmt.Printf("semver error: %s\n", err)
	}
}
