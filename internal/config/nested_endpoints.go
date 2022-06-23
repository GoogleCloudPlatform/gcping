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

// nestedEndpoint struct and the associated structs describe the format of the
// json object retrieved from the Cloud Run Admin API, and enable to unmarshal
// the retreoved json object for unnesting.
type nestedEndpoint struct {
	Metadata nestedMetadata `json:"metadata"`
	Status   nestedStatus   `json:"status"`
	Spec     nestedSpec     `json:"spec"`
}

type nestedMetadata struct {
	Labels nestedLabels `json:"labels"`
	Name   string       `json:"name"`
}

type nestedStatus struct {
	Address nestedAddress `json:"address"`
}

type nestedSpec struct {
	Template nestedTemplate `json:"template"`
}

type nestedLabels struct {
	Location string `json:"cloud.googleapis.com/location"`
}

type nestedAddress struct {
	URL string `json:"url"`
}

type nestedTemplate struct {
	Metadata nestedTemplateMetadata `json:"metadata"`
}

type nestedTemplateMetadata struct {
	Annotations nestedAnnotations `json:"annotations"`
}

type nestedAnnotations struct {
	RegionName string `json:"region-name"`
}

// unNestEndpoint is used to unnest the nested JSON objects retreived from the
// Cloud Run Admin API into an Endpoint struct
func unNestEndpoint(ne nestedEndpoint) Endpoint {
	e := Endpoint{
		URL:        ne.Status.Address.URL,
		Region:     ne.Metadata.Labels.Location,
		RegionName: ne.Spec.Template.Metadata.Annotations.RegionName,
	}
	return e
}
