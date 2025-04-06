module github.com/advdv/trustd

go 1.24.1

require (
	connectrpc.com/connect v1.18.1
	github.com/advdv/stdgo v0.0.112
	github.com/rs/cors v1.11.1
	go.uber.org/fx v1.23.0
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/caarlos0/env/v11 v11.3.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)
