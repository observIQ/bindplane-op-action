#!/usr/bin/env bash

# This script initializes the BindPlane API with an organization or account
# depending on the version. Pre 1.59.0 was account based, post 1.59.0 is
# organization based.

set -ex

BINDPLANE_VERSION=$1
ORGS=false

if [ -z "$BINDPLANE_VERSION" ]; then
    echo "Usage: $0 <BINDPLANE_VERSION>"
fi

echo "Initializing BindPlane $BINDPLANE_VERSION"

if [[ $(printf '%s\n' "$BINDPLANE_VERSION" "1.58.0" | sort -V | head -n1) == "1.58.0" && "$BINDPLANE_VERSION" != "1.58.0" ]]; then
    ORGS=true
elif [[ "$BINDPLANE_VERSION" == "latest" ]]; then
    ORGS=true
fi

if [[ $ORGS == "true" ]]; then
    curl -v -k \
        -u admin:admin \
        https://localhost:3001/v1/organizations \
        -X POST \
        -d '{"organizationName": "init", "accountName": "project", "eulaAccepted":true}' \
        --key step/bindplane.key \
        --cert step/bindplane.crt
else
    curl -v -k \
        -u admin:admin \
        https://localhost:3001/v1/accounts \
        -X POST \
        -d '{"displayName": "init"}' \
        --key step/bindplane.key \
        --cert step/bindplane.crt
fi
