#!/usr/bin/env bash

set -e

curl -L -s -o step.tar.gz \
    https://dl.step.sm/gh-release/cli/gh-release-header/v0.22.0/step_linux_0.22.0_amd64.tar.gz
tar -xzf step.tar.gz
mv step_0.22.0/bin/step /usr/local/bin/step
rm -f step.tar.gz
rm -rf step_0.22.0

mkdir step/
chmod -R 0755 step/

step certificate create \
    ca.internal \
    step/ca.crt step/ca.key \
    --profile root-ca \
    --no-password \
    --insecure \
    --not-after=8760h

step certificate create \
    bindplane.internal \
    step/bindplane.crt step/bindplane.key \
    --profile leaf \
    --not-after 2160h \
    --no-password \
    --insecure \
    --ca step/ca.crt \
    --ca-key step/ca.key

chmod 0644 step/*
