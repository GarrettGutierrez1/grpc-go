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

// Package grpclog (internal) defines depth logging for grpc.
package grpclog

import (
	"os"
	"strconv"
)

// Component is an alias for string, used to identify a component.
type Component string

// Logger is the logger used for the non-depth log functions.
var Logger LoggerV2

// DepthLogger is the logger used for the depth log functions.
var DepthLogger DepthLoggerV2

// Component identifiers.
const (
	Channelz Component = "GRPC_COMPONENT_CHANNELZ"
	Other    Component = "GRPC_COMPONENT_OTHER"
)

// defaultVerbosity specifies the defaul verbosity for any component that uses
// this as the verbosity in its componentData definition.
const defaultVerbosity = 0

// componentData specifies the data for a component relevant to logging.
type componentData struct {
	tag       string
	verbosity int
}

// components maps component IDs to component logging data.
var components = map[Component]componentData{
	Channelz: componentData{"[channelz]", defaultVerbosity},
	Other:    componentData{"[other]", defaultVerbosity},
}

// init uses environment variables to set the verbosity of each component.
func init() {
	for c, cData := range components {
		if varString, ok := os.LookupEnv(string(c)); ok {
			if v, err := strconv.Atoi(varString); err == nil {
				cData.verbosity = v
			}
		}
	}
}

// CInfo info logs for a specified component.
// Log will only occur is v is <= the component's verbosity.
func CInfo(c Component, v int, d int, args ...interface{}) {
	if cData, ok := components[c]; ok {
		args = append([]interface{}{cData.tag}, args...)
		if v <= cData.verbosity {
			InfoDepth(d, args)
		}
	} else {
		Logger.Error("component", c, "not identified")
	}
}

// CWarning warning logs for a specified component.
func CWarning(c Component, d int, args ...interface{}) {
	if cData, ok := components[c]; ok {
		args = append([]interface{}{cData.tag}, args...)
		WarningDepth(d, args)
	} else {
		Logger.Error("component", c, "not identified")
	}
}

// CError error logs for a specified component.
func CError(c Component, d int, args ...interface{}) {
	if cData, ok := components[c]; ok {
		args = append([]interface{}{cData.tag}, args...)
		ErrorDepth(d, args)
	} else {
		Logger.Error("component", c, "not identified")
	}
}

// InfoDepth logs to the INFO log at the specified depth.
func InfoDepth(depth int, args ...interface{}) {
	if DepthLogger != nil {
		DepthLogger.InfoDepth(depth, args...)
	} else {
		Logger.Info(args...)
	}
}

// WarningDepth logs to the WARNING log at the specified depth.
func WarningDepth(depth int, args ...interface{}) {
	if DepthLogger != nil {
		DepthLogger.WarningDepth(depth, args...)
	} else {
		Logger.Warning(args...)
	}
}

// ErrorDepth logs to the ERROR log at the specified depth.
func ErrorDepth(depth int, args ...interface{}) {
	if DepthLogger != nil {
		DepthLogger.ErrorDepth(depth, args...)
	} else {
		Logger.Error(args...)
	}
}

// FatalDepth logs to the FATAL log at the specified depth.
func FatalDepth(depth int, args ...interface{}) {
	if DepthLogger != nil {
		DepthLogger.FatalDepth(depth, args...)
	} else {
		Logger.Fatal(args...)
	}
}

// LoggerV2 does underlying logging work for grpclog.
// This is a copy of the LoggerV2 defined in the external grpclog package. It
// is defined here to avoid a circular dependency.
type LoggerV2 interface {
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
	Infoln(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warning(args ...interface{})
	// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
	Warningln(args ...interface{})
	// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warningf(format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	Errorln(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(args ...interface{})
	// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalln(args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalf(format string, args ...interface{})
	// V reports whether verbosity level l is at least the requested verbose level.
	V(l int) bool
}

// DepthLoggerV2 logs at a specified call frame. If a LoggerV2 also implements
// DepthLoggerV2, the below functions will be called with the appropriate stack
// depth set for trivial functions the logger may ignore.
// This is a copy of the DepthLoggerV2 defined in the external grpclog package.
// It is defined here to avoid a circular dependency.
//
// This API is EXPERIMENTAL.
type DepthLoggerV2 interface {
	// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Print.
	InfoDepth(depth int, args ...interface{})
	// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Print.
	WarningDepth(depth int, args ...interface{})
	// ErrorDetph logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Print.
	ErrorDepth(depth int, args ...interface{})
	// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Print.
	FatalDepth(depth int, args ...interface{})
}
