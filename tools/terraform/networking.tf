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
  managed_domains = var.domain_alias_flag ? [
    "www.${var.domain}",
    "global.${var.domain}",
    "${var.domain}",
    "www.${var.domain_alias}",
    "${var.domain_alias}",
    ] : [
    "www.${var.domain}",
    "global.${var.domain}",
    "${var.domain}",
  ]
}

// Enable Compute Engine API.
resource "google_project_service" "compute" {
  service = "compute.googleapis.com"
}


// Reserve a global static IP address.
resource "google_compute_global_address" "global" {
  name = "address"

  depends_on = [
    google_project_service.compute
  ]
}

resource "google_compute_global_forwarding_rule" "global" {
  name       = "global"
  target     = google_compute_target_https_proxy.global.id
  port_range = "443"
  ip_address = google_compute_global_address.global.address
  depends_on = [
    google_project_service.compute
  ]
}

resource "google_compute_target_https_proxy" "global" {
  provider = google-beta

  name             = "global"
  url_map          = google_compute_url_map.global.id
  ssl_certificates = [google_compute_managed_ssl_certificate.global.id]
}

resource "google_compute_url_map" "global" {

  provider = google-beta

  name            = "global"
  description     = "a description"
  default_service = google_compute_backend_service.global.id


  // Create a host rule to match traffic to alias (gcpping.com)
  host_rule {
    hosts        = ["*"]
    path_matcher = "endpoints-config-bucket"
  }

  dynamic "host_rule" {
    for_each = var.domain_alias_flag ? [1] : []

    content {
      hosts = [
        var.domain_alias,
      ]
      path_matcher = "alt-redirect"
    }
  }

  path_matcher {
    name            = "endpoints-config-bucket"
    default_service = google_compute_backend_service.global.self_link

    path_rule {
      // TODO: set path to /api/endpoints once LB url maps and configs validated
      paths   = ["/api/gcs-endpoints"]
      service = google_compute_backend_bucket.endpoints_backend.id
    }
  }
  // 301 redirect traffic from gcpping.com to gcping.com
  dynamic "path_matcher" {
    for_each = var.domain_alias_flag ? [1] : []

    content {
      name = "alt-redirect"



      default_url_redirect {
        host_redirect          = var.domain
        https_redirect         = false
        redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
        strip_query            = false
      }
    }
  }
  depends_on = [
    google_project_service.compute
  ]
}

// Create a global backend service with a backend for each regional NEG.
resource "google_compute_backend_service" "global" {
  name       = "global"
  enable_cdn = true

  // Add a backend for each regional NEG.
  dynamic "backend" {
    for_each = google_compute_region_network_endpoint_group.regions
    content {
      group = backend.value["id"]
    }
  }
}

// Create an HTTP->HTTPS upgrade rule.
resource "google_compute_url_map" "https_redirect" {
  name = "https-redirect"

  default_url_redirect {
    https_redirect         = true
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    strip_query            = false
  }
  depends_on = [
    google_project_service.compute
  ]

}

resource "google_compute_target_http_proxy" "https_redirect" {
  name    = "https-redirect"
  url_map = google_compute_url_map.https_redirect.id
}

resource "google_compute_global_forwarding_rule" "https_redirect" {
  name = "https-redirect"

  target     = google_compute_target_http_proxy.https_redirect.id
  port_range = "80"
  ip_address = google_compute_global_address.global.address
}


// create a reandom id for the SSL cert
resource "random_id" "certificate" {
  byte_length = 2
  prefix      = "global-"
  keepers = {
    domains = join(",", local.managed_domains)
  }
}

// create a managed SSL cert
resource "google_compute_managed_ssl_certificate" "global" {
  provider = google-beta

  name = random_id.certificate.hex
  managed {
    domains = local.managed_domains
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    google_project_service.compute
  ]
}
