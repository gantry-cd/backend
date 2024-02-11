run:
	sh ./scripts/run.sh

down:
	docker-compose down

genproto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/**/*.proto

lintproto:
	protolint lint proto/**/*.proto

.PHONY: create-network run