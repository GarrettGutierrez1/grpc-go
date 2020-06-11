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

var logger = grpclog.Component("channelz")

func addTraceEventLogger(l grpclog.DepthLoggerV2, id int64, depth int, desc *TraceEventDesc) {
	for d := desc; d != nil; d = d.Parent {
		switch d.Severity {
		case CtUNKNOWN:
			l.InfoDepth(depth+1, d.Desc)
		case CtINFO:
			l.InfoDepth(depth+1, d.Desc)
		case CtWarning:
			l.WarningDepth(depth+1, d.Desc)
		case CtError:
			l.ErrorDepth(depth+1, d.Desc)
		}
	}
	if getMaxTraceEntry() == 0 {
		return
	}
	db.get().traceEvent(id, desc)
}

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

// InfoToLogger is equivalent to Info except it logs to a specified component logger.
func InfoToLogger(l grpclog.DepthLoggerV2, id int64, args ...interface{}) {
	if IsOn() {
		addTraceEventLogger(l, id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtINFO,
		})
	} else {
		l.InfoDepth(1, args...)
	}
}

// InfofToLogger is equivalent to Infof except it logs to a specified component logger.
func InfofToLogger(l grpclog.DepthLoggerV2, id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		addTraceEventLogger(l, id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtINFO,
		})
	} else {
		l.InfoDepth(1, msg)
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

// WarningToLogger is equivalent to Warning except it logs to a specified component logger.
func WarningToLogger(l grpclog.DepthLoggerV2, id int64, args ...interface{}) {
	if IsOn() {
		addTraceEventLogger(l, id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtWarning,
		})
	} else {
		l.WarningDepth(1, args...)
	}
}

// WarningfToLogger is equivalent to Warningf except it logs to a specified component logger.
func WarningfToLogger(l grpclog.DepthLoggerV2, id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		addTraceEventLogger(l, id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtWarning,
		})
	} else {
		l.WarningDepth(1, msg)
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

// ErrorToLogger is equivalent to Error except it logs to a specified component logger.
func ErrorToLogger(l grpclog.DepthLoggerV2, id int64, args ...interface{}) {
	if IsOn() {
		addTraceEventLogger(l, id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtError,
		})
	} else {
		l.ErrorDepth(1, args...)
	}
}

// ErrorfToLogger is equivalent to Errorf except it logs to a specified component logger.
func ErrorfToLogger(l grpclog.DepthLoggerV2, id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		addTraceEventLogger(l, id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtError,
		})
	} else {
		l.ErrorDepth(1, msg)
	}
}
