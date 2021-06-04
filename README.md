# gcping

[![Build Status](https://travis-ci.com/GoogleCloudPlatform/gcping.svg?branch=master)](https://travis-ci.com/GoogleCloudPlatform/gcping)

gcpping is both a command line tool and a webapp that reports median latency to
Google Cloud regions. The webapp is hosted at [gcping.com](http://gcping.com).

Note: This is not an official Google product.

This began as a side project by [@imjasonh][] who built the first versions of the
webapp. Later [@rakyll][] build the CLI version.

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

* Linux 64-bit: https://storage.googleapis.com/gcping-release/gcping_linux_amd64_0.0.3
  ```
  $ curl https://storage.googleapis.com/gcping-release/gcping_linux_amd64_0.0.3 > gcping && chmod +x gcping
  ```
* Mac 64-bit: https://storage.googleapis.com/gcping-release/gcping_darwin_amd64_0.0.3
* Windows 64-bit: https://storage.googleapis.com/gcping-release/gcping_windows_amd64_0.0.3

[@imjasonh]: https://github.com/imjasonh
[@rakyll]:   https://github.com/rakyll
