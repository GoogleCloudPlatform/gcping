# Development

Deploy Cloud Run site using Terraform

### Prerequisites

Install
[Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) and
[`ko`](https://github.com/google/ko), and set the `KO_DOCKER_REPO` env var to
the GCR repository you'd like to deploy to (e.g.,
`KO_DOCKER_REPO=gcr.io/gcping`)

### Deploy using Terraform

```
$ gcloud auth login                      # Used by ko
$ gcloud auth application-default login  # Used by Terraform
```

```
$ terraform init
$ terraform apply -var image=$(ko publish -P ./cmd/ping/)
```

This deploys the ping service to all Cloud Run regions and configures a global
HTTPS Load Balancer with Google-managed SSL certificate for
`global.gcping.com`.

### Run frontend locally

```
docker run -p 8080:8080 $(KO_DOCKER_REPO=ko.local ko publish -P ./cmd/ping/)
```

And browse to http://localhost:8080/

This connects to real regional backends and the global LB backend.

### Regenerate list of URLs

```
./genconfig.sh
```

This transforms Terraform output to a form usable by the JS frontend. If there
are any changes, you need to redeploy to see them.
