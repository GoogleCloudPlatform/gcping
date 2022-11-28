// Copyright 2022 Google Inc. All Rights Reserved.
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

// Print each service URL.
output "services" {
  value = {
    for svc in google_cloud_run_service.regions :
    svc.name => svc.status[0].url
  }
}

// Prepares a config output in JSON format compatible with output of /api/endpoints
output "config" {
  value = merge(zipmap(keys(google_cloud_run_service.regions), [for k, v in google_cloud_run_service.regions : { URL = v.status[0].url, RegionName = google_cloud_run_service.regions[k].template[0].metadata[0].annotations["gcping.com/region-name"], Region = k }]), zipmap(["global"], [{ URL = "https://global.${var.domain}", Region = "global", RegionName = "Global External HTTPS Load Balancer" }]))
}

// Print global LB IP address.
output "global" {
  value = google_compute_global_address.global.address
}
