#!/usr/bin/env bash

set -e

bindplane_remote_url=$1
bindplane_api_key=$2
bindplane_username=$3
bindplane_password=$4
destination_path=$5
configuration_path=$6
configuration_output_dir=$7
target_branch=$8
token=$9

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

  if [ -z "$bindplane_remote_url" ]; then
    echo "bindplane_remote_url is not set."
    exit 1
  else 
    profile_args="$profile_args --remote-url $bindplane_remote_url"
  fi

  if [ -n "$bindplane_username" ] && [ -z "$bindplane_password" ]; then
    echo "bindplane_password is required when bindplane_username is not set."
    exit 1
  else 
    profile_args="$profile_args --bindplane_username $bindplane_username --bindplane_password $bindplane_password"
  fi

  if [ -z "$bindplane_username" ] && [ -z "$bindplane_api_key" ]; then
    echo "api key is required when bindplane_username is not set."
    exit 1
  elif [ -n "$bindplane_api_key" ]; then
    profile_args="$profile_args --api-key $bindplane_api_key"
  fi

  if [ -n "$target_branch" ]; then
    if [ -z "$configuration_output_dir" ]; then
      echo "configuration_output_dir is required when target_branch is set."
      exit 1
    fi

    if [ -z "$token" ]; then
      echo "token is required when target_branch is set."
      exit 1
    fi

    # GITHUB_ACTOR and GITHUB_REPOSITORY are set by the github actions runtime
    # but we do expect them to be set.

    if [ -z "$GITHUB_ACTOR" ]; then
      echo "GITHUB_ACTOR is required when target_branch is set. This is likely a bug in the action. Please reach out to obserIQ support."
      exit 1
    fi

    if [ -z "$GITHUB_REPOSITORY" ]; then
      echo "GITHUB_REPOSITORY is required when target_branch is set. This is likely a bug in the action. Please reach out to obserIQ support."
      exit 1
    fi
  fi

  eval bindplane profile set "action" "$profile_args"
  bindplane profile use "action"
}

write_back() {
  # Clone the repo on the current branch
  # and use depth 1 to avoid cloning the entire history.
  git clone \
    --depth 1 \
    --branch "$BRANCH_NAME" \
    "https://${GITHUB_ACTOR}:${token}@github.com/${GITHUB_REPOSITORY}.git" \
    ../out_repo

  cd "../out_repo"

  mkdir -p "$configuration_output_dir"

  for config in $(bindplane get config | awk 'NR>1 {print $1}'); do
    out_file="$configuration_output_dir/$config.yaml"
    bindplane get config "$config" -o raw > "$out_file"
    git add "$out_file"
  done

  # check if git status is clean, return early
  if [[ -z $(git status --porcelain) ]]; then
    echo "No changes detected. Skipping commit."
    return
  fi

  git config --global user.email "bindplane-op-action"
  git config --global user.name "bindplane-op-action"
  git commit -m "BindPlane OP Action: Update OTEL Configs"
  git push
}

install_bindplane_cli
validate

if [ "$BRANCH_NAME" != "$target_branch" ]; then
  echo "Skipping apply and repo write. Current branch ${BRANCH_NAME} does not match target branch ${target_branch}."
else
  bindplane apply "$destination_path"
  bindplane apply "$configuration_path"
  write_back
fi


