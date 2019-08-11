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
	"sort"
	"text/tabwriter"
	"time"
)

var endpoints = map[string]string{}

// TODO(jbd): Add more zones.
var gcpendpoints = map[string]string{
	"global":                  "http://35.186.221.153/ping",
	"asia-east1":              "http://104.155.201.52/ping",
	"asia-east2":              "http://35.220.162.209/ping",
	"asia-northeast1":         "http://104.198.86.148/ping",
	"asia-northeast2":         "http://34.97.196.51/ping",
	"asia-south1":             "http://35.200.186.152/ping",
	"asia-southeast1":         "http://35.185.179.198/ping",
	"australia-southeast1":    "http://35.189.6.113/ping",
	"europe-north1":           "http://35.228.170.201/ping",
	"europe-west1":            "http://104.199.82.109/ping",
	"europe-west2":            "http://35.189.67.146/ping",
	"europe-west3":            "http://35.198.78.172/ping",
	"europe-west4":            "http://35.204.93.82/ping",
	"europe-west6":            "http://34.65.3.254/ping",
	"northamerica-northeast1": "http://35.203.57.164/ping",
	"southamerica-east1":      "http://35.198.10.68/ping",
	"us-central1":             "http://104.197.165.8/ping",
	"us-east1":                "http://104.196.161.21/ping",
	"us-east4":                "http://35.186.168.152/ping",
	"us-west1":                "http://104.199.116.74/ping",
	"us-west2":                "http://35.236.45.25/ping",
}

var mutlicloudendpoints = map[string]string{
	"global":                    "http://35.186.221.153/ping",
	"asia-east1":                "http://104.155.201.52/ping",
	"asia-east2":                "http://35.220.162.209/ping",
	"asia-northeast1":           "http://104.198.86.148/ping",
	"asia-northeast2":           "http://34.97.196.51/ping",
	"asia-south1":               "http://35.200.186.152/ping",
	"asia-southeast1":           "http://35.185.179.198/ping",
	"australia-southeast1":      "http://35.189.6.113/ping",
	"europe-north1":             "http://35.228.170.201/ping",
	"europe-west1":              "http://104.199.82.109/ping",
	"europe-west2":              "http://35.189.67.146/ping",
	"europe-west3":              "http://35.198.78.172/ping",
	"europe-west4":              "http://35.204.93.82/ping",
	"europe-west6":              "http://34.65.3.254/ping",
	"northamerica-northeast1":   "http://35.203.57.164/ping",
	"southamerica-east1":        "http://35.198.10.68/ping",
	"us-central1":               "http://104.197.165.8/ping",
	"us-east1":                  "http://104.196.161.21/ping",
	"us-east4":                  "http://35.186.168.152/ping",
	"us-west1":                  "http://104.199.116.74/ping",
	"us-west2":                  "http://35.236.45.25/ping",
	"aws-us-east1":              "http://dynamodb.us-east-1.amazonaws.com/",
	"aws-us-east2":              "http://dynamodb.us-east-2.amazonaws.com/",
	"aws-us-west1":              "http://dynamodb.us-west-1.amazonaws.com/",
	"aws-us-west2":              "http://dynamodb.us-east-2.amazonaws.com/",
	"aws-ap-east-1":             "http://dynamodb.ap-east-1.amazonaws.com/",
	"aws-ap-south-1":            "http://dynamodb.ap-south-1.amazonaws.com/",
	"aws-ap-northeast-3":        "http://dynamodb.ap-northeast-3.amazonaws.com/",
	"aws-ap-northeast-2":        "http://dynamodb.ap-northeast-2.amazonaws.com/",
	"aws-ap-southeast-1":        "http://dynamodb.ap-southeast-1.amazonaws.com/",
	"aws-ap-southeast-2":        "http://dynamodb.ap-southeast-2.amazonaws.com/",
	"aws-ap-northeast-1":        "http://dynamodb.ap-northeast-1.amazonaws.com/",
	"aws-ca-central-1":          "http://dynamodb.ca-central-1.amazonaws.com/",
	"aws-cn-north-1":            "http://dynamodb.cn-north-1.amazonaws.com.cn/",
	"aws-cn-northwest-1":        "http://dynamodb.cn-northwest-1.amazonaws.com.cn/",
	"aws-eu-central-1":          "http://dynamodb.eu-central-1.amazonaws.com/",
	"aws-eu-west-1":             "http://dynamodb.eu-west-1.amazonaws.com/",
	"aws-eu-west-2":             "http://dynamodb.eu-west-2.amazonaws.com/",
	"aws-eu-west-3":             "http://dynamodb.eu-west-3.amazonaws.com/",
	"aws-eu-north-1":            "http://dynamodb.eu-north-1.amazonaws.com/",
	"aws-me-south-1":            "http://dynamodb.me-south-1.amazonaws.com/",
	"aws-sa-east-1":             "http://dynamodb.sa-east-1.amazonaws.com/",
	"aws-us-gov-east-1":         "http://dynamodb.us-gov-east-1.amazonaws.com/",
	"aws-us-gov-west-1":         "http://dynamodb.us-gov-west-1.amazonaws.com/",
	"azure-us-central":          "http://centralus.blob.core.windows.net/",
	"azure-southafrica-north":   "http://southafricanorth.blob.core.windows.net/",
	"azure-uae-north":           "http://unitedarabemirates.blob.core.windows.net/",
	"azurefrance-central":       "http://francecentral.blob.core.windows.net/",
	"azure-uk-west":             "http://ukwest.blob.core.windows.net/",
	"azure-uk-south":            "http://uksouth.blob.core.windows.net/",
	"azure-europe-west":         "http://westeurope.blob.core.windows.net/",
	"azure-europe-north":        "http://northeurope.blob.core.windows.net/",
	"azure-australia-central2":  "http://australiacentral2.blob.core.windows.net/",
	"azure-korea":               "http://korea.blob.core.windows.net/",
	"azure-india-west":          "http://westindia.blob.core.windows.net/",
	"azure-india-central":       "http://centralindia.blob.core.windows.net/",
	"azure-australia-southeast": "http://australiasoutheast.blob.core.windows.net/",
	"azure-australia-east":      "http://australiaeast.blob.core.windows.net/",
	"azure-japan-east":          "http://japaneast.blob.core.windows.net/",
	"azure-japan-west":          "http://japanwest.blob.core.windows.net/",
	"azure-asia-southeast":      "http://southeastasia.blob.core.windows.net/",
	"azure-asia-east":           "http://eastasia.blob.core.windows.net/",
	"azure-us-west2":            "http://uswest2.blob.core.windows.net/",
	"azure-us-westcentral":      "http://westcentralus.blob.core.windows.net/",
	"azure-canada-east":         "http://canadaeast.blob.core.windows.net/",
	"azure-canada-central":      "http://canadacentral.blob.core.windows.net/",
	"azure-brazil-south":        "http://brazil.blob.core.windows.net/",
	"azure-us-southcentral":     "http://southcentralus.blob.core.windows.net/",
	"azure-us-nourthcentral":    "http://northcentralus.blob.core.windows.net/",
	"azure-us-west":             "http://westus.blob.core.windows.net/",
	"azure-us-east2":            "http://eastus2.blob.core.windows.net/",
	"azure-us-east":             "http://eastus.blob.core.windows.net/",
}

