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
	"reflect"
	"testing"
)

func toString(m *map[string]*ComponentData) string {
	result := ""
	for k, v := range *m {
		result += fmt.Sprintf("\t%v: %v\n", k, *v)
	}
	return result
}

func parseAndCompare(t *testing.T, envVar string, envVars, preVars map[string]*ComponentData) {
	envVarsResult, preVarsResult := parseEnvironmentVar(envVar)
	if !reflect.DeepEqual(envVars, envVarsResult) {
		t.Errorf("Failed to parse environment variable '%v'.\nExpected standard variables:\n%vParsed:\n%v", envVar, toString(&envVars), toString(&envVarsResult))
	}
	if !reflect.DeepEqual(preVars, preVarsResult) {
		t.Errorf("Failed to parse environment variable '%v'.\nExpected standard variables:\n%vParsed:\n%v", envVar, toString(&preVars), toString(&preVarsResult))
	}
}

var parserTests = []struct {
	name    string
	in      string
	envVars map[string]*ComponentData
	preVars map[string]*ComponentData
}{
	{"Level", "INFO:INFO,WARNING:WARNING,ERROR:ERROR", map[string]*ComponentData{
		"INFO":    {"INFO", sentinel, levelInfo},
		"WARNING": {"WARNING", sentinel, levelWarning},
		"ERROR":   {"ERROR", sentinel, levelError},
	}, map[string]*ComponentData{}},
	{"Verbosity", "INFO:INFO,-1:INFO_-1,0:INFO_0,1:INFO_1", map[string]*ComponentData{
		"INFO": {"INFO", sentinel, levelInfo},
		"-1":   {"-1", -1, levelInfo},
		"0":    {"0", 0, levelInfo},
		"1":    {"1", 1, levelInfo},
	}, map[string]*ComponentData{}},
	{"PrefixLevel", "PRE_INFO*:INFO,PRE_WARNING*:WARNING,PRE_ERROR*:ERROR", map[string]*ComponentData{}, map[string]*ComponentData{
		"PRE_INFO":    {"PRE_INFO*", sentinel, levelInfo},
		"PRE_WARNING": {"PRE_WARNING*", sentinel, levelWarning},
		"PRE_ERROR":   {"PRE_ERROR*", sentinel, levelError},
	}},
	{"PrefixVerbosity", "PRE_INFO*:INFO,PRE_-1*:INFO_-1,PRE_0*:INFO_0,PRE_1*:INFO_1", map[string]*ComponentData{}, map[string]*ComponentData{
		"PRE_INFO": {"PRE_INFO*", sentinel, levelInfo},
		"PRE_-1":   {"PRE_-1*", -1, levelInfo},
		"PRE_0":    {"PRE_0*", 0, levelInfo},
		"PRE_1":    {"PRE_1*", 1, levelInfo},
	}},
}

func TestEnvironmentParser(t *testing.T) {
	for _, tt := range parserTests {
		t.Run(tt.name, func(t *testing.T) {
			parseAndCompare(t, tt.in, tt.envVars, tt.preVars)
		})
	}
}
