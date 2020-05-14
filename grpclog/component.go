/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpclog

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/internal/grpclog"
)

// Component is a unique identifier for a component. Created with NewComponent.
type Component string

// componentData records the settings for a component.
type componentData struct {
	verbosity int
	level     int
}

const envName = "GRPC_GO_LOG_LEVEL"
const sentinel = math.MinInt32

// grpc components
const (
	CORE       = Component("CORE")
	ALTS       = Component("ALTS")
	GRPCLB     = Component("BALANCER_GRPCLB")
	RLS        = Component("BALANCER_RLS")
	ROUNDROBIN = Component("BALANCER_ROUNDROBIN")
	BINARYLOG  = Component("BINARYLOG")
	CHANNELZ   = Component("CHANNELZ")
	DNS        = Component("RESOLVER_DNS")
	TRANSPORT  = Component("TRANSPORT")
	XDS        = Component("XDS")
)

var components = map[Component]*componentData{}
var grpcComponents = []Component{
	CORE,
	ALTS,
	GRPCLB,
	RLS,
	ROUNDROBIN,
	BINARYLOG,
	CHANNELZ,
	DNS,
	TRANSPORT,
	XDS,
}
var environmentVars = map[string]*componentData{}
var prefixVars = map[string]*componentData{}

// init creates all grpc components and any components named in the environment variable and applies all settings
// specified in the environment variable. Notably, settings for all components with a given prefix are applies before
// regular component settings.
func init() {
	// Initialize the grpc components
	for _, c := range grpcComponents {
		components[c] = &componentData{0, 0}
	}
	// Pull environment variable data
	if envVar, ok := os.LookupEnv(envName); ok && len(envVar) > 0 {
		varList := strings.Split(envVar, ",")
		for _, varPair := range varList {
			varPairList := strings.Split(varPair, ":")
			if len(varPairList) != 2 {
				fmt.Fprintf(os.Stderr, "error: could not parse %v value '%v', unrecognized key-value pair '%v'\n", envName, envVar, varPair)
				os.Exit(1)
			}
			if cData, ok := parseVar(varPairList[1]); ok {
				if varPrefix, ok := getPrefix(varPairList[0]); ok {
					prefixVars[varPrefix] = &componentData{cData.verbosity, cData.level}
				} else {
					environmentVars[varPairList[0]] = &componentData{cData.verbosity, cData.level}
				}
			} else {
				fmt.Fprintf(os.Stderr, "error: could not parse %v value '%v', unrecognized value '%v'\n", envName, envVar, varPairList[1])
				os.Exit(1)
			}
		}
	}
	// Create any components named in the environment variable that do not exist
	for name := range environmentVars {
		c := Component(name)
		if _, ok := components[c]; !ok {
			components[c] = &componentData{0, 0}
		}
	}
	// Apply prefixes
	for prefix, pData := range prefixVars {
		for c, cData := range components {
			if strings.HasPrefix(string(c), prefix) {
				cData.apply(pData)
			}
		}
	}
	// Apply non-prefixes
	for name, vData := range environmentVars {
		c := Component(name)
		if cData, ok := components[c]; ok {
			cData.apply(vData)
		}
	}
}

// apply applies the parameter componentData to the receiver componentData.
// Sentinel values in the parameter will not be applied.
func (cData *componentData) apply(applyData *componentData) {
	if applyData.verbosity != sentinel {
		cData.verbosity = applyData.verbosity
	}
	if applyData.level != sentinel {
		cData.level = applyData.level
	}
}

// parseVar parses the value in a key-value pair from the environment variable.
// The resulting componentData will have a sentinel value for the verbosity if it is unspecified in the value.
// Second return value is false if there was an error parsing.
func parseVar(s string) (componentData, bool) {
	result := componentData{sentinel, sentinel}
	switch {
	case s == "INFO":
		result.level = 0
	case s == "WARNING":
		result.level = 1
	case s == "ERROR":
		result.level = 2
	case strings.HasPrefix(s, "INFO_"):
		result.level = 0
		vStr := strings.TrimPrefix(s, "INFO_")
		v, err := strconv.Atoi(vStr)
		if err != nil {
			return result, false
		}
		result.verbosity = v
	default:
		return result, false
	}
	return result, true
}

