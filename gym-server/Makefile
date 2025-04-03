export CONFIG_PATH=$(shell pwd)/gym-server/config/dev.yaml

# TODO: Сделать pretty тесты и логи.

# Генерация grpc кода, для работы c Go
gen:
	mkdir -p ./protogen
	protoc -I ../api ../api/v1/auth/auth.proto --go_out=./protogen/ --go_opt=paths=source_relative \
	--go-grpc_out=./protogen/ --go-grpc_opt=paths=source_relative

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