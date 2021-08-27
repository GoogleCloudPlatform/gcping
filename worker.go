// Copyright 2010 Google LLC
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
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/GoogleCloudPlatform/gcping/internal/config"
)

type input struct {
	region   string
	endpoint string
}

func (i *input) HTTP() output {
	return i.benchmark(func() error {
		req, _ := http.NewRequest("GET", i.endpoint+"/ping", nil)
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

func (i *input) benchmark(fn func() error) output {
	if verbose {
		fmt.Printf("Pinging %q\n", i.region)
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	o := output{
		region:    i.region,
		durations: []time.Duration{duration},
	}
	if err != nil {
		o.errors++
	}

	if verbose {
		fmt.Printf("Ping to %q completed in %v\n", i.region, duration)
	}

	if csv {
		fmt.Printf("%v,%v,%v,%v\n", i.region, i.endpoint, duration.Nanoseconds(), err != nil)
	}

	return o
}

type output struct {
	region    string
	durations []time.Duration
	errors    int

	med time.Duration // median of durations; calculated on first call to median()
}

func (o *output) median() time.Duration {
	if o.med == 0 {
		// Sort durations and pick the middle one.
		sort.Slice(o.durations, func(i, j int) bool {
			return o.durations[i] < o.durations[j]
		})
		o.med = o.durations[len(o.durations)/2]
	}
	return o.med

}

type worker struct {
	inputs  chan input
	outputs chan output
}

func (w *worker) start() {
	for worker := 0; worker < concurrency; worker++ {
		go func() {
			for m := range w.inputs {
				o := m.HTTP()
				w.outputs <- o
			}
		}()
	}
}

func (w *worker) sortOutput() []output {
	m := make(map[string]output)
	for i := 0; i < w.size(region); i++ {
		o := <-w.outputs

		a := m[o.region]

		a.region = o.region
		a.durations = append(a.durations, o.durations[0])
		a.errors += o.errors

		m[o.region] = a
	}
	all := make([]output, 0, len(m))
	for _, t := range m {
		all = append(all, t)
	}

	// sort all by median duration.
	sort.Slice(all, func(i, j int) bool {
		return all[i].median() < all[j].median()
	})
	return all
}

func (w *worker) reportAll() {
	w.inputs = make(chan input, concurrency)
	w.outputs = make(chan output, w.size(region))
	for i := 0; i < number; i++ {
		for r, e := range config.AllEndpoints {
			w.inputs <- input{region: r, endpoint: e.URL}
		}
	}
	close(w.inputs)

	sorted := w.sortOutput()
	tr := tabwriter.NewWriter(os.Stdout, 3, 2, 2, ' ', 0)
	for i, a := range sorted {
		fmt.Fprintf(tr, "%2d.\t[%v]\t%v", i+1, a.region, a.median())
		if a.errors > 0 {
			fmt.Fprintf(tr, "\t(%d errors)", a.errors)
		}
		fmt.Fprintln(tr)
	}
	tr.Flush()
}

func (w *worker) reportCSV() {
	w.inputs = make(chan input, concurrency)
	w.outputs = make(chan output, w.size(region))
	for i := 0; i < number; i++ {
		for r, e := range config.AllEndpoints {
			w.inputs <- input{region: r, endpoint: e.URL}
		}
	}
	close(w.inputs)

	sorted := w.sortOutput()
	fmt.Println("region,latency_ns,errors")
	for _, a := range sorted {
		fmt.Printf("%v,%v,%v\n", a.region, a.median().Nanoseconds(), a.errors)
	}
}

func (w *worker) reportTop() {
	w.inputs = make(chan input, concurrency)
	w.outputs = make(chan output, w.size(region))
	for i := 0; i < number; i++ {
		for r, e := range config.AllEndpoints {
			w.inputs <- input{region: r, endpoint: e.URL}
		}
	}
	close(w.inputs)

	sorted := w.sortOutput()
	t := sorted[0].region
	if t == "global" {
		t = sorted[1].region
	}
	fmt.Print(t)
	return
}

func (w *worker) reportRegion(region string) {
	w.inputs = make(chan input, concurrency)
	w.outputs = make(chan output, w.size(region))
	for i := 0; i < number; i++ {
		e, _ := config.AllEndpoints[region]
		w.inputs <- input{region: region, endpoint: e.URL}
	}
	close(w.inputs)

	sorted := w.sortOutput()
	fmt.Print(sorted[0].median())

}

func (w *worker) size(region string) int {
	if region != "" {
		return number
	}
	return number * len(config.AllEndpoints)
}
