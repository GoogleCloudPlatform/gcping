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
	// Lat is the latitude of the physical region
	Lat float32
	// Lng is the longitude of the physical region
	Lng float32
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
		Lat:        23.69781,
		Lng:        120.960515,
	},
	"asia-east2": {
		URL:        "https://asia-east2-5tkroniexa-df.a.run.app",
		Region:     "asia-east2",
		RegionName: "Hong Kong",
		Lat:        22.3193039,
		Lng:        114.1693611,
	},
	"asia-northeast1": {
		URL:        "https://asia-northeast1-5tkroniexa-an.a.run.app",
		Region:     "asia-northeast1",
		RegionName: "Tokyo",
		Lat:        35.6761919,
		Lng:        139.6503106,
	},
	"asia-northeast2": {
		URL:        "https://asia-northeast2-5tkroniexa-dt.a.run.app",
		Region:     "asia-northeast2",
		RegionName: "Osaka",
		Lat:        34.6937249,
		Lng:        135.5022535,
	},
	"asia-northeast3": {
		URL:        "https://asia-northeast3-5tkroniexa-du.a.run.app",
		Region:     "asia-northeast3",
		RegionName: "Seoul",
		Lat:        37.566535,
		Lng:        126.9779692,
	},
	"asia-south1": {
		URL:        "https://asia-south1-5tkroniexa-el.a.run.app",
		Region:     "asia-south1",
		RegionName: "Mumbai",
		Lat:        19.0759837,
		Lng:        72.8776559,
	},
	"asia-south2": {
		URL:        "https://asia-south2-5tkroniexa-em.a.run.app",
		Region:     "asia-south2",
		RegionName: "Delhi",
		Lat:        28.7040592,
		Lng:        77.10249019999999,
	},
	"asia-southeast1": {
		URL:        "https://asia-southeast1-5tkroniexa-as.a.run.app",
		Region:     "asia-southeast1",
		RegionName: "Singapore",
		Lat:        1.352083,
		Lng:        103.819836,
	},
	"asia-southeast2": {
		URL:        "https://asia-southeast2-5tkroniexa-et.a.run.app",
		Region:     "asia-southeast2",
		RegionName: "Jakarta",
		Lat:        -6.2087634,
		Lng:        106.845599,
	},
	"australia-southeast1": {
		URL:        "https://australia-southeast1-5tkroniexa-ts.a.run.app",
		Region:     "australia-southeast1",
		RegionName: "Sydney",
		Lat:        -33.8688197,
		Lng:        151.2092955,
	},
	"australia-southeast2": {
		URL:        "https://australia-southeast2-5tkroniexa-km.a.run.app",
		Region:     "australia-southeast2",
		RegionName: "Melbourne",
		Lat:        -37.8136276,
		Lng:        144.9630576,
	},
	"europe-central2": {
		URL:        "https://europe-central2-5tkroniexa-lm.a.run.app",
		Region:     "europe-central2",
		RegionName: "Warsaw",
		Lat:        52.2329172,
		Lng:        20.9911553
	},
	"europe-north1": {
		URL:        "https://europe-north1-5tkroniexa-lz.a.run.app",
		Region:     "europe-north1",
		RegionName: "Finland",
		Lat:        61.92410999999999,
		Lng:        25.7481511,
	},
	"europe-west1": {
		URL:        "https://europe-west1-5tkroniexa-ew.a.run.app",
		Region:     "europe-west1",
		RegionName: "Belgium",
		Lat:        50.503887,
		Lng:        4.469936,
	},
	"europe-west2": {
		URL:        "https://europe-west2-5tkroniexa-nw.a.run.app",
		Region:     "europe-west2",
		RegionName: "London",
		Lat:        51.5073509,
		Lng:        -0.1277583,
	},
	"europe-west3": {
		URL:        "https://europe-west3-5tkroniexa-ey.a.run.app",
		Region:     "europe-west3",
		RegionName: "Frankfurt",
		Lat:        50.1109221,
		Lng:        8.6821267,
	},
	"europe-west4": {
		URL:        "https://europe-west4-5tkroniexa-ez.a.run.app",
		Region:     "europe-west4",
		RegionName: "Netherlands",
		Lat:        52.132633,
		Lng:        5.291265999999999,
	},
	"europe-west6": {
		URL:        "https://europe-west6-5tkroniexa-oa.a.run.app",
		Region:     "europe-west6",
		RegionName: "Zurich",
		Lat:        47.3768866,
		Lng:        8.541694,
	},
	"northamerica-northeast1": {
		URL:        "https://northamerica-northeast1-5tkroniexa-nn.a.run.app",
		Region:     "northamerica-northeast1",
		RegionName: "Montréal",
		Lat:        45.5016889,
		Lng:        -73.567256,
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
		Lat:        -21.2922457,
		Lng:        -50.3428431,
	},
	"us-central1": {
		URL:        "https://us-central1-5tkroniexa-uc.a.run.app",
		Region:     "us-central1",
		RegionName: "Iowa",
		Lat:        41.8780025,
		Lng:        -93.097702,
	},
	"us-east1": {
		URL:        "https://us-east1-5tkroniexa-ue.a.run.app",
		Region:     "us-east1",
		RegionName: "South Carolina",
		Lat:        33.836081,
		Lng:        -81.1637245,
	},
	"us-east4": {
		URL:        "https://us-east4-5tkroniexa-uk.a.run.app",
		Region:     "us-east4",
		RegionName: "North Virginia",
		Lat:        32.817108,
		Lng:        -96.94944799999999
	},
	"us-west1": {
		URL:        "https://us-west1-5tkroniexa-uw.a.run.app",
		Region:     "us-west1",
		RegionName: "Oregon",
		Lat:        34.0522342,
		Lng:        -118.2436849,
	},
	"us-west2": {
		URL:        "https://us-west2-5tkroniexa-wl.a.run.app",
		Region:     "us-west2",
		RegionName: "Los Angeles",
		Lat:        34.0522342,
		Lng:        -118.2436849,
	},
	"us-west3": {
		URL:        "https://us-west3-5tkroniexa-wm.a.run.app",
		Region:     "us-west3",
		RegionName: "Salt Lake City",
		Lat:        40.7607793,
		Lng:        -111.8910474,
	},
	"us-west4": {
		URL:        "https://us-west4-5tkroniexa-wn.a.run.app",
		Region:     "us-west4",
		RegionName: "Las Vegas",
		Lat:        36.1699412,
		Lng:        -115.1398296,
	},
}
