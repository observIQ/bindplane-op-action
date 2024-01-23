# bindplane-op-action

## Usage

### Export Resources

```bash
bindplane get destination -o yaml --export > destination.yaml
bindplane get configuration -o yaml --export > configuration.yaml
```

### Workflow

The following workflow can be used as an example.

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
          bindplane_api_key: "" # Optional replacement for username and password

          destination_path: test/resources/destinations/resource.yaml
          configuration_path: test/resources/configurations/resource.yaml

          # Write raw OTEL configs back to the repo.
          enable_otel_config_write_back: true
          configuration_output_dir: test/otel/${{ matrix.bindplane_versions }}
          configuration_output_branch: dev
          token: ${{ secrets.GITHUB_TOKEN }}
```
