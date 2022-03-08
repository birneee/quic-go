# QUIC-GO with H-QUIC and XSE-QUIC Extension

> :warning: Experimental!

This repository contains a modified version of [quic-go](https://github.com/lucas-clemente/quic-go).

## Security
Due to lack of security measures, this implementation is intended for research purposes only and should not be deployed on the internet.

## Changes to original QUIC-GO
- client migration
  - tbd: address validation with path challenge
- server migration
  - tbd: address validation with path challenge
- change udp socket during live session
- options to set initial, minimum and maximum congestion window
- additional qlog events
  - path updates (connection migration)
- H-QUIC extension
  - store and restore session
  - use non-transparent encryption-breaking proxies
- XSE-QUIC extension
  - additional encryption QUIC stream content
  - additional qlog events
    - received XSE records (TLS records)

## Guides

Running tests:
```bash
go test ./...
```

Generating code:
```bash
go generate ./...
```
