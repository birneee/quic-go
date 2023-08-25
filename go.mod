module github.com/quic-go/quic-go

go 1.20

replace (
	github.com/quic-go/qtls-go1-20 => github.com/birneee/qtls-go1-20 v0.0.0-20230822163111-cfa94cb80061
	github.com/tinylib/msgp => github.com/birneee/msgp v0.0.0-20230807002656-18d07944fa3d
)

require (
	github.com/francoispqt/gojay v1.2.13
	github.com/golang/mock v1.6.0
	github.com/json-iterator/go v1.1.12
	github.com/onsi/ginkgo/v2 v2.9.5
	github.com/onsi/gomega v1.27.6
	github.com/quic-go/qpack v0.4.0
	github.com/quic-go/qtls-go1-20 v0.3.2
	github.com/stretchr/testify v1.6.1
	github.com/tinylib/msgp v1.1.8
	golang.org/x/crypto v0.4.0
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db
	golang.org/x/net v0.10.0
	golang.org/x/sync v0.2.0
	golang.org/x/sys v0.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
