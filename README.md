# gcping

[![Build Status](https://travis-ci.com/GoogleCloudPlatform/gcping.svg?branch=master)](https://travis-ci.com/GoogleCloudPlatform/gcping)

gcping is a command line tools that reports latency to
Google Cloud regions. It is inspired by [gcping.com](https://gcping.com).

An example output:

```
$ gcping
0. [global] 17.282896ms
1. [us-west2] 24.831978ms
2. [us-west1] 37.200678ms
3. [us-central1] 62.908829ms
4. [us-east4] 76.461448ms
5. [us-east1] 87.213011ms
6. [northamerica-northeast1] 87.672202ms
7. [asia-northeast1] 128.187124ms
8. [asia-northeast2] 134.71903ms
9. [europe-west2] 151.353133ms
10. [europe-west1] 152.090048ms
11. [asia-east1] 154.469779ms
12. [europe-west4] 156.138144ms
13. [europe-west3] 161.17031ms
14. [australia-southeast1] 161.458364ms
15. [asia-east2] 167.030351ms
16. [europe-west6] 168.109585ms
17. [southamerica-east1] 186.049558ms
18. [europe-north1] 186.685738ms
19. [asia-southeast1] 189.678346ms
20. [asia-south1] 248.994198ms
```

## Installation

* Linux 64-bit: https://storage.googleapis.com/gcping-release/gcping_linux_amd64
  ```
  $ curl https://storage.googleapis.com/gcping-release/gcping_linux_amd64 > gcping && chmod +x gcping
  ```
* Mac 64-bit: https://storage.googleapis.com/gcping-release/gcping_darwin_amd64
* Windows 64-bit: https://storage.googleapis.com/gcping-release/gcping_windows_amd64

Note: This is not an official Google product.
