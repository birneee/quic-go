# QUIC-GO with XADS-QUIC Extension

> :warning: Experimental!

This repository contains a modified version of [quic-go](https://github.com/quic-go/quic-go).

## Security
Due to lack of security measures, this implementation is intended for research purposes only and should not be deployed on the internet.

## Changes to original QUIC-GO
- XADS-QUIC extension: additional TLS encryption of QUIC stream content

## Requirements
- Go 1.20

## Guides

Running tests:
```bash
go test ./...
```

Generating code:
```bash
go install github.com/golang/mock/mockgen@latest
go install golang.org/x/tools/cmd/goimports@latest
go generate ./...
```

Build:
```bash
go build .
```