var azureendpoints = map[string]string{
	"azure-us-central":          "http://centralus.blob.core.windows.net/",
	"azure-southafrica-north":   "http://southafricanorth.blob.core.windows.net/",
	"azure-uae-north":           "http://unitedarabemirates.blob.core.windows.net/",
	"azurefrance-central":       "http://francecentral.blob.core.windows.net/",
	"azure-uk-west":             "http://ukwest.blob.core.windows.net/",
	"azure-uk-south":            "http://uksouth.blob.core.windows.net/",
	"azure-europe-west":         "http://westeurope.blob.core.windows.net/",
	"azure-europe-north":        "http://northeurope.blob.core.windows.net/",
	"azure-australia-central2":  "http://australiacentral2.blob.core.windows.net/",
	"azure-korea":               "http://korea.blob.core.windows.net/",
	"azure-india-west":          "http://westindia.blob.core.windows.net/",
	"azure-india-central":       "http://centralindia.blob.core.windows.net/",
	"azure-australia-southeast": "http://australiasoutheast.blob.core.windows.net/",
	"azure-australia-east":      "http://australiaeast.blob.core.windows.net/",
	"azure-japan-east":          "http://japaneast.blob.core.windows.net/",
	"azure-japan-west":          "http://japanwest.blob.core.windows.net/",
	"azure-asia-southeast":      "http://southeastasia.blob.core.windows.net/",
	"azure-asia-east":           "http://eastasia.blob.core.windows.net/",
	"azure-us-west2":            "http://uswest2.blob.core.windows.net/",
	"azure-us-westcentral":      "http://westcentralus.blob.core.windows.net/",
	"azure-canada-east":         "http://canadaeast.blob.core.windows.net/",
	"azure-canada-central":      "http://canadacentral.blob.core.windows.net/",
	"azure-brazil-south":        "http://brazil.blob.core.windows.net/",
	"azure-us-southcentral":     "http://southcentralus.blob.core.windows.net/",
	"azure-us-nourthcentral":    "http://northcentralus.blob.core.windows.net/",
	"azure-us-west":             "http://westus.blob.core.windows.net/",
	"azure-us-east2":            "http://eastus2.blob.core.windows.net/",
	"azure-us-east":             "http://eastus.blob.core.windows.net/",
}
var awsendpoints = map[string]string{
	"aws-us-east1":       "http://dynamodb.us-east-1.amazonaws.com/",
	"aws-us-east2":       "http://dynamodb.us-east-2.amazonaws.com/",
	"aws-us-west1":       "http://dynamodb.us-west-1.amazonaws.com/",
	"aws-us-west2":       "http://dynamodb.us-east-2.amazonaws.com/",
	"aws-ap-east-1":      "http://dynamodb.ap-east-1.amazonaws.com/",
	"aws-ap-south-1":     "http://dynamodb.ap-south-1.amazonaws.com/",
	"aws-ap-northeast-3": "http://dynamodb.ap-northeast-3.amazonaws.com/",
	"aws-ap-northeast-2": "http://dynamodb.ap-northeast-2.amazonaws.com/",
	"aws-ap-southeast-1": "http://dynamodb.ap-southeast-1.amazonaws.com/",
	"aws-ap-southeast-2": "http://dynamodb.ap-southeast-2.amazonaws.com/",
	"aws-ap-northeast-1": "http://dynamodb.ap-northeast-1.amazonaws.com/",
	"aws-ca-central-1":   "http://dynamodb.ca-central-1.amazonaws.com/",
	"aws-cn-north-1":     "http://dynamodb.cn-north-1.amazonaws.com.cn/",
	"aws-cn-northwest-1": "http://dynamodb.cn-northwest-1.amazonaws.com.cn/",
	"aws-eu-central-1":   "http://dynamodb.eu-central-1.amazonaws.com/",
	"aws-eu-west-1":      "http://dynamodb.eu-west-1.amazonaws.com/",
	"aws-eu-west-2":      "http://dynamodb.eu-west-2.amazonaws.com/",
	"aws-eu-west-3":      "http://dynamodb.eu-west-3.amazonaws.com/",
	"aws-eu-north-1":     "http://dynamodb.eu-north-1.amazonaws.com/",
	"aws-me-south-1":     "http://dynamodb.me-south-1.amazonaws.com/",
	"aws-sa-east-1":      "http://dynamodb.sa-east-1.amazonaws.com/",
	"aws-us-gov-east-1":  "http://dynamodb.us-gov-east-1.amazonaws.com/",
	"aws-us-gov-west-1":  "http://dynamodb.us-gov-west-1.amazonaws.com/",
}
var (
	top         bool
	number      int // number of requests for each region
	concurrency int
	timeout     time.Duration
	csv         bool
	verbose     bool
	multicloud  bool
	azure       bool
	aws         bool
	// TODO(jbd): Add payload options such as body size.

	client  *http.Client // TODO(jbd): One client per worker?
	inputs  chan input
	outputs chan output
)

