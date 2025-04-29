export CONFIG_PATH=$(shell pwd)/gym-server/config/dev.yaml
PROTO_FILES = $(shell find ../api/v1 -name '*.proto')


# TODO: Сделать тесты.
# Добавить docker volumes для сервера, чтобы не перебилдить каждый раз
# Исправить кроссплатформенность в make
# Наладить envoy

# Написать тесты на SMTP; Auth
# Добавить Password Reset
# Подумать как разграничить методы на уровне БД

# Генерация gRPC Go Интерфейсов.
gen:
	mkdir -p ./protogen
	protoc -I ../api $(PROTO_FILES) \
	 --go_out=./protogen/ \
	  --go_opt=paths=source_relative \
	--go-grpc_out=./protogen/ \
	 --go-grpc_opt=paths=source_relative

# Путь конфига указывается двумя способами:
# 1) flag --config=./config/dev.yaml
# 2) С помощью env CONFIG_PATH, который экспортируется перед запуском make
# Если flag не указан, используется путь из env.
bundle:
	go mod tidy
	go run ./cmd/gymserver/main.go --config=./config/dev.yaml

# go test ./... - для поиска тестов во всех пакетах
test:
	go mod tidy
	go test health/lib/jwt

migrate:
	go run ./cmd/migrator/migrator.go --storage-url=postgres:root@localhost:5436/health?sslmode=disable --migrations-path=migrations
