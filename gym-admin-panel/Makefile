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
    if exist protogen rmdir /s /q protogen
	PROTO_FILES = ../api/v1/auth/auth.proto ../api/v1/users/admin.proto ../api/v1/users/employee.proto ../api/v1/schedule/schedule.proto
	MAKE_DIR = mkdir protogen
else
	PROTO_FILES = $(shell find ../api/v1 -name '*.proto')
	MAKE_DIR = mkdir -p ./protogen
endif

# Генерация gRPC TypeScript Интерфейсов. 
gen:
	$(MAKE_DIR)
	npx protoc -I ../api $(PROTO_FILES) \
    --ts_out ./protogen \
    --ts_opt ts_nocheck \
    --ts_opt eslint_disable \

clean:
ifeq ($(SHELL_TYPE),windows)
	if exist protogen rmdir /s /q protogen
else
	rm -rf ./protogen
endif