func main() {
	flag.BoolVar(&top, "top", false, "")
	flag.IntVar(&number, "n", 10, "")
	flag.IntVar(&concurrency, "c", 10, "")
	flag.DurationVar(&timeout, "t", time.Duration(0), "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&csv, "csv", false, "")
	//add flag for multicloud
	flag.BoolVar(&multicloud, "mc", false, "")
	//add flag for azure
	flag.BoolVar(&azure, "azure", false, "")
	//add flag for aws
	flag.BoolVar(&aws, "aws", false, "")

	flag.Usage = usage
	flag.Parse()

	if number < 0 || concurrency <= 0 {
		usage()
	}
	if csv {
		verbose = false // if output is CSV, no need for verbose output
	}
	if multicloud {
		endpoints = mutlicloudendpoints //use multicloud endpoints map
	} else if azure {
		endpoints = azureendpoints //use azure only endpoints
	} else if aws {
		endpoints = awsendpoints //use aws only endpoints
	} else {
		endpoints = gcpendpoints //original gcp endpoints
	}

	client = &http.Client{
		Timeout: timeout,
	}

	go start()
	inputs = make(chan input, concurrency)
	outputs = make(chan output, number*len(endpoints))
	for i := 0; i < number; i++ {
		for r, e := range endpoints {
			inputs <- input{region: r, endpoint: e}
		}
	}
	close(inputs)
	report()
}

func start() {
	for worker := 0; worker < concurrency; worker++ {
		go func() {
			for m := range inputs {
				m.HTTP()
			}
		}()
	}
}

func report() {
	m := make(map[string]output)
	for i := 0; i < number*len(endpoints); i++ {
		o := <-outputs

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

	if top {
		t := all[0].region
		if t == "global" {
			t = all[1].region
		}
		fmt.Print(t)
		return
	}

	tr := tabwriter.NewWriter(os.Stdout, 3, 2, 2, ' ', 0)
	for i, a := range all {
		fmt.Fprintf(tr, "%2d.\t[%v]\t%v", i+1, a.region, a.median())
		if a.errors > 0 {
			fmt.Fprintf(tr, "\t(%d errors)", a.errors)
		}
		fmt.Fprintln(tr)
	}
	tr.Flush()
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
-t   Timeout. By default, no timeout.

     Examples: "500ms", "1s", "1s500ms".

-top If true, only the top (non-global) region is printed
	 
-aws     Use only aws regions
-azure   Use only azure regions. Recommend to limit number of req to < 10
-mc      All your clouds are belong to us (multicloud)


-csv CSV output; disables verbose output.
-v   Verbose output.

Need a website version? See gcping.com
`
