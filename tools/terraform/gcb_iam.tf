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

// Enable IAM API
resource "google_project_service" "iam" {
  service = "iam.googleapis.com"
}

// Create custom Service Account for Cloud Build
resource "google_service_account" "gcb" {
  account_id   = "cloudbuild"
  display_name = "Service account for Cloud Build jobs"
}

resource "google_project_iam_member" "gcb_run_admin" {
  project = var.project
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.gcb.email}"
}

resource "google_project_iam_member" "gcb_builder" {
  project = var.project
  role    = "roles/cloudbuild.builds.builder"
  member  = "serviceAccount:${google_service_account.gcb.email}"
}

resource "google_project_iam_member" "gcb_sa_user" {
  project = var.project
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.gcb.email}"
}
