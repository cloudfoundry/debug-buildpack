/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package debug

import (
	"github.com/buildpacks/libbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
)

// Dependency indicates that a JVM application should be run with debugging enabled.
const Dependency = "debug"

// Debug represents the debug configuration for a JVM application.
type Debug struct {
	info  buildpack.Info
	layer layers.HelperLayer
}

// Contribute makes the contribution to launch.
func (d Debug) Contribute() error {
	return d.layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
		layer.Logger.LaunchConfiguration("Set $BPL_DEBUG_PORT to configure", "8000")
		layer.Logger.LaunchConfiguration("Set $BPL_DEBUG_SUSPEND to configure", "n")

		return layer.WriteProfile("debug", `PORT=${BPL_DEBUG_PORT:=8000}
SUSPEND=${BPL_DEBUG_SUSPEND:=n}

printf "Debugging enabled on port ${PORT}"

if [[ "${SUSPEND}" = "y" ]]; then
  printf ", suspended on start\n"
else
  printf "\n"
fi

export JAVA_OPTS="${JAVA_OPTS} -agentlib:jdwp=transport=dt_socket,server=y,address=${PORT},suspend=${SUSPEND}"
`)
	}, layers.Launch)
}

// NewDebug creates a new Debug instance. OK is true if build plan contains "debug" dependency, otherwise false.
func NewDebug(build build.Build) (Debug, bool) {
	if !build.Plans.Has(Dependency) {
		return Debug{}, false
	}

	return Debug{build.Buildpack.Info, build.Layers.HelperLayer(Dependency, "Debug")}, true
}
