module github.com/quic-go/quic-go

go 1.21

replace (
	github.com/quic-go/qtls-go1-20 => github.com/birneee/qtls-go1-20 v0.0.0-20231205151535-000383df7d5a
	github.com/tinylib/msgp => github.com/birneee/msgp v0.0.0-20240302001452-8e6189e615cd
)

require (
	github.com/francoispqt/gojay v1.2.13
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress v1.17.7
	github.com/onsi/ginkgo/v2 v2.9.5
	github.com/onsi/gomega v1.27.6
	github.com/quic-go/qpack v0.4.0
	github.com/stretchr/testify v1.6.1
	github.com/tinylib/msgp v1.1.8
	go.uber.org/mock v0.4.0
	golang.org/x/crypto v0.21.0
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225
	golang.org/x/net v0.22.0
	golang.org/x/sync v0.6.0
	golang.org/x/sys v0.18.0
	gonum.org/v1/gonum v0.14.0
	golang.org/x/time v0.5.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
