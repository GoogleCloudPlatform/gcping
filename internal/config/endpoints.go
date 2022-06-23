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

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/run/v1"
)

// Endpoint represents a Cloud Run service deploy in a particular region.
type Endpoint struct {
	// URL is the HTTPS URL of the service
	URL string
	// Region is the programmatic name of the region where the endpoint is
	// deployed, e.g., us-central1.
	Region string
	// RegionName is the geographic name of the region, e.g., Iowa.
	RegionName string
}

// GenerateConfigFromEndpoints is used by the cli to generate an Endpoint map
// using a precompiled list served by the gcping endpoints.
func GenerateConfigFromEndpoints(ctx context.Context) map[string]Endpoint {

	EndpointsMap := make(map[string]Endpoint)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://global.gcping.com/api/endpoints",
		nil,
	)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&EndpointsMap)

	return EndpointsMap
}

// GenerateConfigFromAPI is used to generate the endpoint config through the
// metadat provided by the Cloud Run Admin API.
func GenerateConfigFromAPI(ctx context.Context) (map[string]Endpoint, error) {
	log.Print("Using Cloud Run Admin API to generate Endpoints config.")
	runService, err := run.NewService(ctx)
	// TODO: Get project name from Cloud Run metadata service if not defined in env variable
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// List Services
	resp, err := runService.Namespaces.Services.List("namespaces/" + projectID).Fields("items(status/address/url,metadata(labels,name),spec(template/metadata/annotations))").LabelSelector("env=prod").Do()

	s, _ := json.MarshalIndent(resp.Items, "", "\t")

	var nestedEndpointsMap []nestedEndpoint
	json.Unmarshal(s, &nestedEndpointsMap)
	var EndpointsMap = make(map[string]Endpoint)

	// Add global endpoint to map if env is defined.
	globalURL := os.Getenv("GLOBAL_ENDPOINT")
	if globalURL != "" {
		Global := Endpoint{
			URL:        os.Getenv("GLOBAL_ENDPOINT"),
			Region:     "global",
			RegionName: "Global External HTTPS Load Balancer",
		}
		EndpointsMap[Global.Region] = Global
	}

	for _, nestedEndpoint := range nestedEndpointsMap {
		e := unNestEndpoint(nestedEndpoint)
		EndpointsMap[e.Region] = e
	}
	return EndpointsMap, err
}
