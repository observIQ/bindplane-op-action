#!/usr/bin/env bash

set -e

REMOTE_URL=$1
API_KEY=$2
USERNAME=$3
PASSWORD=$4
DESTINATION_PATH=$5
CONFIG_PATH=$6

install_bindplane_cli() {
  curl -Ls \
    -o bindplane.zip \
    https://storage.googleapis.com/bindplane-op-releases/bindplane/latest/bindplane-ee-linux-amd64.zip

  mkdir -p ~/bin
  export PATH=$PATH:~/bin

  unzip bindplane.zip -d ~/bin

  bindplane --help > /dev/null
}

# Validate will ensure that all required variables are set
# and generates the bindplane profile.
validate() {
  profile_args=""

  if [ -z "$REMOTE_URL" ]; then
    echo "REMOTE_URL is not set."
    exit 1
  else 
    profile_args="$profile_args --remote-url $REMOTE_URL"
  fi

  if [ -n "$USERNAME" ] && [ -z "$PASSWORD" ]; then
    echo "password is required when username is not set."
    exit 1
  else 
    profile_args="$profile_args --username $USERNAME --password $PASSWORD"
  fi

  if [ -z "$USERNAME" ] && [ -z "$API_KEY" ]; then
    echo "api key is required when username is not set."
    exit 1
  elif [ -n "$API_KEY" ]; then
    profile_args="$profile_args --api-key $API_KEY"
  fi  

  eval bindplane profile set "action" "$profile_args"
  bindplane profile use "action"
}

install_bindplane_cli
validate
# Apply will apply resources in the correct order. Re-usable
# resources must exist before they can be referenced by
# a configuration.
bindplane apply "$DESTINATION_PATH"
bindplane apply "$CONFIG_PATH"
