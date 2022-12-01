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


// Create a bucket for CLI releases
resource "google_storage_bucket_object" "endpoints_config" {
  // TODO: set path to /api/endpoints once LB url maps and configs validated
  name         = "api/gcs-endpoints"
  bucket       = google_storage_bucket.config.id
  content_type = "application/json"

  content = jsonencode(
    merge(
      zipmap(
        keys(local.regions),
        [for k, v in local.regions :
          {
            URL        = google_cloud_run_service.regions[k].status[0].url,
            RegionName = v,
            Region     = k
          }
        ]
      ),
      zipmap(
        ["global"],
        [
          {
            URL        = "https://global.${var.domain}",
            Region     = "global",
            RegionName = "Global External HTTPS Load Balancer"
          }
        ]
      )
    )
  )

}
