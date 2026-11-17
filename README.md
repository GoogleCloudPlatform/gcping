# gcping

[![Build Status](https://github.com/GoogleCloudPlatform/gcping/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/GoogleCloudPlatform/gcping/actions/workflows/tests.yml)

> **🛑 Service Discontinued & Repository Archived:** The hosted `gcping.com` web application, `global.gcping.com/api/endpoints` service, regional ping endpoints, and pre-built CLI binary downloads (`gs://gcping-release`) have been shut down, and this repository is now archived in read-only mode.
>
> To measure network latency to Google Cloud regions, please use the official [Google Cloud Network Intelligence Center Performance Dashboard](https://docs.cloud.google.com/network-intelligence-center/docs/performance-dashboard/concepts/overview). If you wish to run GCPing yourself, the open-source code and Terraform configuration below remain available to deploy and self-host in your own Google Cloud project.

gcpping is both a command line tool and a webapp that reports median latency to
Google Cloud regions.

Note: This is not an official Google product.

## CLI Usage

```
gcping [options...]

Options:
-n       Number of requests to be made to each region.
         By default 10; can't be negative.
-c       Max number of requests to be made at any time.
         By default 10; can't be negative or zero.
-r       Report latency for an individual region.
-t       Timeout. By default, no timeout.
         Examples: "500ms", "1s", "1s500ms".
-top     If true, only the top (non-global) region is printed.
-csv-cum If true, cumulative value is printed in CSV; disables default report.
-url     URL of endpoint list. (Self-hosted deployments must pass their own endpoints URL.)

-csv     CSV output; disables verbose output.
-v       Verbose output.
```

An example output:

```
$ gcping -url https://example.com/api/endpoints
 1.  [global]                   11.17568ms
 2.  [us-central1]              12.373109ms
 3.  [us-west3]                 29.203499ms
 4.  [northamerica-northeast2]  30.615139ms
 5.  [us-east4]                 33.401098ms
 6.  [northamerica-northeast1]  38.612769ms
 7.  [us-west1]                 43.041808ms
 8.  [us-east1]                 46.847258ms
 9.  [us-west4]                 53.438688ms
10.  [us-west2]                 57.659108ms
11.  [europe-west2]             103.371016ms
12.  [europe-west4]             111.966565ms
13.  [europe-west1]             112.327356ms
14.  [europe-west3]             114.245525ms
15.  [europe-west6]             118.966225ms
16.  [europe-central2]          128.008935ms
17.  [europe-north1]            136.796505ms
18.  [asia-northeast1]          142.480775ms
19.  [southamerica-east1]       147.324384ms
20.  [asia-northeast2]          156.088594ms
21.  [asia-northeast3]          168.205243ms
22.  [asia-east2]               170.763954ms
23.  [australia-southeast2]     188.310153ms
24.  [southamerica-west1]       206.412352ms
25.  [asia-southeast2]          210.029872ms
26.  [asia-south1]              256.2628ms
27.  [asia-south2]              276.434709ms
28.  [australia-southeast1]     396.915245ms
29.  [asia-east1]               417.147963ms
30.  [asia-southeast1]          496.648151ms
```

```
$ gcping -url https://example.com/api/endpoints -r us-east1
502.068712ms
```

```
$ gcping -url https://example.com/api/endpoints -top
us-west2
```

## Building from Source & Self-Hosting

Pre-built binaries at `storage.googleapis.com/gcping-release` have been retired. To build the CLI from source with Go 1.21+:

```
go install github.com/GoogleCloudPlatform/gcping@latest
```

See [DEVELOPMENT.md](DEVELOPMENT.md) and [`tools/terraform/`](tools/terraform/) for instructions on deploying your own regional Cloud Run services and endpoints configuration.
