.PHONY: docs
# Основной Makefile для иницализации и запуска Docker, third_party и основных сервисов.
# Возможно для запуска пайплайнов
PROTO_FILES = $(shell find ./api -name '*.proto')
DOCS_FILE = $(pwd)/docs:/usr/share/nginx/html

bundle-up:
	docker compose up -d

docs:
	mkdir -p ./docs
	protoc -I ./api --doc_out=./docs --doc_opt=html,index.html $(PROTO_FILES)
	docker compose up gym-docs -d --no-deps
	@echo "Документация доступна по адресу: http://localhost:8083"

clean-volumes:
	docker compose down --volumes --remove-orphans

bundle-gen:
	make -C gym-admin-panel gen
	make -C gym-server gen
	docker compose up --build --force-recreate

bundle:
	docker compose down
	docker compose build
	docker compose up

rebundle-admin:
	docker compose build gym-admin
	docker compose up -d --no-deps gym-admin
