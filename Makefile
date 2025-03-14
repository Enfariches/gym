
# Генерация grpc кода, для работы на Go
gen:
	mkdir -p protogen
	protoc -I api api/proto/auth/auth.proto --go_out=./protogen/ --go_opt=paths=source_relative \
	--go-grpc_out=./protogen/ --go-grpc_opt=paths=source_relative