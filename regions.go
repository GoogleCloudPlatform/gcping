package main

import "sort"

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

// getSortedRegions sorts and returns a list of regions.
func getSortedRegions() []string {
	regions := make([]string, 0, len(endpoints))
	for region := range endpoints {
		regions = append(regions, region)
	}
	sort.Strings(regions) //sort regions alphabetically.
	return regions
}