// getPrefix Gets the prefix if s has a wildcard.
// If s does not have a wildcard, returns "", false.
func getPrefix(s string) (string, bool) {
	if strings.HasSuffix(s, "*") {
		return strings.TrimSuffix(s, "*"), true
	}
	return "", false
}

// InfoDepth performs an info log of args at depth to the component, conditioned on the component's log level.
func (c Component) InfoDepth(depth int, args ...interface{}) {
	if cData, ok := components[c]; ok {
		if cData.level > 0 {
			return
		}
		args = append([]interface{}{"[" + string(c) + "]"}, args...)
		grpclog.InfoDepth(depth, args...)
	} else {
		grpclog.Logger.Error("component", c, "not identified")
	}
}

// WarningDepth performs a warning log of args at depth to the component, conditioned on the component's log level.
func (c Component) WarningDepth(depth int, args ...interface{}) {
	if cData, ok := components[c]; ok {
		if cData.level > 1 {
			return
		}
		args = append([]interface{}{"[" + string(c) + "]"}, args...)
		grpclog.WarningDepth(depth, args...)
	} else {
		grpclog.Logger.Error("component", c, "not identified")
	}
}

// ErrorDepth performs an error log of args at depth to the component.
func (c Component) ErrorDepth(depth int, args ...interface{}) {
	if _, ok := components[c]; ok {
		args = append([]interface{}{"[" + string(c) + "]"}, args...)
		grpclog.ErrorDepth(depth, args...)
	} else {
		grpclog.Logger.Error("component", c, "not identified")
	}
}

// FatalDepth performs a fatal log of args at depth to the component and then exits the application in accordance with the logger's fatal behavior.
func (c Component) FatalDepth(depth int, args ...interface{}) {
	if _, ok := components[c]; ok {
		args = append([]interface{}{"[" + string(c) + "]"}, args...)
		grpclog.FatalDepth(depth, args...)
	} else {
		grpclog.Logger.Error("component", c, "not identified")
		grpclog.Logger.Fatal(args)
	}
}

// Info performs an InfoDepth log at depth 0.
func (c Component) Info(args ...interface{}) {
	c.InfoDepth(0, args...)
}

// Warning performs a WarningDepth log at depth 0.
func (c Component) Warning(args ...interface{}) {
	c.WarningDepth(0, args...)
}

// Error performs an ErrorDepth log at depth 0.
func (c Component) Error(args ...interface{}) {
	c.ErrorDepth(0, args...)
}

// Fatal performs a FatalDepth log at depth 0.
func (c Component) Fatal(args ...interface{}) {
	c.FatalDepth(0, args...)
}

// Infof formats the string and performs and InfoDepth log at depth 0.
func (c Component) Infof(format string, args ...interface{}) {
	c.InfoDepth(0, fmt.Sprintf(format, args...))
}

// Warningf formats the string and performs and WarningDepth log at depth 0.
func (c Component) Warningf(format string, args ...interface{}) {
	c.WarningDepth(0, fmt.Sprintf(format, args...))
}

// Errorf formats the string and performs and ErrorDepth log at depth 0.
func (c Component) Errorf(format string, args ...interface{}) {
	c.ErrorDepth(0, fmt.Sprintf(format, args...))
}

// Fatalf formats the string and performs and FatalDepth log at depth 0.
func (c Component) Fatalf(format string, args ...interface{}) {
	c.FatalDepth(0, fmt.Sprintf(format, args...))
}

// V reports whether thbe verbosity level of the component is at least l.
func (c Component) V(l int) bool {
	if cData, ok := components[c]; ok {
		return cData.verbosity >= l
	}
	grpclog.Logger.Error("component", c, "not identified")
	return true
}

// NewComponent creates a new component and returns its identifier used for logging.
// If a component with the name already exists, nothing will be created and its identifier will be returned.
func NewComponent(componentName string) Component {
	c := Component(componentName)
	if _, ok := components[c]; !ok {
		// The component does not exist, so create it
		// A component may already exist upon execution of this function if its name was encountered in the environment variable
		cData := componentData{0, 0}
		components[c] = &cData
		// Apply any prefix settings to it
		for prefix, pData := range prefixVars {
			if strings.HasPrefix(componentName, prefix) {
				cData.apply(pData)
			}
		}
		// Apply non-prefix settings
		if vData, ok := environmentVars[string(c)]; ok {
			cData.apply(vData)
		}
	}
	return c
}
