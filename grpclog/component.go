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

// ComponentData records the settings for a component.
type ComponentData struct {
	name      string
	verbosity int
	level     int
}

const envName = "GRPC_GO_LOG_LEVEL"
const sentinel = math.MinInt32

// log levels
const (
	levelInfo    = 0
	levelWarning = 1
	levelError   = 2
)

// grpc components
var (
	Core       = ComponentData{"CORE", 0, 0}
	Alts       = ComponentData{"ALTS", 0, 0}
	Grpclb     = ComponentData{"BALANCER_GRPCLB", 0, 0}
	Rls        = ComponentData{"BALANCER_RLS", 0, 0}
	RoundRobin = ComponentData{"BALANCER_ROUNDROBIN", 0, 0}
	BinaryLog  = ComponentData{"BINARYLOG", 0, 0}
	Channelz   = ComponentData{"CHANNELZ", 0, 0}
	DNS        = ComponentData{"RESOLVER_DNS", 0, 0}
	Transport  = ComponentData{"TRANSPORT", 0, 0}
	Xds        = ComponentData{"XDS", 0, 0}
)

var grpcComponents = []*ComponentData{
	&Core,
	&Alts,
	&Grpclb,
	&Rls,
	&RoundRobin,
	&BinaryLog,
	&Channelz,
	&DNS,
	&Transport,
	&Xds,
}
var environmentVars = map[string]*ComponentData{}
var prefixVars = map[string]*ComponentData{}

// init creates all grpc components and any components named in the environment
// variable and applies all settings specified in the environment variable.
// Notably, settings for all components with a given prefix are applies before
// regular component settings.
func init() {
	// Pull environment variable data and put in environmentVars and prefixVars
	if envVar, ok := os.LookupEnv(envName); ok && len(envVar) > 0 {
		varList := strings.Split(envVar, ",")
		for _, varPair := range varList {
			varPairList := strings.Split(varPair, ":")
			if len(varPairList) != 2 {
				fmt.Fprintf(os.Stderr, "error: could not parse %v value '%v', unrecognized key-value pair '%v'\n", envName, envVar, varPair)
				os.Exit(1)
			}
			if cData, ok := parseVar(varPairList[0], varPairList[1]); ok {
				if varPrefix, ok := getPrefix(varPairList[0]); ok {
					prefixVars[varPrefix] = &cData
				} else {
					environmentVars[varPairList[0]] = &cData
				}
			} else {
				fmt.Fprintf(os.Stderr, "error: could not parse %v value '%v', unrecognized value '%v'\n", envName, envVar, varPairList[1])
				os.Exit(1)
			}
		}
	}
	// Apply environment variables to grpcComponents
	for _, cData := range grpcComponents {
		// Apply prefix settings
		for prefix, pData := range prefixVars {
			if strings.HasPrefix(cData.name, prefix) {
				cData.apply(pData)
			}
		}
		// Apply specific settings
		if vData, ok := environmentVars[cData.name]; ok {
			cData.apply(vData)
		}
	}
}

// apply applies the parameter componentData to the receiver ComponentData.
// Sentinel values in the parameter will not be applied.
func (c *ComponentData) apply(applyData *ComponentData) {
	if applyData.verbosity != sentinel {
		c.verbosity = applyData.verbosity
	}
	if applyData.level != sentinel {
		c.level = applyData.level
	}
}

// parseVar parses a key-value pair from the environment variable. The resulting
// ComponentData will have a sentinel value for the verbosity if it is
// unspecified in the value. Second return value is false if there was an error
// parsing.
func parseVar(key string, value string) (ComponentData, bool) {
	result := ComponentData{key, sentinel, sentinel}
	switch {
	case value == "INFO":
		result.level = levelInfo
	case value == "WARNING":
		result.level = levelWarning
	case value == "ERROR":
		result.level = levelError
	case strings.HasPrefix(value, "INFO_"):
		result.level = levelInfo
		vStr := strings.TrimPrefix(value, "INFO_")
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

// getPrefix Gets the prefix if s has a wildcard. If s does not have a wildcard,
// returns "", false.
func getPrefix(s string) (string, bool) {
	if strings.HasSuffix(s, "*") {
		return strings.TrimSuffix(s, "*"), true
	}
	return "", false
}

// InfoDepth performs an info log of args at depth to the component, conditioned
// on the component's log level.
func (c *ComponentData) InfoDepth(depth int, args ...interface{}) {
	if c.level > levelInfo {
		return
	}
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.InfoDepth(depth, args...)
}

// WarningDepth performs a warning log of args at depth to the component,
// conditioned on the component's log level.
func (c *ComponentData) WarningDepth(depth int, args ...interface{}) {
	if c.level > levelWarning {
		return
	}
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.WarningDepth(depth, args...)
}

// ErrorDepth performs an error log of args at depth to the component.
func (c *ComponentData) ErrorDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.ErrorDepth(depth, args...)
}

// FatalDepth performs a fatal log of args at depth to the component and then
// exits the application in accordance with the logger's fatal behavior.
func (c *ComponentData) FatalDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.FatalDepth(depth, args...)
}

// Info performs an InfoDepth log at depth 0.
func (c *ComponentData) Info(args ...interface{}) {
	c.InfoDepth(0, args...)
}

// Warning performs a WarningDepth log at depth 0.
func (c *ComponentData) Warning(args ...interface{}) {
	c.WarningDepth(0, args...)
}

// Error performs an ErrorDepth log at depth 0.
func (c *ComponentData) Error(args ...interface{}) {
	c.ErrorDepth(0, args...)
}

// Fatal performs a FatalDepth log at depth 0.
func (c *ComponentData) Fatal(args ...interface{}) {
	c.FatalDepth(0, args...)
}

// Infof formats the string and performs and InfoDepth log at depth 0.
func (c *ComponentData) Infof(format string, args ...interface{}) {
	c.InfoDepth(0, fmt.Sprintf(format, args...))
}

// Warningf formats the string and performs and WarningDepth log at depth 0.
func (c *ComponentData) Warningf(format string, args ...interface{}) {
	c.WarningDepth(0, fmt.Sprintf(format, args...))
}

// Errorf formats the string and performs and ErrorDepth log at depth 0.
func (c *ComponentData) Errorf(format string, args ...interface{}) {
	c.ErrorDepth(0, fmt.Sprintf(format, args...))
}

// Fatalf formats the string and performs and FatalDepth log at depth 0.
func (c *ComponentData) Fatalf(format string, args ...interface{}) {
	c.FatalDepth(0, fmt.Sprintf(format, args...))
}

// V reports whether thbe verbosity level of the component is at least l.
func (c *ComponentData) V(l int) bool {
	return c.verbosity >= l
}

// Component creates a new component and returns its identifier used for
// logging. If a component with the name already exists, nothing will be created
// and its identifier will be returned.
func Component(componentName string) ComponentData {
	c := ComponentData{componentName, 0, 0}
	// Apply prefix settings
	for prefix, pData := range prefixVars {
		if strings.HasPrefix(c.name, prefix) {
			c.apply(pData)
		}
	}
	// Apply specific settings
	if vData, ok := environmentVars[c.name]; ok {
		c.apply(vData)
	}
	return c
}
