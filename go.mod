module github.com/kulycloud/api-server

go 1.15

require (
	github.com/kulycloud/common v1.0.0
	github.com/kulycloud/protocol v1.0.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	google.golang.org/grpc v1.32.0
)

replace github.com/kulycloud/common v1.0.0 => ../common

replace github.com/kulycloud/protocol v1.0.0 => ../protocol
