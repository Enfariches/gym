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

ifeq ($(SHELL_TYPE),windows)
	PROTO_FILES = ../api/v1/auth/auth.proto ../api/v1/user/admin.proto ../api/v1/user/employee.proto ../api/v1/schedule/schedule.proto
	MAKE_DIR = mkdir protogen
else
	PROTO_FILES = '$(shell find ../api/v1 -name '*.proto')'
	MAKE_DIR = mkdir -p ./protogen
endif

# Генерация gRPC JS Интерфейсов.
gen:
	$(MAKE_DIR)
	protoc -I ../api $(PROTO_FILES) \
	--js_out=import_style=commonjs,binary:./protogen/ \
	--grpc-web_out=import_style=typescript,mode=grpcwebtext:./protogen/ \
	--ts_out=./protogen/

clean:
	if exist protogen rmdir /s /q protogen