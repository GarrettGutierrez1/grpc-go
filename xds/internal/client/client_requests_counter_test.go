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

package client_test

import (
	"sync"
	"testing"

	"google.golang.org/grpc/xds/internal/client"
)

type counterTest struct {
	name            string
	circuitBreaking bool
	maxRequests     uint32
	numRequests     uint32
	errorExpected   bool
}

func testCounter(t *testing.T, test counterTest) {
	counter := client.ServiceRequestsCounter{ServiceName: test.name}
	counter.UpdateCounter(&test.maxRequests)
	if !test.circuitBreaking {
		counter.UpdateCounter(nil)
	}
	requestsStartedWg := sync.WaitGroup{}
	requestsStartedWg.Add(1)
	requestsSent := sync.WaitGroup{}
	requestsSent.Add(int(test.numRequests))
	requestsDoneWg := sync.WaitGroup{}
	requestsDoneWg.Add(int(test.numRequests))
	var firstError error = nil
	errorMu := sync.Mutex{}
	fail := func(err error) {
		errorMu.Lock()
		defer errorMu.Unlock()
		if firstError == nil {
			firstError = err
		}
	}
	for i := 0; i < int(test.numRequests); i++ {
		go func() {
			defer requestsDoneWg.Done()
			if err := counter.StartRequest(); err != nil {
				fail(err)
				requestsSent.Done()
				return
			}
			requestsSent.Done()
			requestsStartedWg.Wait()
			if err := counter.EndRequest(); err != nil {
				fail(err)
			}
		}()
	}
	requestsSent.Wait()
	requestsStartedWg.Done()
	requestsDoneWg.Wait()
	if test.errorExpected && firstError == nil {
		t.Error("no error when error expected")
	}
	if !test.errorExpected && firstError != nil {
		t.Errorf("error starting request: %v", firstError)
	}
}

func (s) TestRequestsCounter(t *testing.T) {
	tests := []counterTest{
		{
			name:            "cb-on-no-exceed",
			circuitBreaking: true,
			maxRequests:     1024,
			numRequests:     1024,
			errorExpected:   false,
		},
		{
			name:            "cb-off-exceeds",
			circuitBreaking: false,
			maxRequests:     32,
			numRequests:     64,
			errorExpected:   false,
		},
		{
			name:            "cb-on-exceeds",
			circuitBreaking: true,
			maxRequests:     32,
			numRequests:     64,
			errorExpected:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCounter(t, test)
		})
	}
}
