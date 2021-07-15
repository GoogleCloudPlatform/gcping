////// Providers

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.54.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "3.54.0"
    }
  }
}

provider "google" {
  project = var.project
}

provider "google-beta" {
  project = var.project
}

////// Variables

variable "image" {
  type = string
}

variable "project" {
  type    = string
  default = "gcping"
}

variable "domain" {
  type    = string
  default = "gcping.com"
}

data "google_cloud_run_locations" "available" {
}

resource "google_service_account" "minimal" {
  account_id = "minimal-service-account"
  display_name = "Minimal Service Account"
}

////// Cloud Run

// Enable Cloud Run API.
resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

// Enable Compute Engine API.
resource "google_project_service" "compute" {
  service = "compute.googleapis.com"
}

// Deploy image to each region.
resource "google_cloud_run_service" "regions" {
  for_each = toset(data.google_cloud_run_locations.available.locations)
  name     = each.key
  location = each.key

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "3" // Control costs.
        "run.googleapis.com/launch-stage"  = "BETA"
        // This gets added and causes diffs, but must be removed before adding a new service...
        "run.googleapis.com/sandbox"       = "gvisor"
      }
    }
    spec {
      service_account_name = google_service_account.minimal.email
      containers {
        image = var.image
        env {
          name  = "REGION"
          value = each.key
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.run]
}

// Print each service URL.
output "services" {
  value = {
    for svc in google_cloud_run_service.regions :
    svc.name => svc.status[0].url
  }
}

// Make each service invokable by all users.
resource "google_cloud_run_service_iam_member" "allUsers" {
  for_each = toset(data.google_cloud_run_locations.available.locations)

  service  = google_cloud_run_service.regions[each.key].name
  location = each.key
  role     = "roles/run.invoker"
  member   = "allUsers"

  depends_on = [google_cloud_run_service.regions]
}

// Create a regional network endpoint group (NEG) for each regional Cloud Run service.
resource "google_compute_region_network_endpoint_group" "regions" {
  for_each = toset(data.google_cloud_run_locations.available.locations)

  name                  = each.key
  network_endpoint_type = "SERVERLESS"
  region                = each.key
  cloud_run {
    service = google_cloud_run_service.regions[each.key].name
  }

  depends_on = [google_project_service.compute]
}

////// Global Domain + Load Balancer config

// Reserve a global static IP address.
resource "google_compute_global_address" "global" {
  name = "address"
}

resource "google_compute_global_forwarding_rule" "global" {
  name       = "global"
  target     = google_compute_target_https_proxy.global.id
  port_range = "443"
  ip_address = google_compute_global_address.global.address
}

// Print global LB IP address.
output "global" {
  value = google_compute_global_address.global.address
}

resource "google_compute_managed_ssl_certificate" "global" {
  provider = google-beta

  name = "global"
  managed {
    domains = [
      "global.${var.domain}.",
      "${var.domain}.",
    ]
  }
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
