# bindplane-op-action

## Usage

### Export Resources

```bash
bindplane get destination -o yaml --export > destination.yaml
bindplane get configuration -o yaml --export > configuration.yaml
```

### Workflow

```yaml
name: goreleaser

on:
  pull_request:
  push:

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
