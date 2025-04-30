# Gymnastic Monorepo

запуск сервера

make gen
docker compose up или (make bundle из gym)
make migrate 
make bundle (из gym-server)

(если ошибка protoc-gen-doc)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

пути для devs:

http://localhost:8083/ - docs
