#!/usr/bin/env bash

set -e

wget https://dl.smallstep.com/cli/docs-cli-install/latest/step-cli_amd64.deb
sudo apt-get install -y -f ./step-cli_amd64.deb

mkdir tls/

step certificate create \
    ca.internal \
    tls/ca.crt tls/ca.key \
    --profile root-ca \
    --no-password \
    --insecure \
    --not-after=8760h

step certificate create \
    bindplane.internal \
    tls/bindplane.crt tls/bindplane.key \
    --profile leaf \
    --not-after 2160h \
    --no-password \
    --insecure \
    --ca tls/ca.crt \
    --ca-key tls/ca.key

sudo chown -R $USER tls/
sudo chmod -R 0644 tls/
