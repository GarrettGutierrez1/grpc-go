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

// componentData records the settings for a component.
type componentData struct {
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

var environmentVars = map[string]*componentData{}
var prefixVars = map[string]*componentData{}
var cache = map[string]*componentData{}

// init creates all grpc components and any components named in the environment
// variable and applies all settings specified in the environment variable.
// Notably, settings for all components with a given prefix are applies before
// regular component settings.
func init() {
	// Pull environment variable data and put in environmentVars and prefixVars
	if envVar, ok := os.LookupEnv(envName); ok && len(envVar) > 0 {
		environmentVars, prefixVars = parseEnvironmentVar(envVar)
	}
}

// parseEnvironmentVar parses an environment variable string and pulls the component settings data.
func parseEnvironmentVar(envVar string) (map[string]*componentData, map[string]*componentData) {
	envVars := map[string]*componentData{}
	preVars := map[string]*componentData{}
	varList := strings.Split(envVar, ",")
	for _, varPair := range varList {
		varPairList := strings.Split(varPair, ":")
		if len(varPairList) != 2 {
			fmt.Fprintf(os.Stderr, "error: could not parse '%v' value '%v', unrecognized key-value pair '%v'\n", envName, envVar, varPair)
			continue
		}
		if cData, ok := parseVar(varPairList[0], varPairList[1]); ok {
			if varPrefix, ok := getPrefix(varPairList[0]); ok {
				preVars[varPrefix] = &cData
			} else {
				envVars[varPairList[0]] = &cData
			}
		} else {
			fmt.Fprintf(os.Stderr, "error: could not parse '%v' value '%v', unrecognized value '%v'\n", envName, envVar, varPairList[1])
		}
	}
	return envVars, preVars
}

// apply applies the parameter componentData to the receiver componentData.
// Sentinel values in the parameter will not be applied.
func (c *componentData) apply(applyData *componentData) {
	if applyData.verbosity != sentinel {
		c.verbosity = applyData.verbosity
	}
	if applyData.level != sentinel {
		c.level = applyData.level
	}
}

// parseVar parses a key-value pair from the environment variable. The resulting
// componentData will have a sentinel value for the verbosity if it is
// unspecified in the value. Second return value is false if there was an error
// parsing.
func parseVar(key string, value string) (componentData, bool) {
	result := componentData{key, sentinel, sentinel}
	value = strings.ToUpper(value)
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
func (c *componentData) InfoDepth(depth int, args ...interface{}) {
	if c.level > levelInfo {
		return
	}
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.InfoDepth(depth, args...)
}

// WarningDepth performs a warning log of args at depth to the component,
// conditioned on the component's log level.
func (c *componentData) WarningDepth(depth int, args ...interface{}) {
	if c.level > levelWarning {
		return
	}
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.WarningDepth(depth, args...)
}

// ErrorDepth performs an error log of args at depth to the component.
func (c *componentData) ErrorDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.ErrorDepth(depth, args...)
}

// FatalDepth performs a fatal log of args at depth to the component and then
// exits the application in accordance with the logger's fatal behavior.
func (c *componentData) FatalDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	grpclog.FatalDepth(depth, args...)
}

// Info performs an InfoDepth log at depth 0.
func (c *componentData) Info(args ...interface{}) {
	c.InfoDepth(0, args...)
}

// Warning performs a WarningDepth log at depth 0.
func (c *componentData) Warning(args ...interface{}) {
	c.WarningDepth(0, args...)
}

// Error performs an ErrorDepth log at depth 0.
func (c *componentData) Error(args ...interface{}) {
	c.ErrorDepth(0, args...)
}

// Fatal performs a FatalDepth log at depth 0.
func (c *componentData) Fatal(args ...interface{}) {
	c.FatalDepth(0, args...)
}

// Infof formats the string and performs and InfoDepth log at depth 0.
func (c *componentData) Infof(format string, args ...interface{}) {
	c.InfoDepth(0, fmt.Sprintf(format, args...))
}

// Warningf formats the string and performs and WarningDepth log at depth 0.
func (c *componentData) Warningf(format string, args ...interface{}) {
	c.WarningDepth(0, fmt.Sprintf(format, args...))
}

// Errorf formats the string and performs and ErrorDepth log at depth 0.
func (c *componentData) Errorf(format string, args ...interface{}) {
	c.ErrorDepth(0, fmt.Sprintf(format, args...))
}

// Fatalf formats the string and performs and FatalDepth log at depth 0.
func (c *componentData) Fatalf(format string, args ...interface{}) {
	c.FatalDepth(0, fmt.Sprintf(format, args...))
}

// Infoln performs an InfoDepth log at depth 0.
func (c *componentData) Infoln(args ...interface{}) {
	c.Info(args...)
}

// Warningln performs a WarningDepth log at depth 0.
func (c *componentData) Warningln(args ...interface{}) {
	c.Warning(args...)
}

// Errorln performs an ErrorDepth log at depth 0.
func (c *componentData) Errorln(args ...interface{}) {
	c.Error(args...)
}

// Fatalln performs a FatalDepth log at depth 0.
func (c *componentData) Fatalln(args ...interface{}) {
	c.Fatal(args...)
}

// V reports whether thbe verbosity level of the component is at least l.
func (c *componentData) V(l int) bool {
	return c.verbosity >= l
}

// Component creates a new component and returns it for logging. If a component
// with the name already exists, nothing will be created and it will be
// returned.
func Component(componentName string) DepthLoggerV2 {
	if cData, ok := cache[componentName]; ok {
		return cData
	}
	c := componentData{componentName, 0, 0}
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
	cache[componentName] = &c
	return &c
}
