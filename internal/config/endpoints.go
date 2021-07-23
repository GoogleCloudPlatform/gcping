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
		URL:        "https://asia-east1-bmlfzs4h6a-de.a.run.app",
		Region:     "asia-east1",
		RegionName: "Taiwan",
	},
	"asia-east2": {
		URL:        "https://asia-east2-bmlfzs4h6a-df.a.run.app",
		Region:     "asia-east2",
		RegionName: "Hong Kong",
	},
	"asia-northeast1": {
		URL:        "https://asia-northeast1-bmlfzs4h6a-an.a.run.app",
		Region:     "asia-northeast1",
		RegionName: "Tokyo",
	},
	"asia-northeast2": {
		URL:        "https://asia-northeast2-bmlfzs4h6a-dt.a.run.app",
		Region:     "asia-northeast2",
		RegionName: "Osaka",
	},
	"asia-northeast3": {
		URL:        "https://asia-northeast3-bmlfzs4h6a-du.a.run.app",
		Region:     "asia-northeast3",
		RegionName: "Seoul",
	},
	"asia-south1": {
		URL:        "https://asia-south1-bmlfzs4h6a-el.a.run.app",
		Region:     "asia-south1",
		RegionName: "Mumbai",
	},
	"asia-south2": {
		URL:        "https://asia-south2-ezn5kimndq-em.a.run.app",
		Region:     "asia-south2",
		RegionName: "Delhi",
	},
	"asia-southeast1": {
		URL:        "https://asia-southeast1-bmlfzs4h6a-as.a.run.app",
		Region:     "asia-southeast1",
		RegionName: "Singapore",
	},
	"asia-southeast2": {
		URL:        "https://asia-southeast2-bmlfzs4h6a-et.a.run.app",
		Region:     "asia-southeast2",
		RegionName: "Jakarta",
	},
	"australia-southeast1": {
		URL:        "https://australia-southeast1-ezn5kimndq-ts.a.run.app",
		Region:     "australia-southeast1",
		RegionName: "Sydney",
	},
	"australia-southeast2": {
		URL:        "https://australia-southeast2-ezn5kimndq-km.a.run.app",
		Region:     "australia-southeast2",
		RegionName: "Melbourne",
	},
	"europe-central2": {
		URL:        "https://europe-central2-ezn5kimndq-lm.a.run.app",
		Region:     "europe-central",
		RegionName: "",
	},
	"europe-north1": {
		URL:        "https://europe-north1-bmlfzs4h6a-lz.a.run.app",
		Region:     "europe-north1",
		RegionName: "Finland",
	},
	"europe-west1": {
		URL:        "https://europe-west1-bmlfzs4h6a-ew.a.run.app",
		Region:     "europe-west1",
		RegionName: "Belgium",
	},
	"europe-west2": {
		URL:        "https://europe-west2-bmlfzs4h6a-nw.a.run.app",
		Region:     "europe-west2",
		RegionName: "London",
	},
	"europe-west3": {
		URL:        "https://europe-west3-bmlfzs4h6a-ey.a.run.app",
		Region:     "europe-west3",
		RegionName: "Frankfurt",
	},
	"europe-west4": {
		URL:        "https://europe-west4-bmlfzs4h6a-ez.a.run.app",
		Region:     "europe-west4",
		RegionName: "Netherlands",
	},
	"europe-west6": {
		URL:        "https://europe-west6-bmlfzs4h6a-oa.a.run.app",
		Region:     "europe-west6",
		RegionName: "Zurich",
	},
	"northamerica-northeast1": {
		URL:        "https://northamerica-northeast1-bmlfzs4h6a-nn.a.run.app",
		Region:     "northamerica-northeast1",
		RegionName: "Montréal",
	},
	"southamerica-east1": {
		URL:        "https://southamerica-east1-bmlfzs4h6a-rj.a.run.app",
		Region:     "southamerica-east1",
		RegionName: "São Paulo",
	},
	"us-central1": {
		URL:        "https://us-central1-bmlfzs4h6a-uc.a.run.app",
		Region:     "us-central1",
		RegionName: "Iowa",
	},
	"us-east1": {
		URL:        "https://us-east1-bmlfzs4h6a-ue.a.run.app",
		Region:     "us-east1",
		RegionName: "South Carolina",
	},
	"us-east4": {
		URL:        "https://us-east4-bmlfzs4h6a-uk.a.run.app",
		Region:     "us-east4",
		RegionName: "N. Virgina",
	},
	"us-west1": {
		URL:        "https://us-west1-bmlfzs4h6a-uw.a.run.app",
		Region:     "us-west1",
		RegionName: "Oregon",
	},
	"us-west2": {
		URL:        "https://us-west2-ezn5kimndq-wl.a.run.app",
		Region:     "us-west2",
		RegionName: "Los Angeles",
	},
	"us-west3": {
		URL:        "https://us-west3-ezn5kimndq-wm.a.run.app",
		Region:     "us-west3",
		RegionName: "Salt Lake City",
	},
	"us-west4": {
		URL:        "https://us-west4-ezn5kimndq-wn.a.run.app",
		Region:     "us-west4",
		RegionName: "Las Vegas",
	},
}
