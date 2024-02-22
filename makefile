KEYCLOAK_HOST=http://localhost:8080
ADMIN_USER=admin
ADMIN_PASSWORD=admin
KEYCLOAK_CLIENT_ID=front
KEYCLOAK_REDIRECT_URIS=/api/auth/callback/keycloak
KEYCLOAK_WEB_ORIGIN=http://localhost:3000
KEYCLOAK_ROOT_URL=http://localhost:3000

init-keycloak: run-keycloak
	KEYCLOAK_HOST=$(KEYCLOAK_HOST) \
	ADMIN_USER=$(ADMIN_USER) \
	ADMIN_PASSWORD=$(ADMIN_PASSWORD) \
	KEYCLOAK_CLIENT_ID=$(KEYCLOAK_CLIENT_ID) \
	KEYCLOAK_REDIRECT_URIS=$(KEYCLOAK_REDIRECT_URIS) \
	KEYCLOAK_WEB_ORIGIN=$(KEYCLOAK_WEB_ORIGIN) \
	KEYCLOAK_ROOT_URL=$(KEYCLOAK_ROOT_URL) \
	sh ./scripts/init-keycloak.sh

run-keycloak:
	docker-compose up --build -d

run: init-keycloak
	echo "keycloak is running on $(KEYCLOAK_HOST)\n"

down: 
	docker-compose down

genproto:
	protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

genrestproto:
	protoc -I . --grpc-gateway_out . \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    proto/bff*.proto

genswag:
	protoc -I . --openapiv2_out ./openapiv2 --openapiv2_opt=allow_merge=true openapiv2/bff_v1.proto
lintproto:
	protolint lint proto/*.proto

.PHONY: create-network run

build_bff:
	docker build -t harbor.seafood-dev.com/dev/bff:latest -f cmd/bff/Dockerfile .

build_controller:
	docker build -t harbor.seafood-dev.com/dev/controller:latest -f cmd/controller/Dockerfile .

push_bff: build_bff
	docker push harbor.seafood-dev.com/dev/bff:latest

push_controller: build_controller
	docker push harbor.seafood-dev.com/dev/controller:latest
