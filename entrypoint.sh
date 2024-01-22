#!/usr/bin/env bash

set -e

REMOTE_URL=$1
API_KEY=$2
USERNAME=$3
PASSWORD=$4
DESTINATION_PATH=$5
CONFIG_PATH=$6
OUTPUT_DIR=$7
OUTPUT_BRANCH=$8
GITHUB_TOKEN=$9

BRANCH_NAME=${GITHUB_REF#refs/heads/}
echo "Current branch is $BRANCH_NAME"

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

  if [ -n "$OUTPUT_BRANCH" ]; then
    if [ -z "$OUTPUT_DIR" ]; then
      echo "OUTPUT_DIR is required when OUTPUT_BRANCH is set."
      exit 1
    fi

    if [ -z "$GITHUB_TOKEN" ]; then
      echo "GITHUB_TOKEN is required when OUTPUT_BRANCH is set."
      exit 1
    fi

    # GITHUB_ACTOR and GITHUB_REPOSITORY are set by the github actions runtime
    # but we do expect them to be set.

    if [ -z "$GITHUB_ACTOR" ]; then
      echo "GITHUB_ACTOR is required when OUTPUT_BRANCH is set. This is likely a bug in the action. Please reach out to obserIQ support."
      exit 1
    fi

    if [ -z "$GITHUB_REPOSITORY" ]; then
      echo "GITHUB_REPOSITORY is required when OUTPUT_BRANCH is set. This is likely a bug in the action. Please reach out to obserIQ support."
      exit 1
    fi
  fi

  eval bindplane profile set "action" "$profile_args"
  bindplane profile use "action"
}

write_back() {
  if [ "$BRANCH_NAME" != "$OUTPUT_BRANCH" ]; then
    echo "Skipping repo write. Current branch ${BRANCH_NAME} does not match output branch ${OUTPUT_BRANCH}."
    return
  fi

  git clone https://$GITHUB_ACTOR:$GITHUB_TOKEN@github.com/$GITHUB_REPOSITORY.git

  mkdir -p "$GITHUB_REPOSITORY/$OUTPUT_DIR"

  for config in $(bindplane get config | awk 'NR>1 {print $1}'); do
    out_file="$GITHUB_REPOSITORY/$OUTPUT_DIR/$config.yaml"
    bindplane get config "$config" -o raw > "$out_file"
    git add "$out_file"
  done

  git config --global user.email "bindplane-op-action"
  git config --global user.name "bindplane-op-action"
  git commit -m "BindPlane OP Action: Update OTEL Configs"
  git push "https://$GITHUB_ACTOR:$GITHUB_TOKEN@github.com/$GITHUB_REPOSITORY.git" "HEAD:$OUTPUT_BRANCH"
}

install_bindplane_cli
validate
# Apply will apply resources in the correct order. Re-usable
# resources must exist before they can be referenced by
# a configuration.
bindplane apply "$DESTINATION_PATH"
bindplane apply "$CONFIG_PATH"
write_back

