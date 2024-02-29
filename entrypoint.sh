#!/usr/bin/env bash

set -e

# This order must match the order in acitons.yml
bindplane_remote_url=${1}
bindplane_api_key=${2}
bindplane_username=${3}
bindplane_password=${4}
target_branch=${5}
destination_path=${6}
configuration_path=${7}
enable_otel_config_write_back=${8}
configuration_output_dir=${9}
token=${10}
enable_auto_rollout=${11}
configuration_output_branch=${12}
tls_ca_cert=${13}

# This branch name will be compared to target_branch to determine if the action
# should apply or write back configurations.
BRANCH_NAME=${GITHUB_REF#refs/heads/}
echo "Current branch is $BRANCH_NAME"

install_bindplane_cli() {
  curl -Ls \
    -o bindplane.zip \
    https://storage.googleapis.com/bindplane-op-releases/bindplane/1.46.0/bindplane-ee-linux-amd64.zip

  mkdir -p ~/bin
  export PATH=$PATH:~/bin

  unzip bindplane.zip -d ~/bin > /dev/null

  bindplane --help > /dev/null
}

# Validate will ensure that all required variables are set
# and generates the bindplane profile.
validate() {
  profile_args=""

  # Target branch is always required. When not set, the script will not
  # know which branch it should apply configurations from or write back
  # raw otel configs.
  if [ -z "$target_branch" ]; then
    echo "target_branch is required when enable_otel_config_write_back is true."
    exit 1
  fi

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
    profile_args="$profile_args --username $bindplane_username --password $bindplane_password"
  fi

  if [ -z "$bindplane_username" ] && [ -z "$bindplane_api_key" ]; then
    echo "api key is required when bindplane_username is not set."
    exit 1
  elif [ -n "$bindplane_api_key" ]; then
    profile_args="$profile_args --api-key $bindplane_api_key"
  fi

  if [ -n "$tls_ca_cert" ]; then
    echo "tls_ca_cert is set, adding to profile."
    echo "$tls_ca_cert" > ca.pem
    profile_args="$profile_args --tls-ca ca.pem"
  fi

  # configuration_output_dir, target_branch, and token are only required
  # when enable_otel_config_write_back is true.
  if [ "$enable_otel_config_write_back" = true ]; then
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
  # The configuration_output_branch is optional. If not set, the
  # write back branch will be the same as the target branch.
  write_back_branch=${configuration_output_branch:-$target_branch}

  # Clone the repo on the current branch
  # and use depth 1 to avoid cloning the entire history.
  git clone \
    --depth 1 \
    --branch "$write_back_branch" \
    "https://${GITHUB_ACTOR}:${token}@github.com/${GITHUB_REPOSITORY}.git" \
    ../out_repo

  cd "../out_repo"

  mkdir -p "$configuration_output_dir"

  for config in $(bindplane get config | awk 'NR>1 {print $1}'); do
    out_file="$configuration_output_dir/$config.yaml"
    # It is safe to always ask for "latest". BindPlane will return
    # the current version if there is no latest version.
    bindplane get config "${config}:latest" -o raw > "$out_file"
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

main() {
  # Short circuit if the current branch does not match the target branch,
  # there is nothing to do.
  if [ "$BRANCH_NAME" != "$target_branch" ]; then
    echo "Skipping apply and repo write. Current branch ${BRANCH_NAME} does not match target branch ${target_branch}."
    exit 0
  fi

  echo "Target branch ${target_branch} matches current branch ${BRANCH_NAME}."

  # Install the CLI right away in order to construct
  # a client profile.
  echo "Installing bindplane CLI."
  install_bindplane_cli

  # Ensure required options are set and configure
  # the client profile.
  echo "Validating options and configuring client profile."
  validate

  # Apply resources in the correct order.
  echo "Applying resources."

  echo "Applying destination path: $destination_path"
  bindplane apply "$destination_path"

  echo "Applying configuration path: $configuration_path"
  bindplane apply "$configuration_path" > configuration.out
  cat configuration.out

  # When auto rollout is enabled
  if [ "$enable_auto_rollout" = true ]; then
    echo "Auto rollout enabled."
    awk '{print $2}' < configuration.out | while IFS= read -r config
    do
      status=$(bindplane rollout status "${config}" -o json | jq .status)
      case "$status" in
        0)
          echo "Configuration ${config} has a pending rollout, triggering rollout."
          bindplane rollout start "$config"
          ;;
        4)
          echo "Configuration ${config} is stable, skipping rollout."
          ;;
        *)
          echo "Configuration ${config} has an unknown status, skipping rollout."
          ;;
      esac
    done
  fi

  # When write back is enabled, write the raw otel configs
  # back to the repository.
  if [ "$enable_otel_config_write_back" = true ]; then
    echo "Writing back raw otel configs."
    write_back
  fi

  echo "Done."
}

main
