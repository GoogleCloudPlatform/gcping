# Development

Deploy Cloud Run site using Terraform

### Prerequisites

Install
[Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) and
[`ko`](https://github.com/google/ko), and set the `KO_DOCKER_REPO` env var to
the GCR repository you'd like to deploy to (e.g.,
`KO_DOCKER_REPO=gcr.io/gcping`)

The frontend requires [Node.js](https://nodejs.org/en/).

### Build the frontend

```
$ cd web
$ npm install
$ npm run build  # generate the frontend
```

### Deploy using Terraform

```
$ gcloud auth login                      # Used by ko
$ gcloud auth application-default login  # Used by Terraform
```

```
$ terraform init # necessary only the first time
$ terraform apply -var image=$(ko publish -P ./cmd/ping/)
```

This deploys the ping service to all Cloud Run regions and configures a global
HTTPS Load Balancer with Google-managed SSL certificate for
`global.gcping.com`.

### Run frontend locally

First, start the Ping server with:

``` shell
# starts a server on localhost:8080
go run ./cmd/ping/main.go
```

Next, within the `web` directory, run:

``` shell
# starts a frontend server at localhost:1234
npm start
```

The frontend server when run locally is configured to proxy all API requests to
`localhost:8080`.

