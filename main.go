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
	"global":                  "35.186.221.153",
	"asia-east1":              "104.155.201.52",
	"asia-east2":              "35.220.162.209",
	"asia-northeast1":         "104.198.86.148",
	"asia-northeast2":         "34.97.196.51",
	"asia-south1":             "35.200.186.152",
	"asia-southeast1":         "35.185.179.198",
	"australia-southeast1":    "35.189.6.113",
	"europe-north1":           "35.228.170.201",
	"europe-west1":            "104.199.82.109",
	"europe-west2":            "35.189.67.146",
	"europe-west3":            "35.198.78.172",
	"europe-west4":            "35.204.93.82",
	"europe-west6":            "34.65.3.254",
	"northamerica-northeast1": "35.203.57.164",
	"southamerica-east1":      "35.198.10.68",
	"us-central1":             "104.197.165.8",
	"us-east1":                "104.196.161.21",
	"us-east4":                "35.186.168.152",
	"us-west1":                "104.199.116.74",
	"us-west2":                "35.236.45.25",
}

var (
	top         bool
	number      int // number of requests for each region
	concurrency int
	timeout     time.Duration
	csv         bool
	verbose     bool
	region      string
	// TODO(jbd): Add payload options such as body size.

	client *http.Client // TODO(jbd): One client per worker?
)

func main() {
	flag.IntVar(&concurrency, "c", 10, "")
	flag.BoolVar(&csv, "csv", false, "")
	flag.IntVar(&number, "n", 10, "")
	flag.StringVar(&region, "r", "", "")
	flag.DurationVar(&timeout, "t", time.Duration(0), "")
	flag.BoolVar(&top, "top", false, "")
	flag.BoolVar(&verbose, "v", false, "")

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
-n   Number of requests to be made to each region.
     By default 10; can't be negative.
-c   Max number of requests to be made at any time.
     By default 10; can't be negative or zero.
-r   Report latency for an individual region.
-t   Timeout. By default, no timeout.
     Examples: "500ms", "1s", "1s500ms".
-top If true, only the top (non-global) region is printed.

-csv CSV output; disables verbose output.
-v   Verbose output.

Need a website version? See gcping.com
`
