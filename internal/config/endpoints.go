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
	"fmt"
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
// using json served by the gcping endpoints.
func GenerateConfigFromEndpoints(ctx context.Context, endpointsURL string) (map[string]Endpoint, error) {

	endpointsMap := make(map[string]Endpoint)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpointsURL,
		nil,
	)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return endpointsMap, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%v %s", resp.Status, endpointsURL)
		return endpointsMap, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&endpointsMap); err != nil {
		return endpointsMap, err
	}

	return endpointsMap, err
}

// GenerateConfigFromAPI is used to generate the endpoint config through the
// metadata provided by the Cloud Run Admin API.
func GenerateConfigFromAPI(ctx context.Context) (map[string]Endpoint, error) {
	var endpointsMap = make(map[string]Endpoint)
	r, err := run.NewService(ctx)
	// TODO: Get project name from Cloud Run metadata service if not defined in env variable
	//projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	projectID := "kr-gcping"
	if projectID == "" {
		err := fmt.Errorf("could not retrieve Google Cloud Project ID from $GOOGLE_CLOUD_PROJECT")
		return endpointsMap, err
	}

	// List Services
	resp, err := r.Namespaces.Services.List("namespaces/"+projectID).
		Fields("items(status/address/url",
			"metadata(labels,name)",
			"spec(template/metadata/annotations))").
		LabelSelector("env=prod").Do()

	for _, v := range resp.Items {
		e := Endpoint{
			URL:        v.Status.Address.Url,
			Region:     v.Metadata.Labels["cloud.googleapis.com/location"],
			RegionName: v.Spec.Template.Metadata.Annotations["region-name"],
		}
		endpointsMap[e.Region] = e
	}

	// Add global endpoint to map if env is defined.
	globalURL := os.Getenv("GLOBAL_ENDPOINT")
	if globalURL != "" {
		g := Endpoint{
			URL:        globalURL,
			Region:     "global",
			RegionName: "Global External HTTPS Load Balancer",
		}
		endpointsMap[g.Region] = g
	}

	return endpointsMap, err
}
