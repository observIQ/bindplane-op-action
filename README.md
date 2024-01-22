# bindplane-op-action

## Usage

### Export Resources

```bash
bindplane get destination -o yaml --export > destination.yaml
bindplane get configuration -o yaml --export > configuration.yaml
```

### Workflow

```yaml
name: CI

# When raw config write back is configured, it is important
# to run this workflow only when changes to the resources
# are detected. This can prevent a CI infinite loop.
on:
  push:
    branches:
      - main
    paths:
      - 'test/resources/**'


permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          remote_url: ${{ secrets.BINDPLANE_REMOTE_URL }}
          api_key: ${{ secrets.BINDPLANE_API_KEY }}
```
