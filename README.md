# QUIC-GO with H-QUIC Extension

> :warning: Experimental!

This repository contains a modified version of [quic-go](https://github.com/quic-go/quic-go).

## Security
Due to lack of security measures, this implementation is intended for research purposes only and should not be deployed on the internet.

## Changes to original QUIC-GO
- client migration
  - tbd: address validation with path challenge
- server migration
  - tbd: address validation with path challenge
- H-QUIC extension
  - store and restore session
  - use non-transparent encryption-breaking proxies

## Guides

Running tests:
```bash
go test ./...
```

Generating code:
```bash
go install github.com/birneee/msgp@latest  # path mismatch requires manual git clone and go install
go install github.com/golang/mock/mockgen@latest
go install golang.org/x/tools/cmd/goimports@latest
go generate ./...
```
