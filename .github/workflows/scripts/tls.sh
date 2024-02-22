#!/usr/bin/env bash

set -e

curl -L -s -o step.tar.gz \
    https://dl.step.sm/gh-release/cli/gh-release-header/v0.22.0/step_linux_0.22.0_amd64.tar.gz
tar -xzf step.tar.gz
mv step_0.22.0/bin/step /usr/local/bin/step

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

sudo chmod -R 0644 /tmp/step

sudo ls -la /tmp/step

whoami
echo $USER
