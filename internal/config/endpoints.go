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
	// deloyed, e.g., us-central1.
	Region string
	// RegionName is the geographic name of the region, e.g., Iowa.
	RegionName string
}

func GenerateConfigFromEndpoints(ctx context.Context) map[string]Endpoint {
	//Used by CLI to pull endpoint configs from Cloud Run enpoints.

	type innerEndpoint struct {
		//Used to unmarshal properly formatted JSON obtained from /api/endpoint JSON
		URL        string
		Region     string
		RegionName string
	}

	e := make(map[string]Endpoint)
	ie := new(map[string]innerEndpoint)

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
	json.NewDecoder(resp.Body).Decode(ie)

	for _, s := range *ie {
		e[s.Region] = Endpoint(s)
	}
	return e
}

func GenerateConfigFromAPI(ctx context.Context) (map[string]Endpoint, error) {
	//Used by Cloud Run Endpoints to pull endpoint configs from Cloud Run Admin API
	log.Print("Using Cloud Run Admin API to generate Endpoints config.")
	runService, err := run.NewService(ctx)
	// TODO: Get project name from Cloud Run metadata service if not defined in env variable
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// List Services
	resp, err := runService.Namespaces.Services.List("namespaces/" + projectID).Fields("items(status/address/url,metadata(labels,name),spec(template/metadata/annotations))").LabelSelector("env=prod").Do()

	s, _ := json.MarshalIndent(resp.Items, "", "\t")

	var endpoints []Endpoint

	json.Unmarshal([]byte(s), &endpoints)

	var EndpointsMap = make(map[string]Endpoint)

	// Add global endpoint to map if env is defined.
	globalURL := os.Getenv("GLOBAL_ENDPOINT")
	if globalURL != "" {
		var Global Endpoint
		Global.Region = "global"
		Global.RegionName = "Global External HTTPS Load Balancer"
		Global.URL = os.Getenv("GLOBAL_ENDPOINT")
		EndpointsMap[Global.Region] = Global
	}

	for _, s := range endpoints {
		EndpointsMap[s.Region] = s
	}

	return EndpointsMap, err
}

func (es *Endpoint) UnmarshalJSON(data []byte) error {
	//Custom unmarshal function for Endpoint data
	type labelsInner struct {
		Location string `json:"cloud.googleapis.com/location"`
	}

	type annotationsInner struct {
		RegionName string `json:"region-name"`
	}

	type templateMetadataInner struct {
		Annotations annotationsInner `json:"annotations`
	}

	type templateInner struct {
		Metadata templateMetadataInner `json:"metadata"`
	}

	type specInner struct {
		Template templateInner `json:"template"`
	}

	type metadataInner struct {
		Labels labelsInner `json:"labels"`
		Name   string      `json:"name"`
	}

	type addressInner struct {
		URL string `json:"url"`
	}

	type statusInner struct {
		Address addressInner `json:"address"`
	}

	type nestedService struct {
		Metadata metadataInner `json:"metadata"`
		Status   statusInner   `json:"status"`
		Spec     specInner     `json:"spec"`
	}

	var ns nestedService

	if err := json.Unmarshal(data, &ns); err != nil {
		return err
	}

	// create the struct in desired format
	tmp := &Endpoint{
		URL:        ns.Status.Address.URL,
		Region:     ns.Metadata.Labels.Location,
		RegionName: ns.Spec.Template.Metadata.Annotations.RegionName,
	}

	// reassign the method receiver pointer
	// to store the values in the struct
	*es = *tmp
	return nil
}
