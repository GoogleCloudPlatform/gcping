#!/bin/bash

# Cloud Shell and IPv6 don't play well together. This script patches /etc/hosts
# to use IPv4 for all googleapis we use.
#
# See the following issue for more details:
# https://github.com/hashicorp/terraform-provider-google/issues/6782#issuecomment-874574409

export APIS="googleapis.com
cloudresourcemanager.googleapis.com
container.googleapis.com
iam.googleapis.com
iam.googleapis.com
run.googleapis.com
storage.googleapis.com
www.googleapis.com
asia-east1-run.googleapis.com
asia-east2-run.googleapis.com
asia-northeast1-run.googleapis.com
asia-northeast2-run.googleapis.com
asia-northeast3-run.googleapis.com
asia-northest2-run.googleapis.com
asia-south1-run.googleapis.com
asia-southeast1-run.googleapis.com
asia-southeast2-run.googleapis.com
australia-southeast1-run.googleapis.com
australia-southeast2-run.googleapis.com
europe-central2-run.googleapis.com
europe-north1-run.googleapis.com
europe-west1-run.googleapis.com
europe-west2-run.googleapis.com
europe-west3-run.googleapis.com
europe-west3-run.googleapis.com
europe-west4-run.googleapis.com
europe-west10-run.googleapis.com
northamerica-northeast1-run.googleapis.com
southamerica-east1-run.googleapis.com
southamerica-east1.googleapis.com
us-central1-run.googleapis.com
us-east4-run.googleapis.com
us-east4-run.googleapis.com
us-west2-run.googleapis.com
us-west3-run.googleapis.com
us-west4-run.googleapis.com"

for name in $APIS; do
    echo Configuring IPv4 for $name...
    ipv4=$(getent ahostsv4 "$name" | head -n 1 | awk '{ print $1 }')
    grep -q "$name" /etc/hosts || ([ -n "$ipv4" ] && sudo sh -c "echo '$ipv4 $name' >> /etc/hosts")
done
