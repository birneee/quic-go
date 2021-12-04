# Fork of QUIC-GO

## Security
Due to lack of security measures, this implementation is intended for research purposes only and should not be deployed on the internet.

## Changes to Standard QUIC-GO
- client migration
- server migration
- store and restore session
- use encryption breaking proxy

## Guides

*We currently support Go 1.16.x and Go 1.17.x.*

Running tests:
```bash
go test ./...
```

Generating code:
```bash
go generate ./...
```
