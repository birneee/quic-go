# Fork of QUIC-GO

This repository contains a modified version of [quic-go](https://github.com/lucas-clemente/quic-go).

## Security
Due to lack of security measures, this implementation is intended for research purposes only and should not be deployed on the internet.

## Changes to original QUIC-GO
- H-QUIC
  - client migration
  - server migration
  - store and restore session
  - use encryption breaking proxy
- XSE-QUIC
  - todo

## Guides

Running tests:
```bash
go test ./...
```

Generating code:
```bash
go generate ./...
```
