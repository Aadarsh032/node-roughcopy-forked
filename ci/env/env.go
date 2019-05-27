/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package env

import (
	"fmt"
	"os"
	"strconv"
)

const devReleaseVersion = "0.0.0-dev"
const ppaDevReleaseVersion = "0.0.0"

type envVar struct {
	key BuildVar
	val string
}

// BuildVar env variable for CI
type BuildVar string

const (
	// ReleaseBuild indicates release build
	ReleaseBuild = BuildVar("RELEASE_BUILD")

	// MasterBuild indicates master branch (non-release) build
	MasterBuild = BuildVar("MASTER_BUILD")

	// PrBuild indicates pull-request build
	PrBuild = BuildVar("PR_BUILD")

	// BuildVersion stores build version
	BuildVersion = BuildVar("BUILD_VERSION")

	// PpaVersion stores build version for PPA
	PpaVersion = BuildVar("PPA_VERSION")

	// BuildNumber stores CI build number
	BuildNumber = BuildVar("BUILD_NUMBER")

	// GithubOwner stores github repository's owner
	GithubOwner = BuildVar("GITHUB_OWNER")

	// GithubRepository stores github repository name
	GithubRepository = BuildVar("GITHUB_REPO")

	// GithubDevBuildRepository stores github repository name for dev build
	GithubDevBuildRepository = BuildVar("GITHUB_REPO_DEV_BUILD")

	// GithubApiToken is used for accessing github API
	GithubApiToken = BuildVar("GITHUB_API_TOKEN")
)

// GenerateEnvFile for sourcing in other stages
func GenerateEnvFile() error {
	vars := []envVar{
		{ReleaseBuild, strconv.FormatBool(isReleaseBuild())},
		{MasterBuild, strconv.FormatBool(isMasterBuild())},
		{PrBuild, strconv.FormatBool(isPullRequestBuild())},
		{BuildVersion, buildVersion()},
		{PpaVersion, ppaVersion()},
		{BuildNumber, os.Getenv(string(BuildNumber))},
		{GithubOwner, os.Getenv(string(GithubOwner))},
		{GithubRepository, os.Getenv(string(GithubRepository))},
		{GithubDevBuildRepository, os.Getenv(string(GithubDevBuildRepository))},
	}
	return writeEnvVars(vars)
}

func isReleaseBuild() bool {
	return releaseVersion() != ""
}

func isMasterBuild() bool {
	return os.Getenv("BUILD_BRANCH") == "master" && !isReleaseBuild()
}

func isPullRequestBuild() bool {
	return !isMasterBuild() && !isReleaseBuild()
}

func releaseVersion() string {
	return os.Getenv("BUILD_TAG")
}

func buildVersion() string {
	if isReleaseBuild() {
		return releaseVersion()
	}
	return devReleaseVersion + "_" + os.Getenv("BUILD_COMMIT")
}

func ppaVersion() string {
	if isReleaseBuild() {
		return releaseVersion()
	}
	return ppaDevReleaseVersion
}

func writeEnvVars(vars []envVar) error {
	_ = os.Mkdir("./build", 0755)
	file, err := os.Create("./build/env.sh")
	if err != nil {
		return err
	}
	defer file.Close()
	for _, v := range vars {
		_, err := fmt.Fprintf(file, "export %s=%s;\n", v.key, v.val)
		if err != nil {
			return err
		}
	}
	return nil
}
