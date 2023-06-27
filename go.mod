module aviatus/vividora

go 1.20

replace (
	aviatus/vividora/api => ./api
	aviatus/vividora/internal/store => ./internel/store
)

require github.com/quic-go/quic-go v0.36.0

require (
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/pprof v0.0.0-20230602150820-91b7bce49751 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20230524184225-eabc099b10ab // indirect
	github.com/onsi/ginkgo/v2 v2.11.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.3.2 // indirect
	github.com/quic-go/qtls-go1-20 v0.3.0 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
)
