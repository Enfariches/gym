# Determine shell type
ifeq ($(OS),Windows_NT)
SHELL_TYPE = windows
else
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
SHELL_TYPE = linux
endif
ifeq ($(UNAME_S),Darwin)
SHELL_TYPE = darwin
endif
endif

# Set config path
CONFIG_PATH=$(CURDIR)/gym-server/config/dev.yaml

# Export config path
export CONFIG_PATH

# Set variables based on shell type
ifeq ($(SHELL_TYPE),windows)
PROTO_FILES = ../api/v1/auth/auth.proto ../api/v1/users/admin.proto ../api/v1/users/employee.proto ../api/v1/schedule/schedule.proto
MAKE_DIR = mkdir protogen
RM_CMD = if exist protogen rmdir /s /q protogen
PROTOC_CMD = npx protoc -I ../api $(PROTO_FILES) \
	--go_out=./protogen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./protogen/ \
	--go-grpc_opt=paths=source_relative
else
PROTO_FILES = $(shell find ../api/v1 -name '*.proto')
MAKE_DIR = mkdir -p ./protogen
PROTOC_CMD = protoc -I ../api $(PROTO_FILES) \
	--go_out=./protogen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./protogen/ \
	--go-grpc_opt=paths=source_relative
endif

# TODO: Сделать тесты.
# Написать тесты на SMTP; Auth

# Генерация gRPC Go Интерфейсов.
gen:
	$(RM_CMD)
	$(MAKE_DIR)
	$(PROTOC_CMD)

clean:
	$(RM_CMD)

# Путь конфига указывается двумя способами:
# 1) flag --config=./config/dev.yaml
# 2) С помощью env CONFIG_PATH, который экспортируется перед запуском make
# Если flag не указан, используется путь из env.
bundle:
	docker compose up gym-server

# go test ./... - для поиска тестов во всех пакетах
test:
	go mod tidy
	go test health/lib/jwt

migrate:
	go run ./cmd/migrator/migrator.go --storage-url=postgres:root@localhost:5436/health?sslmode=disable --migrations-path=migrations

opendb:
	docker exec -it postgres psql -U postgres health
