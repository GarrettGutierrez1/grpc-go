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

package channelz

import (
	"fmt"

	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.Component("CHANNELZ")

// Info logs and adds a trace event if channelz is on.
func Info(id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtINFO,
		})
	} else {
		logger.InfoDepth(1, args...)
	}
}

// Infof logs and adds a trace event if channelz is on.
func Infof(id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtINFO,
		})
	} else {
		logger.InfoDepth(1, msg)
	}
}

// Warning logs and adds a trace event if channelz is on.
func Warning(id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtWarning,
		})
	} else {
		logger.WarningDepth(1, args...)
	}
}

// Warningf log and adds a trace event if channelz is on.
func Warningf(id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtWarning,
		})
	} else {
		logger.WarningDepth(1, msg)
	}
}

// Error logs and adds a trace event if channelz is on.
func Error(id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtError,
		})
	} else {
		logger.ErrorDepth(1, args...)
	}
}

// Errorf logs and adds a trace event if channelz is on.
func Errorf(id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtError,
		})
	} else {
		logger.ErrorDepth(1, msg)
	}
}
