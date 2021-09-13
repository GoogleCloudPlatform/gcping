# gcping

[![Build Status](https://github.com/GoogleCloudPlatform/gcping/actions/workflows/tests.yml/badge.svg)](https://github.com/GoogleCloudPlatform/gcping/actions/workflows/tests.yml)


gcpping is both a command line tool and a webapp that reports median latency to
Google Cloud regions. The webapp is hosted at [gcping.com](http://gcping.com).
To install the command line tool, see below.

Note: This is not an official Google product.

## CLI Usage

```
gcping [options...]

Options:
-n   Number of requests to be made to each region.
     By default 10; can't be negative.
-c   Max number of requests to be made at any time.
     By default 10; can't be negative or zero.
-r   Report latency for an individual region.
-t   Timeout. By default, no timeout.
     Examples: "500ms", "1s", "1s500ms".
-top If true, only the top (non-global) region is printed.

-csv CSV output; disables verbose output.
-v   Verbose output.

Need a website version? See gcping.com
```

An example output:

```
$ gcping
 1.  [global]                   36.752191ms
 2.  [us-east4]                 37.091976ms
 3.  [northamerica-northeast1]  51.918669ms
 4.  [us-central1]              75.488941ms
 5.  [us-east1]                 75.928857ms
 6.  [us-west2]                 148.998964ms
 7.  [us-west1]                 157.899518ms
 8.  [europe-west2]             166.42703ms
 9.  [europe-west1]             174.226927ms
10.  [europe-west4]             179.802812ms
11.  [europe-west3]             195.430189ms
12.  [europe-west6]             208.143331ms
13.  [europe-north1]            252.823482ms
14.  [southamerica-east1]       311.575344ms
15.  [asia-northeast1]          338.151472ms
16.  [asia-northeast2]          358.787403ms
17.  [asia-east1]               394.165761ms
18.  [asia-east2]               418.293092ms
19.  [australia-southeast1]     425.679503ms
20.  [asia-southeast1]          454.494659ms
21.  [asia-south1]              573.022571ms
```

```
$ gcping -r us-east1
502.068712ms
```

```
$ gcping -top
us-west2
```

## Installation

We build binaries for the following OS's and architectures:

* Linux 64-bit: https://storage.googleapis.com/gcping-release/gcping_linux_amd64_latest
* Mac 64-bit (x86): https://storage.googleapis.com/gcping-release/gcping_darwin_amd64_latest
* Mac 64-bit (Apple Silicon): https://storage.googleapis.com/gcping-release/gcping_darwin_arm64_latest
* Windows 64-bit: https://storage.googleapis.com/gcping-release/gcping_windows_amd64_latest

Installation looks something like this (changing the URL for your system):

```
curl https://storage.googleapis.com/gcping-release/gcping_linux_amd64_latest > gcping && chmod +x gcping
```
