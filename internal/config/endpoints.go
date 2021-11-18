// Copyright 2021 Google LLC
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

package config

// Endpoint represents a Cloud Run service deploy in a particular region.
type Endpoint struct {
	// URL is the HTTPS URL of the service
	URL string
	// Region is the programmatic name of the region where the endpoint is
	// deloyed, e.g., us-central1.
	Region string
	// RegionName is the geographic name of the region, e.g., Iowa.
	RegionName string
}

// AllEndpoints associates a region name with its Cloud Run Endpoint.
var AllEndpoints = map[string]Endpoint{
	"global": {
		URL:        "https://global.gcping.com",
		Region:     "global",
		RegionName: "Global HTTP Load Balancer",
	},
	"asia-east1": {
		URL:        "https://asia-east1-5tkroniexa-de.a.run.app",
		Region:     "asia-east1",
		RegionName: "Taiwan",
	},
	"asia-east2": {
		URL:        "https://asia-east2-5tkroniexa-df.a.run.app",
		Region:     "asia-east2",
		RegionName: "Hong Kong",
	},
	"asia-northeast1": {
		URL:        "https://asia-northeast1-5tkroniexa-an.a.run.app",
		Region:     "asia-northeast1",
		RegionName: "Tokyo",
	},
	"asia-northeast2": {
		URL:        "https://asia-northeast2-5tkroniexa-dt.a.run.app",
		Region:     "asia-northeast2",
		RegionName: "Osaka",
	},
	"asia-northeast3": {
		URL:        "https://asia-northeast3-5tkroniexa-du.a.run.app",
		Region:     "asia-northeast3",
		RegionName: "Seoul",
	},
	"asia-south1": {
		URL:        "https://asia-south1-5tkroniexa-el.a.run.app",
		Region:     "asia-south1",
		RegionName: "Mumbai",
	},
	"asia-south2": {
		URL:        "https://asia-south2-5tkroniexa-em.a.run.app",
		Region:     "asia-south2",
		RegionName: "Delhi",
	},
	"asia-southeast1": {
		URL:        "https://asia-southeast1-5tkroniexa-as.a.run.app",
		Region:     "asia-southeast1",
		RegionName: "Singapore",
	},
	"asia-southeast2": {
		URL:        "https://asia-southeast2-5tkroniexa-et.a.run.app",
		Region:     "asia-southeast2",
		RegionName: "Jakarta",
	},
	"australia-southeast1": {
		URL:        "https://australia-southeast1-5tkroniexa-ts.a.run.app",
		Region:     "australia-southeast1",
		RegionName: "Sydney",
	},
	"australia-southeast2": {
		URL:        "https://australia-southeast2-5tkroniexa-km.a.run.app",
		Region:     "australia-southeast2",
		RegionName: "Melbourne",
	},
	"europe-central2": {
		URL:        "https://europe-central2-5tkroniexa-lm.a.run.app",
		Region:     "europe-central2",
		RegionName: "Warsaw",
	},
	"europe-north1": {
		URL:        "https://europe-north1-5tkroniexa-lz.a.run.app",
		Region:     "europe-north1",
		RegionName: "Finland",
	},
	"europe-west1": {
		URL:        "https://europe-west1-5tkroniexa-ew.a.run.app",
		Region:     "europe-west1",
		RegionName: "Belgium",
	},
	"europe-west2": {
		URL:        "https://europe-west2-5tkroniexa-nw.a.run.app",
		Region:     "europe-west2",
		RegionName: "London",
	},
	"europe-west3": {
		URL:        "https://europe-west3-5tkroniexa-ey.a.run.app",
		Region:     "europe-west3",
		RegionName: "Frankfurt",
	},
	"europe-west4": {
		URL:        "https://europe-west4-5tkroniexa-ez.a.run.app",
		Region:     "europe-west4",
		RegionName: "Netherlands",
	},
	"europe-west6": {
		URL:        "https://europe-west6-5tkroniexa-oa.a.run.app",
		Region:     "europe-west6",
		RegionName: "Zurich",
	},
	"northamerica-northeast1": {
		URL:        "https://northamerica-northeast1-5tkroniexa-nn.a.run.app",
		Region:     "northamerica-northeast1",
		RegionName: "Montréal",
	},
	"northamerica-northeast2": {
		URL:        "https://northamerica-northeast2-5tkroniexa-pd.a.run.app",
		Region:     "northamerica-northeast2",
		RegionName: "Toronto",
	},
	"southamerica-east1": {
		URL:        "https://southamerica-east1-5tkroniexa-rj.a.run.app",
		Region:     "southamerica-east1",
		RegionName: "São Paulo",
	},
	"southamerica-west1": {
		URL:        "https://southamerica-west1-5tkroniexa-tl.a.run.app",
		Region:     "southamerica-west1",
		RegionName: "Santiago",
	},
	"us-central1": {
		URL:        "https://us-central1-5tkroniexa-uc.a.run.app",
		Region:     "us-central1",
		RegionName: "Iowa",
	},
	"us-east1": {
		URL:        "https://us-east1-5tkroniexa-ue.a.run.app",
		Region:     "us-east1",
		RegionName: "South Carolina",
	},
	"us-east4": {
		URL:        "https://us-east4-5tkroniexa-uk.a.run.app",
		Region:     "us-east4",
		RegionName: "North Virginia",
	},
	"us-west1": {
		URL:        "https://us-west1-5tkroniexa-uw.a.run.app",
		Region:     "us-west1",
		RegionName: "Oregon",
	},
	"us-west2": {
		URL:        "https://us-west2-5tkroniexa-wl.a.run.app",
		Region:     "us-west2",
		RegionName: "Los Angeles",
	},
	"us-west3": {
		URL:        "https://us-west3-5tkroniexa-wm.a.run.app",
		Region:     "us-west3",
		RegionName: "Salt Lake City",
	},
	"us-west4": {
		URL:        "https://us-west4-5tkroniexa-wn.a.run.app",
		Region:     "us-west4",
		RegionName: "Las Vegas",
	},
}
