#!/usr/bin/env bash

# This script initializes the Bindplane API with an organization or account
# depending on the version. Pre 1.59.0 was account based, post 1.59.0 is
# organization based.

set -ex

if [ -z "$1" ]; then
    echo "Usage: $0 <bindplane_version>"
    exit 1
fi

BINDPLANE_VERSION=$1
ORGS=false

echo "Initializing Bindplane $BINDPLANE_VERSION"

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

curl -v -k \
        -u admin:admin \
        https://localhost:3001/v1/source-types \
        --key step/bindplane.key \
        --cert step/bindplane.crt
sleep 20
