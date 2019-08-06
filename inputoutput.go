// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"time"
)

type input struct {
	region   string
	endpoint string
}

func (i *input) HTTP() {
	i.benchmark(func() error {
		req, _ := http.NewRequest("GET", "http://"+i.endpoint+"/ping", nil)
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("status code: %v", res.StatusCode)
		}
		return nil
	})
}

func (i *input) benchmark(fn func() error) {
	if verbose {
		fmt.Printf("Pinging %q\n", i.region)
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	o := output{
		region:   i.region,
		duration: duration,
	}
	if err != nil {
		o.errors++
	}
	outputs <- o

	if verbose {
		fmt.Printf("Ping to %q completed in %v\n", i.region, duration)
	}

	if csv {
		fmt.Printf("%v,%v,%v,%v\n", i.region, i.endpoint, duration.Nanoseconds(), err != nil)
	}
}

type output struct {
	region   string
	duration time.Duration
	errors   int
}

type outputSorter struct {
	averages []output
}

// Len is part of sort.Interface.
func (s *outputSorter) Len() int {
	return len(s.averages)
}

// Swap is part of sort.Interface.
func (s *outputSorter) Swap(i, j int) {
	s.averages[i], s.averages[j] = s.averages[j], s.averages[i]
}

// Less is part of sort.Interface.
func (s *outputSorter) Less(i, j int) bool {
	return s.averages[i].duration < s.averages[j].duration
}
