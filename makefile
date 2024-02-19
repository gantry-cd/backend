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
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

lintproto:
	protolint lint proto/**/*.proto

.PHONY: create-network run