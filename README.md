[![CI](https://github.com/observIQ/bindplane-op-action/actions/workflows/ci.yml/badge.svg)](https://github.com/observIQ/bindplane-op-action/actions/workflows/ci.yml)

# bindplane-op-action

The BindPlane OP action can be used to deploy configurations to your BindPlane OP
server. It also supports exporting the OpenTelemetry configurations back to the repository.

## Configuration

| Parameter                     | Default    | Description                     |
| :---------------------------- | :--------- | :------------------------------ |
| bindplane_remote_url          | required   | The endpoint that will be used to connect to BindPalne OP. |
| bindplane_api_key             |            | API key used to authenticate to BindPlane. Required when BindPlane multi account is enabled or when running on BindPlane Cloud |
| bindplane_username            |            | Username used to authenticate to BindPlane. Not required if API key is set. |
| bindplane_password            |            | Password used to authenticate to BindPlane.
| target_branch                 | required   | The branch that the action will use when applying resources to bindplane or when writing otel configs back to the repo. |
| destination_path              |            | Path to the file which contains the BindPlane destination resources |
| configuration_path            |            | Path to the file which contains the BindPlane configuration resources |
| enable_otel_config_write_back | `false`    | Whether or not the action should write the raw OpenTelemetry configurations back to the repository. | 
| configuration_output_dir      |            | When write back is enabled, this is the path that will be written to. |
| token                         |            | The Github token that will be used to write to the repo. Usually secrets.GITHUB_TOKEN is sufficient. Requires the `contents.write` permission. |

## Usage

### Export Resources

To get started, you must handle exporting your existing resources to the repository. Use
the `bindplane get` commands with the `--export` flag.

```bash
bindplane get destination -o yaml --export > destination.yaml
bindplane get configuration -o yaml --export > configuration.yaml
```

With the resources exported to the repository, you can move on to configuring the action
using a new workflow.

### Workflow

The following workflow can be used as an example. It uses the same file paths
created in the [Export Resources](#export-resources) section.

This example will write the raw OTEL configurations back to the repository at the
path `otel/`

```yaml
name: bindplane

on:
  push:
    branches:
      - main

# Write back requires access to the repo
permissions:
  contents: write

# Run commits in order to prevent out of order write back commits.
concurrency:
  group: ${{ github.head_ref || github.ref_name }}
  cancel-in-progress: false

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: observIQ/bindplane-op-action@main
        with:
          bindplane_remote_url: ${{ secrets.BINDPLANE_REMOTE_URL }}
          bindplane_username: ${{ secrets.BINDPLANE_USERNAME }}
          bindplane_password: ${{ secrets.BINDPLANE_PASSWORD }}
          target_branch: main
          destination_path: destination.yaml
          configuration_path: configuration.yaml
          enable_otel_config_write_back: true
          configuration_output_dir: otel/
          token: ${{ secrets.GITHUB_TOKEN }}
```

After the action is executed, you can expect to see OTEL configurations
in the `otel/` directory.

```
otel
├── k8s-cluster.yaml
├── k8s-gateway.yaml
└── k8s-node.yaml
```
