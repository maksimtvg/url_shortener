test_storage:
	@go clean -testcache && go test -v ./internal/storage

run_server:
	@env GRPCPORT=50051 HOST=localhost go run ./cmd/main.go

run_client:
	@env GRPCPORT=50051 HOST=localhost go run ./cmd/client.go

#test_wallet:
#	@go test -v ./internal/app/services

generate_proto:
	@mkdir -p ./internal/pkg
	@protoc  api/proto/*.proto --go-grpc_out=internal/pkg --go_out=internal/pkg
