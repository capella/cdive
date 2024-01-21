---
# This workflow will build a golang project
#
name: Go
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21.x
      - name: Build
        run: go build -v ./...
      - name: Test
        run: go test -v ./...
  lint:
    permissions:
      # for actions/checkout to fetch code
      contents: read
      # for golangci/golangci-lint-action to fetch pull requests
      pull-requests: read
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21.x
      - name: Check Go module tidiness and generated files
        shell: bash
        run: |
          go mod tidy
          STATUS=$(git status --porcelain)
          if [ ! -z "$STATUS" ]; then
            echo "Unstaged files:"
            echo $STATUS
            echo "Run 'go mod tidy' commit them"
            exit 1
          fi
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout=10m
