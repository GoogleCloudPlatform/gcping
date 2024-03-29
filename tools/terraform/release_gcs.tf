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

locals {
  release_bucket = var.release_bucket != "" ? var.release_bucket : "${var.project}${var.release_bucket_suffix}"
}

// Create a bucket for CLI releases
resource "google_storage_bucket" "releases" {
  name                        = local.release_bucket
  uniform_bucket_level_access = true
  location                    = "US"
}

// Make the bucket publically accessible
resource "google_storage_bucket_iam_member" "public_access" {
  bucket = google_storage_bucket.releases.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
