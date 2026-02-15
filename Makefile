include .env

proto-format:
	buf format -w

proto-lint:
	buf lint

proto-gen:
	buf generate

MOCK_GEN=go run go.uber.org/mock/mockgen@v0.6.0

mock-gen:
	$(MOCK_GEN) -source=internal/domain/event/event.go -destination=mocks/domain/event/mock_event.go -package=mocks
	$(MOCK_GEN) -source=internal/domain/event/repository.go -destination=mocks/domain/event/mock_repository.go -package=mocks
	$(MOCK_GEN) -source=internal/application/event/usecase.go -destination=mocks/application/event/mock_usecase.go -package=mocks
	$(MOCK_GEN) -source=internal/application/auth/service.go -destination=mocks/application/auth/mock_service.go -package=mocks
	$(MOCK_GEN) -source=internal/application/user/service.go -destination=mocks/application/user/mock_service.go -package=mocks
	$(MOCK_GEN) -destination=mocks/external/auth/v1/mock_client.go -package=mocks github.com/qkitzero/auth-service/gen/go/auth/v1 AuthServiceClient
	$(MOCK_GEN) -destination=mocks/external/user/v1/mock_client.go -package=mocks github.com/qkitzero/user-service/gen/go/user/v1 UserServiceClient

MIGRATE=migrate -source file://internal/infrastructure/db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_HOST_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

migrate-up:
	$(MIGRATE) up

migrate-up-one:
	$(MIGRATE) up 1

migrate-down:
	$(MIGRATE) down 1

migrate-reset:
	$(MIGRATE) drop -f

migrate-create:
	migrate create -ext sql -dir internal/infrastructure/db/migrations -format 20060102150405 $(name)

migrate-status:
	$(MIGRATE) version

test:
	mkdir -p tmp
	go test -cover ./internal/... -coverprofile=./tmp/cover.out
	go tool cover -func=./tmp/cover.out | tail -n 1
	go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	open ./tmp/cover.html