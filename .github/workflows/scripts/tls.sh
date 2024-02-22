#!/usr/bin/env bash

set -e

wget https://dl.smallstep.com/cli/docs-cli-install/latest/step-cli_amd64.deb
sudo apt-get install -y -f ./step-cli_amd64.deb

mkdir step/

step certificate create \
    ca.internal \
    step/ca.crt step/ca.key \
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
    --ca step/ca.crt \
    --ca-key step/ca.key

chown -R $USER step/
chmod -R 0644 step/
