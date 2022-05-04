test_storage:
	@go clean -testcache && go test -v ./internal/storage

db_up:
	@docker-compose up -d --build db

service_up:
	@docker-compose up -d url-service

app_run:
	@docker-compose up

run_server:
	@env GRPCPORT=50051 HOST="pgdb" DBNAME=urls DBPORT=5432 DBHOST=localhost DBUSERNAME=backend DBPASSWORD=user TIMEOUT=5 go run ./cmd/main.go

run_client:
	@env GRPCPORT=50051 HOST=localhost go run ./client/client.go

server_exec:
	@docker-compose exec -it urlservice sh

db_exec:
	@docker exec -it pgdb sh

generate_proto:
	@mkdir -p ./internal/pkg
	@protoc  api/proto/*.proto --go-grpc_out=internal/pkg --go_out=internal/pkg

mock_repo:
	@mockgen -destination internal/repositories/mock/mock_repository.go -source internal/repositories/repository.go

test_service:
	@go test -v ./internal/services/