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

resource "google_cloudbuild_trigger" "pr-trigger" {
    depends_on = [
        google_project_service.iam,
        google_service_account.gcb
    ]

    provider = google-beta

    service_account = google_service_account.gcb.id
    project         = var.project
    name            = "pr-validation-staging"
    description     = "Build and deploy to a staging endpoint"
    filename        = "tools/cloudbuild/pr-open.yaml"

  github {
    owner = var.github_org
    name  = var.github_repo
    pull_request {
      branch = "^main$"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
    }
  }
}
