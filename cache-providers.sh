#!/bin/sh

set -ex

mkdir -p /root/.terraform.d/plugins

providers="google:5.40.0 google-beta:5.40.0 random:3.6.2"

for provider in $providers; do
    name=$(echo "$provider" | cut -d: -f1)
    version=$(echo "$provider" | cut -d: -f2)
    curl -LO "https://releases.hashicorp.com/terraform-provider-${name}/${version}/terraform-provider-${name}_${version}_linux_amd64.zip"
    unzip -o "terraform-provider-${name}_${version}_linux_amd64.zip" -d /root/.terraform.d/plugins/
    rm "terraform-provider-${name}_${version}_linux_amd64.zip"
done