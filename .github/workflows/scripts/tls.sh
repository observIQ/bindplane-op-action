#!/usr/bin/env bash

set -e

if [ -z "$MAIN_IP" ]; then
    echo "MAIN_IP is not set"
    exit 1
fi

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
    --san "${MAIN_IP}" \
    --profile leaf \
    --not-after 2160h \
    --no-password \
    --insecure \
    --ca step/ca.crt \
    --ca-key step/ca.key

chmod 0644 step/*
