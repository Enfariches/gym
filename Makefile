.PHONY: docs
# Основной Makefile для иницализации и запуска Docker, third_party и основных сервисов.
# Возможно для запуска пайплайнов
PROTO_FILES = $(shell find ./api/v1 -name '*.proto')
DOCS_FILE = $(pwd)/docs:/usr/share/nginx/html

bundle:
	docker compose up -d

docs:
	mkdir -p ./docs
	protoc --doc_out=./docs --doc_opt=html,index.html $(PROTO_FILES)
	docker compose up gym-docs -d --no-deps
	@echo "Документация доступна по адресу: http://localhost:8083"

clean-volumes:
	docker compose down --volumes --remove-orphans

rebundle:
	docker compose down
	docker compose build
	docker compose up -d

rebundle-admin:
	docker compose build gym-admin
	docker compose up -d --no-deps gym-admin
