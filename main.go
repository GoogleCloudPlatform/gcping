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

// Program gcping pings GCP regions and reports about the latency.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TODO(jbd): Add more zones.
var endpoints = map[string]string{
	"global":                  "https://global.gcping.com",
	"asia-east1":              "https://asia-east1-bmlfzs4h6a-de.a.run.app",
	"asia-east2":              "https://asia-east2-bmlfzs4h6a-df.a.run.app",
	"asia-northeast1":         "https://asia-northeast1-bmlfzs4h6a-an.a.run.app",
	"asia-northeast2":         "https://asia-northeast2-bmlfzs4h6a-dt.a.run.app",
	"asia-northeast3":         "https://asia-northeast3-bmlfzs4h6a-du.a.run.app",
	"asia-south1":             "https://asia-south1-bmlfzs4h6a-el.a.run.app",
	"asia-southeast1":         "https://asia-southeast1-bmlfzs4h6a-as.a.run.app",
	"asia-southeast2":         "https://asia-southeast2-bmlfzs4h6a-et.a.run.app",
	"australia-southeast1":    "https://australia-southeast1-bmlfzs4h6a-ts.a.run.app",
	"europe-north1":           "https://europe-north1-bmlfzs4h6a-lz.a.run.app",
	"europe-west1":            "https://europe-west1-bmlfzs4h6a-ew.a.run.app",
	"europe-west2":            "https://europe-west2-bmlfzs4h6a-nw.a.run.app",
	"europe-west3":            "https://europe-west3-bmlfzs4h6a-ey.a.run.app",
	"europe-west4":            "https://europe-west4-bmlfzs4h6a-ez.a.run.app",
	"europe-west6":            "https://europe-west6-bmlfzs4h6a-oa.a.run.app",
	"northamerica-northeast1": "https://northamerica-northeast1-bmlfzs4h6a-nn.a.run.app",
	"southamerica-east1":      "https://southamerica-east1-bmlfzs4h6a-rj.a.run.app",
	"us-central1":             "https://us-central1-bmlfzs4h6a-uc.a.run.app",
	"us-east1":                "https://us-east1-bmlfzs4h6a-ue.a.run.app",
	"us-east4":                "https://us-east4-bmlfzs4h6a-uk.a.run.app",
	"us-west1":                "https://us-west1-bmlfzs4h6a-uw.a.run.app",
}

var (
	top         bool
	number      int // number of requests for each region
	concurrency int
	timeout     time.Duration
	csv         bool
	csvCum      bool
	verbose     bool
	region      string
	// TODO(jbd): Add payload options such as body size.

	client *http.Client // TODO(jbd): One client per worker?
)

func main() {
	flag.BoolVar(&top, "top", false, "")
	flag.IntVar(&number, "n", 10, "")
	flag.IntVar(&concurrency, "c", 10, "")
	flag.DurationVar(&timeout, "t", time.Duration(0), "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&csv, "csv", false, "")
	flag.BoolVar(&csvCum, "csv-cum", false, "")
	flag.StringVar(&region, "r", "", "")

	flag.Usage = usage
	flag.Parse()

	if number < 0 || concurrency <= 0 {
		usage()
	}
	if csv {
		verbose = false // if output is CSV, no need for verbose output
	}

	if region != "" {
		if _, found := endpoints[region]; !found {
			fmt.Printf("region %q is not supported or does not exist\n", region)
			os.Exit(1)
		}
	}

	client = &http.Client{
		Timeout: timeout,
	}

	w := &worker{}
	go w.start()

	switch {
	case region != "":
		w.reportRegion(region)
	case top:
		w.reportTop()
	case csvCum:
		w.reportCSV()
	default:
		w.reportAll()
	}
}

func usage() {
	fmt.Println(usageText)
	os.Exit(0)
}

var usageText = `gcping [options...]

Options:
-n       Number of requests to be made to each region.
         By default 10; can't be negative.
-c       Max number of requests to be made at any time.
         By default 10; can't be negative or zero.
-r       Report latency for an individual region.
-t       Timeout. By default, no timeout.
         Examples: "500ms", "1s", "1s500ms".
-top     If true, only the top (non-global) region is printed.
-csv-cum If true, cumulative value is printed in CSV; disables default report.

-csv     CSV output; disables verbose output.
-v       Verbose output.

Need a website version? See gcping.com
`
