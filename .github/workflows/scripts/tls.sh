#!/usr/bin/env bash

set -e

wget https://dl.smallstep.com/cli/docs-cli-install/latest/step-cli_amd64.deb
sudo apt-get install -y -f ./step-cli_amd64.deb

mkdir /tmp/step

ls -la /tmp/step

step certificate create \
    ca.internal \
    /tmp/step/ca.crt /tmp/step/ca.key \
    --profile root-ca \
    --no-password \
    --insecure \
    --not-after=8760h

step certificate create \
    bindplane.internal \
    /tmp/step/bindplane.crt /tmp/step/bindplane.key \
    --profile leaf \
    --not-after 2160h \
    --no-password \
    --insecure \
    --ca /tmp/step/ca.crt \
    --ca-key /tmp/step/ca.key

chmod -R 0644 /tmp/step
